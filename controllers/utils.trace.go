package controllers

import (
	"bytes"
	"fmt"
	"time"

	"go-practice/libs/tool"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// trace logger
type traceWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (c traceWriter) Write(b []byte) (int, error) {
	c.body.Write(b)
	return c.ResponseWriter.Write(b)
}

type traceLog struct {
	Uri    string `json:"uri"`
	Method string `json:"method"`
	Status int    `json:"status"`
	Client string `json:"client"`
	Body   any    `json:"body"`
	Data   any    `json:"data"`
}

func TraceLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		writer := &traceWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}

		c.Writer = writer

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		client := c.ClientIP()
		status := c.Writer.Status()
		method := c.Request.Method
		uri := c.Request.RequestURI

		body := tool.ClearSpace(string(getMountBody(c)))

		data := writer.body.String()

		trace := traceLog{
			Uri:    uri,
			Method: method,
			Status: status,
			Client: client,
			Body:   tool.UnmarshalJson(body),
			Data:   tool.UnmarshalJson(data),
		}

		log.Trace(fmt.Sprintf("%s %s %s", tool.GetTime(), string(tool.MarshalJson(trace)), latency))
	}
}

func getMountBody(ctx *gin.Context) []byte {
	d, _ := ctx.Get("body")
	result, _ := d.([]byte)

	return result
}
