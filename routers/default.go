package routers

import (
	"go-practice/controllers"
	"go-practice/libs/utils"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func AddDefaultRouter(router *gin.Engine) {
	router.NoRoute(controllers.NotFound)
	router.NoMethod(controllers.NotFound)

	router.GET("/", controllers.SayHi)
	router.GET("/favicon.ico", controllers.SayHi)

	// kubernetes
	router.GET("/healthz", controllers.SayHi)
	router.GET("/livez", controllers.SayHi)
	router.GET("/readyz", controllers.SayHi)

	// prometheus
	router.GET("/metrics", utils.PromHandler(promhttp.Handler()))
}
