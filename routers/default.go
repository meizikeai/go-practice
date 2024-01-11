package routers

import (
	"go-practice/controllers"

	"github.com/gin-gonic/gin"
)

func AddDefaultRouter(router *gin.Engine) {
	router.NoRoute(controllers.NotFound)
	router.NoMethod(controllers.NotFound)

	router.GET("/", controllers.SayHi)
	router.GET("/favicon.ico", controllers.SayHi)
}
