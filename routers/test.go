package routers

import (
	"go-practice/controllers"
	"go-practice/libs/utils"

	"github.com/gin-gonic/gin"
)

func AddTestApiRouter(router *gin.Engine) {
	router.GET("/home", controllers.Home)

	r := router.Group("/api")

	r.Use(utils.ApiAuth())
	{
		r.GET("/test", controllers.ApiTest)
		r.POST("/test", controllers.ApiTest)
	}
}
