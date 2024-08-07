package main

import (
	"os"

	"go-practice/config"
	"go-practice/controllers"
	"go-practice/libs/log"
	"go-practice/libs/tool"
	"go-practice/libs/utils"
	"go-practice/routers"

	"github.com/gin-gonic/gin"
)

var (
	// rules = tool.NewRegexp()
	// share = tool.NewShare()
	// tasks = crontab.NewTasks()
	tools = tool.NewTools()
	// units = tool.NewUnits()
	// chaos  = secret.NewSecret()
	daily = log.NewCreateLog()
	// fetch  = models.NewModelsFetch()
	// jwt    = token.NewJsonWebToken()
	// lion   = fetch.NewFetch()
	logger = log.NewLogger()
	logic  = controllers.NewLogic()
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

		tools.Stdout("The Service is Shutdown")

		os.Exit(0)
	})

	// gin
	router := gin.New()

	// logger
	router.Use(logger.TraceLogger())

	// prometheus
	router.Use(utils.PromMiddleware(&utils.PromOpts{
		ExcludeRegexStatus: "404",
		EndpointLabelMappingFn: func(c *gin.Context) string {
			return logic.EndpointLabelMappingFn(c)
		},
	}))

	// recovery
	router.Use(gin.Recovery())

	routers.HandleRouter(router)

	port := config.GetPort()
	tools.Stdout("The current environment is " + config.GetMode())
	tools.Stdout("The service is running on 127.0.0.1:" + port)
	router.Run(":" + port)
}
