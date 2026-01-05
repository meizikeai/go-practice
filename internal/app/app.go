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
	"go-practice/internal/pkg/fetch"
	"go-practice/internal/pkg/jwt"
	"go-practice/internal/pkg/log"
	"go-practice/internal/pkg/prometheus"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	CFG        *config.Config
	server     *http.Server
	Crypto     *crypto.Manager
	Engine     *gin.Engine
	Jwt        *jwt.Manager
	Log        *zap.Logger
	cache      Storage
	db         Storage
	kafka      Storage
	Repository repository.Repository
	Service    service.Service
}

type Storage interface {
	Close()
}

func NewApp(mock ...Mock) *App {
	cfg := config.Load()

	cache := cache.NewClient(&cfg.Redis)
	db := mysql.NewClient(&cfg.MySQL)
	fetch := fetch.NewClient()
	kafka := kafka.NewClient(&cfg.Kafka)
	record := log.Load(cfg.App.Name, cfg.App.Mode)
	crypto, _ := crypto.NewManager(&cfg.CryptoKey)
	jwt, _ := jwt.NewManager(&cfg.JwtKey)

	app := new(App)

	app.CFG = cfg
	app.Crypto = crypto
	app.Jwt = jwt
	app.Log = record
	app.cacheClient(db, cache, kafka)

	app.Repository = repository.New(
		cache,
		db,
		fetch,
		record,
		cfg.LB,
	)
	app.Service = service.New(
		cfg.App.Mode,
		app.Repository,
	)

	// used for testing(mock)
	for _, fn := range mock {
		fn(app)
	}

	// gin
	gin.SetMode(cfg.App.Mode)
	app.Engine = gin.New()
	app.Engine.Use(
		log.AccessLog(),
		gin.Recovery(),
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
	a.db = db
	a.cache = cache
	a.kafka = kafka
}

func (a *App) Run() {
	a.server = &http.Server{
		Addr:    net.JoinHostPort("", a.CFG.App.Port),
		Handler: a.Engine,
	}

	go func() {
		a.Stdout("Application initialization started in " + a.CFG.App.Mode + " environment")
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

	if a.cache != nil {
		a.Stdout("Redis connection closed")
		if err := a.closeStorage(a.cache, "Redis"); err != nil {
			errs = append(errs, err)
		}
	}

	if a.db != nil {
		a.Stdout("MySQL connection closed")
		if err := a.closeStorage(a.db, "MySQL"); err != nil {
			errs = append(errs, err)
		}
	}

	if a.kafka != nil {
		a.Stdout("Kafka connection closed")
		if err := a.closeStorage(a.kafka, "Kafka"); err != nil {
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
	log := fmt.Sprintf("%s %s %s \n", time.Now().Format(time.DateTime), fmt.Sprintf("[%s]", a.CFG.App.Name), format)

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
