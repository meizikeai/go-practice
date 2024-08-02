package routers

import (
	"go-practice/controllers"

	"github.com/gin-gonic/gin"
)

func AddApiRouter(router *gin.Engine) {
	router.GET("/home", controllers.Home)

	r := router.Group("/api")

	r.Use(controllers.ApiAuth())
	{
		r.GET("/test", controllers.ApiTest)
		r.POST("/test", controllers.ApiTest)
	}
}
