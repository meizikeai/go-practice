package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SayHi(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}

func NotFound(ctx *gin.Context) {
	ctype := ctx.Request.Header.Get("Content-Type")

	if ctype == "application/json" {
		ctx.JSON(http.StatusNotFound, newResponse(http.StatusNotFound, nil))
	} else {
		ctx.JSON(http.StatusNotFound, "Not Found")
	}
}
