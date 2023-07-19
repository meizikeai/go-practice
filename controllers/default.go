package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

func TraceLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		clientIP := c.ClientIP()
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()

		data := fmt.Sprintf("client:%s, latency:%s, status:%v, method:%s, uri:%s",
			clientIP, latencyTime, statusCode, reqMethod, reqUri)

		log.Trace(data)
	}
}
