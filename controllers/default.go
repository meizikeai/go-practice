package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (l *Logic) SayHi(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}

func (l *Logic) NotFound(ctx *gin.Context) {
	ctype := ctx.Request.Header.Get("Content-Type")

	if ctype == "application/json" {
		ctx.JSON(http.StatusNotFound, tiger.newResponse(http.StatusNotFound, nil))
	} else {
		ctx.JSON(http.StatusNotFound, "Not Found")
	}
}
