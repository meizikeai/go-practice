package routers

import (
	"go-practice/controllers"

	"github.com/gin-gonic/gin"
)

func AddWebRouter(router *gin.Engine) {
	// NoRoute
	router.NoRoute(controllers.NotFound)
	// NoMethod
	router.NoMethod(controllers.NotFound)

	// home
	router.GET("/", controllers.Home)
}
