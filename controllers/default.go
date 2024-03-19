package controllers

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func SayHi(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}

func NotFound(ctx *gin.Context) {
	ctype := ctx.Request.Header.Get("Content-Type")

	if ctype == "application/json" {
		ctx.JSON(http.StatusNotFound, newResponse(http.StatusNotFound, nil))
	} else {
		ctx.JSON(http.StatusNotFound, "Not Found")
	}
}

// prometheus
var pointLabel = HandleRouterCompile()

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

func EndpointLabelMappingFn(c *gin.Context) string {
	result := "/unknown"
	url := []byte(c.Request.URL.Path)

	for k, v := range pointLabel {
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
