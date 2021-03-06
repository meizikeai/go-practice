package main

import (
	"os"
	"time"

	"go-practice/libs/log"
	"go-practice/libs/tool"
	"go-practice/routers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	tool.HandleZookeeperConfig()

	// not recommended for use
	// tool.HandleLocalMysqlConfig()
	// tool.HandleLocalRedisConfig()

	tool.HandleMySQLClient()
	tool.HandleRedisClient()

	log.HandleLogger("go-practice")
}

func main() {
	pwd, _ := os.Getwd()
	router := gin.Default()

	router.Static("/public", pwd+"/public")
	router.StaticFile("/favicon.ico", pwd+"/public/favicon.ico")

	router.LoadHTMLGlob("views/*")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.Use(log.AccessLogger("go-practice"))

	routers.HandleRouter(router)

	router.Run(":8080")
}
