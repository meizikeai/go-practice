// internal/app/controller/base.go
package controller

import (
	"go-practice/internal/app"
	"go-practice/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
	Resp *response.Responder
}

func (b *BaseController) Bind(c *gin.Context) *response.Responder {
	b.Resp = response.NewResponder(c)
	return b.Resp
}

type Controller struct {
	*BaseController
	app *app.App
}

func NewController(app *app.App) *Controller {
	return &Controller{
		app:            app,
		BaseController: &BaseController{},
	}
}

func (l *Controller) SayHi(ctx *gin.Context) {
	l.Bind(ctx).Success(ctx, nil)
}

func (l *Controller) NoRoute(ctx *gin.Context) {
	l.Bind(ctx).Fail(ctx, response.CodeNotFound)
}

func (l *Controller) NoMethod(ctx *gin.Context) {
	l.Bind(ctx).Fail(ctx, response.CodeMethodNotAllowed)
}
