// internal/pkg/ginctx/common.go
package ginctx

import (
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetBody(c *gin.Context) []byte {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return []byte{}
	}

	return body
}

func GetReqID(c *gin.Context) string {
	if reqID := c.GetHeader("X-Request-Id"); reqID != "" {
		return reqID
	}
	return ""
}

func GetClientIP(c *gin.Context) string {
	return c.ClientIP()
}

func GetRemoteIP(c *gin.Context) string {
	return c.RemoteIP()
}

func GetUserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}

func GetToken(c *gin.Context) string {
	headers := []string{"Authorization", "Token"}
	for _, h := range headers {
		if val := c.GetHeader(h); val != "" {
			if strings.HasPrefix(strings.ToLower(val), "bearer ") {
				return strings.TrimSpace(val[7:])
			}
			return val
		}
	}
	return c.Query("token")
}
