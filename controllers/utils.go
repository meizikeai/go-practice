package controllers

import (
	"io"
	"net/http"
	"regexp"

	"go-practice/libs/secret"
	"go-practice/libs/token"

	"github.com/gin-gonic/gin"
)

type Logic struct {
	pointLabel map[string]*regexp.Regexp
}

func NewLogic() *Logic {
	return &Logic{
		pointLabel: HandleRouterCompile(),
	}
}

// prometheus
func HandleRouterCompile() map[string]*regexp.Regexp {
	var data = map[string]string{
		"metrics":     "metrics$",
		"healthz":     "healthz$",
		"favicon.ico": "favicon.ico$",
	}

	result := map[string]*regexp.Regexp{}

	for k, v := range data {
		re, err := regexp.Compile(v)

		if err != nil {
			continue
		}

		result[k] = re
	}

	return result
}

func (l *Logic) EndpointLabelMappingFn(c *gin.Context) string {
	result := "/unknown"
	url := []byte(c.Request.URL.Path)

	for k, v := range l.pointLabel {
		match := v.Match(url)

		if match == true {
			result = k
			break
		}
	}

	if c.Writer.Status() == 404 {
		result = "/unknown"
	}

	return result
}

type Tiger struct {
	empty           map[string]any
	responseMessage map[int]string
}

func NewTiger() *Tiger {
	return &Tiger{
		empty: map[string]any{},
		responseMessage: map[int]string{
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
		},
	}
}

var tiger = NewTiger()
var jwt = token.NewJsonWebToken()
var chaos = secret.NewSecret()

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (t *Tiger) newResponse(code int, data any) *response {
	message := ""

	if _, ok := t.responseMessage[code]; ok {
		message = t.responseMessage[code]
	}

	if data == nil {
		data = t.empty
	}

	return &response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func (t *Tiger) getBody(c *gin.Context) []byte {
	body, _ := io.ReadAll(c.Request.Body)

	// Get the value of body
	c.Set("body", body)

	return body
}

func (t *Tiger) getHeader(req *http.Request) map[string]string {
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
