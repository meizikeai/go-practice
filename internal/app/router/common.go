// internal/app/router/common.go
package router

import (
	"go-practice/internal/app"
	"go-practice/internal/app/handler"
	"go-practice/internal/pkg/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Setup(app *app.App) {
	logic := handler.NewHandler(app)

	app.Engine.GET("/", logic.SayHi)
	app.Engine.NoRoute(logic.NoRoute)
	app.Engine.NoMethod(logic.NoMethod)
	app.Engine.GET("/metrics", prometheus.PromHandler(promhttp.Handler()))

	app.Engine.GET("/test/get", logic.TestGet)
	app.Engine.POST("/test/post", logic.TestPost)
}
