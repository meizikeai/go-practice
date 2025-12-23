// internal/pkg/ginctx/common.go
package ginctx

import (
	"net/http"
)

func GetReqID(req *http.Request) string {
	return req.Header.Get("X-Request-Id")
}

func GetClientIP(req *http.Request) string {
	ip := req.Header.Get("x-real-ip")

	if len(ip) == 0 {
		ip = req.Header.Get("x-forwarded-for")
	}

	return ip
}

func GetRemoteIP(req *http.Request) string {
	return req.RemoteAddr
}

func GetUserAgent(req *http.Request) string {
	return req.UserAgent()
}

func GetHeader(req *http.Request) map[string]string {
	ip := GetClientIP(req)
	id := GetReqID(req)

	result := map[string]string{
		"id": id,
		"ip": ip,
		"rp": req.RemoteAddr,
		"ua": req.UserAgent(),
	}

	return result
}
