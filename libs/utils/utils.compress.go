package utils

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// compress
// not recommended for use in production.
type compressWriter struct {
	io.Writer
	gin.ResponseWriter
}

func (c *compressWriter) Write(b []byte) (int, error) {
	return c.Writer.Write(b)
}

var paths = []string{
	"/favicon.ico",
}
var shouldCompress = handlerExcludedPathRegexs(paths)

type excludedPathRegexs []*regexp.Regexp

func handlerExcludedPathRegexs(regexs []string) excludedPathRegexs {
	result := make([]*regexp.Regexp, len(regexs))

	for i, reg := range regexs {
		result[i] = regexp.MustCompile(reg)
	}

	return result
}

func (e excludedPathRegexs) Contains(requestURI string) bool {
	for _, reg := range e {
		if reg.MatchString(requestURI) {
			return true
		}
	}

	return false
}

func ContentEncoding() gin.HandlerFunc {
	return func(c *gin.Context) {
		if shouldCompress.Contains(c.Request.URL.Path) {
			return
		}

		accept := c.Request.Header.Get("Accept-Encoding")

		if strings.Contains(accept, "gzip") {
			// we've made sure compression level is valid in compress/gzip,
			// no need to check same error again.
			gz, err := gzip.NewWriterLevel(c.Writer, 3)

			if err == nil {
				c.Header("Content-Encoding", "gzip")
				c.Header("Vary", "Accept-Encoding")

				defer gz.Close()

				// delete content length after we know we have been written to
				c.Writer.Header().Del("Content-Length")

				c.Writer = &compressWriter{gz, c.Writer}
			}
		} else if strings.Contains(accept, "deflate") {
			// we've made sure compression level is valid in compress/flate,
			// no need to check same error again.
			de, err := flate.NewWriter(c.Writer, 3)

			if err == nil {
				c.Header("Content-Encoding", "deflate")
				c.Header("Vary", "Accept-Encoding")

				defer de.Close()

				// delete content length after we know we have been written to
				c.Writer.Header().Del("Content-Length")

				c.Writer = &compressWriter{de, c.Writer}
			}
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
