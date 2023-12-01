package controllers

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go-practice/libs/jwt"
	"go-practice/libs/tool"

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

type compressWriter struct {
	io.Writer
	gin.ResponseWriter
}

func (c *compressWriter) Write(b []byte) (int, error) {
	return c.Writer.Write(b)
}

func ContentEncoding() gin.HandlerFunc {
	return func(c *gin.Context) {
		accept := c.Request.Header.Get("Accept-Encoding")

		if strings.Contains(accept, "gzip") {
			// we've made sure compression level is valid in compress/gzip,
			// no need to check same error again.
			gz, err := gzip.NewWriterLevel(c.Writer, 3)

			if err != nil {
				panic(err.Error())
			}

			c.Header("Content-Encoding", "gzip")
			c.Header("Vary", "Accept-Encoding")

			defer gz.Close()

			// delete content length after we know we have been written to
			c.Writer.Header().Del("Content-Length")

			c.Writer = &compressWriter{gz, c.Writer}
		} else if strings.Contains(accept, "deflate") {
			// we've made sure compression level is valid in compress/flate,
			// no need to check same error again.
			de, err := flate.NewWriter(c.Writer, 3)

			if err != nil {
				panic(err.Error())
			}

			c.Header("Content-Encoding", "deflate")
			c.Header("Vary", "Accept-Encoding")

			defer de.Close()

			// delete content length after we know we have been written to
			c.Writer.Header().Del("Content-Length")

			c.Writer = &compressWriter{de, c.Writer}
		}

		// golang does not support Brotli
		// else if strings.Contains(accept, "br") {
		// 	fmt.Println("br")
		// }

		// If the request header does not contain gzip / deflate,
		// continue processing the request.
		c.Next()
	}
}

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
