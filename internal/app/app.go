// internal/app/app.go
package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"go-practice/internal/app/repository"
	"go-practice/internal/app/service"
	"go-practice/internal/config"
	"go-practice/internal/pkg/crypto"
	"go-practice/internal/pkg/database/cache"
	"go-practice/internal/pkg/database/kafka"
	"go-practice/internal/pkg/database/mysql"
	"go-practice/internal/pkg/jwt"
	"go-practice/internal/pkg/log"
	"go-practice/internal/pkg/prometheus"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	Crypto     *crypto.Manager
	Engine     *gin.Engine
	Jwt        *jwt.Manager
	Log        *zap.Logger
	cfg        *config.Config
	server     *http.Server
	Cache      Storage
	DB         Storage
	Kafka      Storage
	Repository repository.Repository
	Service    service.Service
}

type Storage interface {
	Close()
}

func NewApp(mock ...Mock) *App {
	cfg := config.Load()

	cache := cache.NewClient(&cfg.Redis)
	data := mysql.NewClient(&cfg.MySQL)
	kafka := kafka.NewClient(&cfg.Kafka)
	crypto, _ := crypto.NewManager(&cfg.CryptoKey)
	jwt, _ := jwt.NewManager(&cfg.JwtKey)
	record := log.Load(cfg.App.Name, cfg.App.Mode)

	app := new(App)

	app.cfg = cfg
	app.Crypto = crypto
	app.Jwt = jwt
	app.Log = record
	app.cacheClient(data, cache, kafka)

	app.Repository = repository.NewRepository(record, data, cache)
	app.Service = service.NewService(record, app.Repository)

	// used for testing(mock)
	for _, fn := range mock {
		fn(app)
	}

	// gin
	gin.SetMode(cfg.App.Mode)
	app.Engine = gin.New()
	app.Engine.Use(
		log.AccessLog(),
		log.RecoveryWithZap(record),
		prometheus.PromMiddleware(&prometheus.PromOpts{
			ExcludeRegexStatus: "404",
			EndpointLabelMappingFn: func(c *gin.Context) string {
				route := c.FullPath()
				if route == "" {
					route = "/unknown"
				}
				return route
			},
		}),
	)

	return app
}

func (a *App) cacheClient(db, cache, kafka Storage) {
	a.DB = db
	a.Cache = cache
	a.Kafka = kafka
}

func (a *App) Run() {
	a.server = &http.Server{
		Addr:    net.JoinHostPort("", a.cfg.App.Port),
		Handler: a.Engine,
	}

	go func() {
		a.Stdout("Application initialization started in " + a.cfg.App.Mode + " environment")
		a.Stdout("Application started successfully and listening on 127.0.0.1" + a.server.Addr)

		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Stdout("Application start failed", err)
		}
	}()
}

func (a *App) Shutdown(ctx context.Context) error {
	a.Stdout("Service shutdown initiated")

	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			a.Stdout("Service forced to shut down", err)
		}
	}

	var errs []error

	if a.Cache != nil {
		a.Stdout("Redis connection closed")
		if err := a.closeStorage(a.Cache, "Redis"); err != nil {
			errs = append(errs, err)
		}
	}

	if a.DB != nil {
		a.Stdout("MySQL connection closed")
		if err := a.closeStorage(a.DB, "MySQL"); err != nil {
			errs = append(errs, err)
		}
	}

	if a.Kafka != nil {
		a.Stdout("Kafka connection closed")
		if err := a.closeStorage(a.Kafka, "Kafka"); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("Service connection closed with error -> %v", errs)
	}

	a.Stdout("Service exited")

	return nil
}

func (a *App) Stdout(format string, v ...any) {
	log := fmt.Sprintf("%s %s %s \n", time.Now().Format(time.DateTime), fmt.Sprintf("[%s]", a.cfg.App.Name), format)

	if _, err := fmt.Fprintf(os.Stdout, log, v...); err != nil {
		fmt.Println(log)
	}
}

func (a *App) closeStorage(s Storage, name string) error {
	defer func() {
		if r := recover(); r != nil {
			a.Stdout("Panic when closing -> %s: %v", name, r)
		}
	}()
	s.Close()
	return nil
}

// unit testing
type Mock func(*App)

func WithRepository(r repository.Repository) Mock {
	return func(a *App) {
		a.Repository = r
	}
}

func WithService(s service.Service) Mock {
	return func(a *App) {
		a.Service = s
	}
}
