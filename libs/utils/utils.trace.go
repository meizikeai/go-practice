package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	Uri     string `json:"uri"`
	Method  string `json:"method"`
	Status  int    `json:"status"`
	Client  string `json:"client"`
	Request string `json:"request"`
	Body    any    `json:"body,omitempty"`
	Data    any    `json:"data,omitempty"`
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

		body := tool.CleanSpace(string(getMountBody(c)))

		data := writer.body.String()

		trace := traceLog{
			Uri:    uri,
			Method: method,
			Status: status,
			Client: client,
			Body:   tool.UnmarshalJson(body),
			Data:   tool.UnmarshalJson(data),
		}

		log.Trace(fmt.Sprintf("%s %s", string(tool.MarshalJson(trace)), latency))
	}
}

func LoggingIllegalEntity(c *gin.Context) {
	body := tool.CleanSpace(string(getMountBody(c)))

	res := traceLog{
		Uri:     c.Request.RequestURI,
		Method:  c.Request.Method,
		Status:  c.Writer.Status(),
		Client:  c.ClientIP(),
		Request: getRequestID(c.Request),
		Body:    tool.UnmarshalJson(body),
		Data:    nil,
	}

	data, _ := json.Marshal(res)

	log.Info(fmt.Sprintf("%s %s", "LoggingIllegalEntity", string(data)))
}

func getMountBody(ctx *gin.Context) []byte {
	d, _ := ctx.Get("body")
	result, _ := d.([]byte)

	return result
}

func getRequestID(req *http.Request) string {
	rid := req.Header.Get("http_x_request_id")

	if len(rid) != 0 {
		return rid
	}

	return ""
}
