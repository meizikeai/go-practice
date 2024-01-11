package controllers

import (
	"io"
	"net/http"
	"strings"

	"go-practice/libs/jwt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func getBody(req *http.Request) []byte {
	data, _ := io.ReadAll(req.Body)
	return data
}

func forbidden(c *gin.Context) {
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

func ApiAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")

		if strings.TrimSpace(authorization) == "" {
			forbidden(c)
			return
		}

		authorizationList := strings.Split(authorization, " ")

		if len(authorizationList) != 2 {
			forbidden(c)
			return
		}

		bearer := authorizationList[0]
		token := authorizationList[1]

		if strings.ToLower(bearer) != "bearer" || strings.TrimSpace(token) == "" {
			forbidden(c)
			return
		}
		// log.Print("get token: ", token)

		// jwt
		j := jwt.NewJWT()
		claims, err := j.DecryptToken(token)

		if err != nil {
			if err == jwt.TokenExpired {
				log.Error("Token is expired")

				forbidden(c)
				return
			}

			log.Error(err)

			forbidden(c)
			return
		}

		// After passing the authentication, get the value of body
		body := getBody(c.Request)

		c.Set("body", body)
		c.Set("claims", claims)
	}
}
