package controllers

import (
	"net/http"

	"github.com/mileusna/useragent"
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

func getUA(req *http.Request) string {
	return req.UserAgent()
}

func getUserAgent(data string) (os, osVersion, device string) {
	ua := useragent.Parse(data)

	return ua.OS, ua.OSVersion, ua.Device
}
