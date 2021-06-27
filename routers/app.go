package routers

import (
	"go-practice/controllers"

	"github.com/gin-gonic/gin"
)

func AddWebRouter(router *gin.Engine) {
	router.NoRoute(controllers.NotFound)
	router.NoMethod(controllers.NotFound)

	router.GET("/", controllers.Home)
}
