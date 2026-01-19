// internal/pkg/response/common.go
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	TraceID string `json:"trace_id,omitempty"`
}

type Code int

const (
	CodeOK                  Code = 200
	CodeBadRequest          Code = 400
	CodeUnauthorized        Code = 401
	CodeForbidden           Code = 403
	CodeNotFound            Code = 404
	CodeMethodNotAllowed    Code = 405
	CodeUnprocessableEntity Code = 422
	CodeInternalServerError Code = 500
	CodeServiceUnavailable  Code = 503

	// Custom code
	CodeUserNotFound Code = 100001
	CodeDBError      Code = 100002
)

var codeMsg = map[Code]string{
	CodeOK:                  "OK",
	CodeBadRequest:          "Bad Request",
	CodeUnauthorized:        "Unauthorized",
	CodeForbidden:           "Forbidden",
	CodeNotFound:            "Not Found",
	CodeMethodNotAllowed:    "Method Not Allowed",
	CodeUnprocessableEntity: "Unprocessable Entity",
	CodeInternalServerError: "Internal Server Error",
	CodeServiceUnavailable:  "Service Unavailable",

	// Custom code
	CodeUserNotFound: "User Not Found",
	CodeDBError:      "Database Error",
}

type Responder struct{}

func NewResponder(c *gin.Context) *Responder {
	return &Responder{}
}

func (r *Responder) JSON(c *gin.Context, httpStatus int, code Code, data any, overrides ...string) {
	msg := codeMsg[code]
	if len(overrides) > 0 {
		msg = overrides[0]
	}

	c.JSON(httpStatus, Response{
		Code:    int(code),
		Message: msg,
		Data:    data,
	})
}

func (r *Responder) Success(c *gin.Context, data any) {
	r.JSON(c, http.StatusOK, CodeOK, data)
}

func (r *Responder) Created(c *gin.Context, data any) {
	r.JSON(c, http.StatusCreated, CodeOK, data)
}

func (r *Responder) Fail(c *gin.Context, code Code, overrides ...string) {
	httpStatus := codeToHTTPStatus(code)
	r.JSON(c, httpStatus, code, nil, overrides...)
}

func (r *Responder) Error(c *gin.Context, err error) {
	r.JSON(c, http.StatusInternalServerError, CodeInternalServerError, nil, err.Error())
}

func codeToHTTPStatus(code Code) int {
	switch code {
	case 400:
		return http.StatusBadRequest
	case 401:
		return http.StatusUnauthorized
	case 403:
		return http.StatusForbidden
	case 404:
		return http.StatusNotFound
	case 405:
		return http.StatusMethodNotAllowed
	case 422:
		return http.StatusUnprocessableEntity
	case 500:
		return http.StatusInternalServerError
	case 503:
		return http.StatusServiceUnavailable
	default:
		if code != CodeOK {
			return http.StatusInternalServerError
		}
	}
	return http.StatusOK
}
