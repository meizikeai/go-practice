// internal/app/handler/base.go
package handler

import (
	"context"

	"go-practice/internal/app"
	"go-practice/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

type BaseHandler struct {
	Resp *response.Responder
}

func (b *BaseHandler) Bind(c *gin.Context) *response.Responder {
	b.Resp = response.NewResponder(c)
	return b.Resp
}

type Handler struct {
	*BaseHandler
	app *app.App
}

func New(app *app.App) *Handler {
	return &Handler{
		app:         app,
		BaseHandler: &BaseHandler{},
	}
}

func (h *Handler) SayHi(c *gin.Context) {
	h.Bind(c).Success(c, nil)
}

func (h *Handler) NoRoute(c *gin.Context) {
	h.Bind(c).Fail(c, response.CodeNotFound)
}

func (h *Handler) NoMethod(c *gin.Context) {
	h.Bind(c).Fail(c, response.CodeMethodNotAllowed)
}
