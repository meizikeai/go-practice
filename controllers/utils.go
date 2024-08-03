package controllers

import (
	"io"
	"net/http"

	"go-practice/libs/secret"
	"go-practice/libs/token"

	"github.com/gin-gonic/gin"
)

var jwt = token.NewJsonWebToken()
var chaos = secret.NewSecret()

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

var empty = map[string]any{}
var responseMessage = map[int]string{
	200: "OK",
	400: "Bad Request",
	401: "StatusUnauthorized",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	415: "Unsupported Media Type",
	422: "Unprocessable Entity",
	500: "Internal Server Error",

	403001: "select / inset / update / delete data failed!",
}

func newResponse(code int, data any) *response {
	message := ""

	if _, ok := responseMessage[code]; ok {
		message = responseMessage[code]
	}

	if data == nil {
		data = empty
	}

	return &response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func getBody(c *gin.Context) []byte {
	body, _ := io.ReadAll(c.Request.Body)

	// Get the value of body
	c.Set("body", body)

	return body
}

func getHeader(req *http.Request) map[string]string {
	ip := req.Header.Get("x-real-ip")

	if len(ip) == 0 {
		ip = req.Header.Get("x-forwarded-for")
	}

	// if len(ip) == 0 {
	// 	ip = req.RemoteAddr
	// }

	result := map[string]string{
		"ip":     ip,
		"rp":     req.RemoteAddr,
		"id":     req.Header.Get("x-request-id"),
		"uid":    req.Header.Get("x-remote-userid"),
		"ua":     req.UserAgent(),
		"uri":    req.RequestURI,
		"method": req.Method,
		"lang":   req.Header.Get("accept-language"),
	}

	return result
}
