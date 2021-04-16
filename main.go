package main

import (
	"go-practice/libs/log"
	"go-practice/libs/tool"
	"go-practice/routers"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	tool.HandleZookeeperConfig()

	tool.HandleLocalMysqlConfig()
	tool.HandleLocalRedisConfig()

	tool.HandleMySQLClient()
	tool.HandleRedisClient()

	log.HandleLogger("go-practice")
}

func main() {
	pwd, _ := os.Getwd()
	router := gin.Default()

	// 静态资源
	router.Static("/public", pwd+"/public")
	router.StaticFile("/favicon.ico", pwd+"/public/favicon.ico")

	// 静态模板
	router.LoadHTMLGlob("views/*")

	// 跨域共享
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

	// 访问记录
	router.Use(log.AccessLogger("go-practice"))

	// 站点路由
	routers.HandleRouter(router)

	// 监听端口
	router.Run(":8080")
}
