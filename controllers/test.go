package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (l *Logic) Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "OK",
		"data":    gin.H{},
	})
}

func (l *Logic) ApiTest(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, tiger.newResponse(http.StatusOK, gin.H{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
		"f": 6,
	}))
}
