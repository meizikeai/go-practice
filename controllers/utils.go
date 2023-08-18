package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"go-practice/libs/types"

	"github.com/gin-gonic/gin"
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

func getBody(req *http.Request) types.RequestBody {
	data := types.RequestBody{}
	body, err := io.ReadAll(req.Body)

	if err != nil {
		return data
	}

	err = json.Unmarshal(body, &data)

	if err != nil {
		return data
	}

	return data
}

func ApiAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.Request.Header.Get("Authorization")
		authorizationList := strings.Split(authorization, " ")

		basic := authorizationList[0]
		accessToken := authorizationList[1]

		if basic != "Basic" {
			ctx.AbortWithStatusJSON(http.StatusOK, newResponse(http.StatusUnprocessableEntity, nil))
			return
		}

		if strings.TrimSpace(accessToken) == "" {
			ctx.AbortWithStatusJSON(http.StatusOK, newResponse(http.StatusUnprocessableEntity, nil))
			return
		}

		body := getBody(ctx.Request)

		// todo
		// 此处进行鉴权相关操作
		state := 0

		if state == 0 {
			ctx.AbortWithStatusJSON(http.StatusOK, newResponse(http.StatusOK, nil))
			return
		}

		// 往下传递 body
		ctx.Set("body", body)
	}
}
