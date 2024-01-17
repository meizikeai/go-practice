package main

import (
	"os"

	"go-practice/config"
	"go-practice/controllers"
	"go-practice/libs/log"
	"go-practice/libs/tool"
	"go-practice/routers"

	"github.com/gin-gonic/gin"
)

func init() {
	// tool.HandleMySQLClient()
	// tool.HandleRedisClient()

	log.HandleLogger("go-practice")
}

func main() {
	tool.SignalHandler(func() {
		// tool.CloseMySQL()
		// tool.CloseRedis()

		tool.Stdout("The Service is Shutdown")

		os.Exit(0)
	})

	// gin
	router := gin.New()

	// logger
	router.Use(controllers.TraceLogger())

	// recovery
	router.Use(gin.Recovery())

	routers.HandleRouter(router)

	port := config.GetPort()
	tool.Stdout("The current environment is " + config.GetMode())
	tool.Stdout("The service is running on 127.0.0.1:" + port)
	router.Run(":" + port)
}
