// internal/pkg/log/gin.zap.go
package log

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	maxTruncateSize = 4 * 1024
	skipBodyRoutes  = "/healthz,/metrics,/favicon.ico"
)

var (
	jsonContentTypes = map[string]bool{
		"application/json":                true,
		"application/json; charset=utf-8": true,
		"application/json;charset=utf-8":  true,
	}
	sensitiveFields  = regexp.MustCompile(`(?i)"(?:password|passwd|id_card|bank_card)":"[^"]*"`)
	skipBodyRouteMap = make(map[string]bool)
)

func init() {
	for _, r := range strings.Split(skipBodyRoutes, ",") {
		if r = strings.TrimSpace(r); r != "" {
			skipBodyRouteMap[r] = true
		}
	}
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqBody := []byte{}
		if shouldLogRequestBody(c) {
			if c.Request.Body != nil {
				reqBody, _ = io.ReadAll(io.LimitReader(c.Request.Body, int64(maxTruncateSize+1024)))
				c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}
		}

		lrw := &loggingResponseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = lrw

		c.Next()

		latency := time.Since(start)

		fields := make([]any, 0)
		fields = append(fields,
			zap.Duration("latency", latency),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("req_size", c.Request.ContentLength),
			zap.Int("resp_size", lrw.body.Len()),
			zap.String("client_ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("req_id", c.GetHeader("X-Request-Id")),
			zap.String("remote_ip", c.RemoteIP()),
			zap.String("route", getRoute(c)),
			zap.String("user_agent", c.GetHeader("User-Agent")),
		)

		if len(reqBody) > 0 {
			body := desensitizeAndCompact(reqBody)
			if len(body) > maxTruncateSize {
				body = append(body[:maxTruncateSize], []byte("...[truncated]")...)
			}
			fields = append(fields, zap.ByteString("req_body", body))
		}

		if shouldLogResponseBody(c, lrw) {
			respBytes := lrw.body.Bytes()
			respStr := desensitizeAndCompact(respBytes)
			if len(respStr) > maxTruncateSize {
				respStr = append(respStr[:maxTruncateSize], []byte("...[truncated]")...)
			}
			fields = append(fields, zap.ByteString("resp_body", respStr))
		} else if lrw.body.Len() > maxTruncateSize {
			fields = append(fields, zap.String("resp_body", "[skipped: too large]"))
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		if c.Writer.Status() >= 500 {
			Sugar.Errorw("access", fields...)
		} else if c.Writer.Status() >= 400 {
			Sugar.Warnw("access", fields...)
		} else {
			Sugar.Infow("access", fields...)
		}
	}
}

type loggingResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func shouldLogRequestBody(c *gin.Context) bool {
	if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
		return false
	}
	if skipBodyRouteMap[c.FullPath()] {
		return false
	}
	if c.Request.ContentLength > maxTruncateSize*2 {
		return false
	}

	ctype := c.Request.Header.Get("Content-Type")
	return strings.HasPrefix(ctype, "application/json") ||
		strings.HasPrefix(ctype, "application/x-www-form-urlencoded")
}

func shouldLogResponseBody(c *gin.Context, lrw *loggingResponseWriter) bool {
	if c.Writer.Status() >= 400 {
		return true
	}

	ctype := c.Writer.Header().Get("Content-Type")
	if !jsonContentTypes[ctype] && !strings.HasPrefix(ctype, "application/json") {
		return false
	}

	return lrw.body.Len() <= maxTruncateSize
}

func desensitizeAndCompact(b []byte) []byte {
	if len(b) == 0 {
		return b
	}

	s := sensitiveFields.ReplaceAllString(string(b), `"$1":"***"`)

	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(s)); err != nil {
		return []byte(s)
	}

	return buf.Bytes()
}

func getRoute(c *gin.Context) string {
	if route := c.FullPath(); route != "" {
		return route
	}
	return c.Request.URL.Path
}
