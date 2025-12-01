// internal/pkg/ginctx/ginctx.go
package ginctx

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

type Context struct {
	c *gin.Context
}

func New(c *gin.Context) *Context {
	return &Context{c: c}
}

func (g *Context) Raw() *gin.Context {
	return g.c
}

func (g *Context) Context() context.Context {
	return g.c.Request.Context()
}

func (g *Context) GetReqID() string {
	if reqID := g.c.GetHeader("X-Request-Id"); reqID != "" {
		return reqID
	}
	// reqID := uuid.New().String()
	// g.c.Header("X-Request-Id", reqID)
	// return reqID
	return ""
}

func (g *Context) GetClientIP() string {
	return g.c.ClientIP()
}

func (g *Context) GetRemoteIP() string {
	return g.c.RemoteIP()
}

func (g *Context) GetUserAgent() string {
	return g.c.GetHeader("User-Agent")
}

func (g *Context) GetToken() string {
	headers := []string{"Authorization", "Token"}
	for _, h := range headers {
		if val := g.c.GetHeader(h); val != "" {
			if strings.HasPrefix(strings.ToLower(val), "bearer ") {
				return strings.TrimSpace(val[7:])
			}
			return val
		}
	}
	return g.c.Query("token")
}

func FromContext(ctx context.Context) (*Context, bool) {
	c, ok := ctx.(*gin.Context)
	if !ok {
		return nil, false
	}
	return New(c), true
}

// // middleware
// func RequestID() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		reqID := c.GetHeader("X-Request-Id")

// 		if reqID == "" {
// 			reqID = uuid.New().String()
// 		}

// 		c.Header("X-Request-Id", reqID)
// 		c.Next()
// 	}
// }
