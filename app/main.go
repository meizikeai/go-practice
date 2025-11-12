package main

import (
	"os"

	"go-practice/config"
	"go-practice/controller"
	"go-practice/libs/tool"
	"go-practice/libs/utils"
	"go-practice/router"

	"github.com/gin-gonic/gin"
)

var (
	daily  = tool.NewCreateLog()
	logger = tool.NewLogger()
	tools  = tool.NewTools()
	logic  = controller.NewLogic()
	// chaos         = tool.NewSecret()
	// jwt           = tool.NewJsonWebToken()
	// lion          = tool.NewFetch()
	// rules         = tool.NewRegexp()
	// share         = tool.NewShare()
	// units         = tool.NewUnits()
	// tasks         = crontab.NewTasks()
	// fetch         = models.NewModelsFetch()
	// kafkaProducer = tool.NewKafkaProducer()
	// kafkaConsumer = tool.NewKafkaConsumer()
)

func init() {
	// tools.HandleMySQLClient()
	// tools.HandleRedisClient()
	daily.HandleLogger("go-practice")
}

func main() {
	tools.SignalHandler(func() {
		// tools.CloseMySQL()
		// tools.CloseRedis()
		tools.Stdout("Service shut down")
		os.Exit(0)
	})

	// gin
	app := gin.New()

	// log
	app.Use(logger.TraceLogger())

	// prometheus
	app.Use(utils.PromMiddleware(&utils.PromOpts{
		ExcludeRegexStatus: "404",
		EndpointLabelMappingFn: func(c *gin.Context) string {
			return logic.EndpointLabelMappingFn(c)
		},
	}))

	// recovery
	app.Use(gin.Recovery())

	router.HandleRouter(app)

	port := config.GetPort()
	tools.Stdout("Starting application in the " + config.GetMode() + " environment")
	tools.Stdout("Application started successfully. Listening on 127.0.0.1:" + port)
	app.Run(":" + port)
}
