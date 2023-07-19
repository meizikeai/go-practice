package main

import (
	"fmt"
	"os"

	"go-practice/controllers"
	"go-practice/libs/log"
	"go-practice/libs/tool"
	"go-practice/routers"

	"github.com/gin-gonic/gin"
)

func init() {
	// tool.HandleZookeeperConfig()

	// tool.HandleMySQLClient()
	// tool.HandleRedisClient()
	// tool.HandleKafkaProducerClient()
	// tool.HandleKafkaConsumerClient()

	log.HandleLogger("go-practice")
}

func main() {
	tool.SignalHandler(func() {
		tool.CloseKafka()
		tool.CloseMySQL()
		tool.CloseRedis()

		tool.Stdout("Server Shutdown")

		os.Exit(0)
	})

	router := gin.Default()

	router.Use(controllers.TraceLogger())

	routers.HandleRouter(router)

	// kafka consumer
	// tool.HandlerKafkaConsumerMessage("broker", "topic")

	fmt.Println("Listen and Server running on 127.0.0.1:8000")
	router.Run(":8000")
}
