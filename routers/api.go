package routers

import (
	"go-practice/controllers"
	"go-practice/libs/jwt"

	"github.com/gin-gonic/gin"
)

func AddApiRouter(router *gin.Engine) {
	router.NoRoute(controllers.NotFound)
	router.NoMethod(controllers.NotFound)

	router.GET("/", controllers.SayHi)
	router.GET("/favicon.ico", controllers.SayHi)

	r := router.Group("/api")

	r.Use(jwt.ApiAuth())
	{
		r.GET("/test", controllers.ApiTest)
		r.GET("/add/person", controllers.ApiAddPerson)
	}
}
