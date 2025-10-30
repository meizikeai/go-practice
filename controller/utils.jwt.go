package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (l *Logic) ApiAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")

		if strings.TrimSpace(authorization) == "" {
			l.forbidden(c)
			return
		}

		authorizationList := strings.Split(authorization, " ")

		if len(authorizationList) != 2 {
			l.forbidden(c)
			return
		}

		bearer := authorizationList[0]
		token := authorizationList[1]

		if strings.ToLower(bearer) != "bearer" || strings.TrimSpace(token) == "" {
			l.forbidden(c)
			return
		}

		// decode token
		token = chaos.HandleServiceDecrypt(token)

		// jwt
		claims, err := jwt.DecryptToken(token)

		if err != nil {
			l.forbidden(c)
			return
		}

		// After passing the authentication, get the value of body
		body := tiger.getBody(c)

		c.Set("body", body)
		c.Set("claims", claims)
	}
}

func (l *Logic) forbidden(c *gin.Context) {
	ctype := c.Request.Header.Get("Content-Type")

	if ctype == "application/json" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  403,
			"message": "Forbidden",
		})
	} else {
		c.Abort()
		c.String(http.StatusForbidden, "Forbidden")
	}
}
