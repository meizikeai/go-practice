package controllers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var empty = map[string]interface{}{}
var responseMessage = map[int]string{
	200: "OK",
	400: "Bad Request",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	415: "Unsupported Media Type",
	422: "Unprocessable Entity",
	500: "Internal Server Error",
}

func newResponse(code int, data interface{}) *response {
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
	result := map[string]string{}

	// ip
	ip := req.Header.Get("x-real-ip")

	if len(ip) == 0 {
		ip = req.Header.Get("x-forwarded-for")
	}

	if len(ip) == 0 {
		ip = req.RemoteAddr
	}

	result["ip"] = ip
	result["rp"] = req.RemoteAddr

	// id
	id := req.Header.Get("x-request-id")

	if len(id) <= 0 {
		id = ""
	}

	result["id"] = id

	// uid
	uid := req.Header.Get("x-remote-userid")

	result["uid"] = uid

	// ua
	result["ua"] = req.UserAgent()

	return result
}
