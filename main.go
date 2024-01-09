package main

import (
	"os"

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
		tool.CloseMySQL()
		tool.CloseRedis()

		tool.Stdout("The Service is Shutdown")

		os.Exit(0)
	})

	router := gin.Default()

	router.Use(controllers.TraceLogger())

	routers.HandleRouter(router)

	tool.Stdout("The service is running on 127.0.0.1:8000")
	router.Run(":8000")
}
