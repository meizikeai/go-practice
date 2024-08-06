package routers

import (
	"go-practice/libs/utils"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func AddDefaultRouter(router *gin.Engine) {
	router.NoRoute(logic.NotFound)
	router.NoMethod(logic.NotFound)

	router.GET("/", logic.SayHi)
	router.GET("/favicon.ico", logic.SayHi)

	// kubernetes
	router.GET("/healthz", logic.SayHi)
	router.GET("/livez", logic.SayHi)
	router.GET("/readyz", logic.SayHi)

	// prometheus
	router.GET("/metrics", utils.PromHandler(promhttp.Handler()))
}
