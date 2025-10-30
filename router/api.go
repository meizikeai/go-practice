package router

import (
	"github.com/gin-gonic/gin"
)

func addApiRouter(router *gin.Engine) {
	router.GET("/home", logic.Home)

	r := router.Group("/api")

	r.Use(logic.ApiAuth())
	{
		r.GET("/test", logic.ApiTest)
		r.POST("/test", logic.ApiTest)
	}
}
