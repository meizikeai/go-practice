package routers

import (
	"go-practice/controllers"

	"github.com/gin-gonic/gin"
)

func AddApiRouter(router *gin.Engine) {
	router.NoRoute(controllers.NotFound)
	router.NoMethod(controllers.NotFound)

	router.GET("/", controllers.SayHi)
	router.GET("/favicon.ico", controllers.SayHi)

	router.GET("/home", controllers.Home)

	r := router.Group("/api")

	r.Use(controllers.ApiAuth())
	{
		r.GET("/test", controllers.ApiTest)
		r.POST("/test", controllers.ApiTest)
	}
}
