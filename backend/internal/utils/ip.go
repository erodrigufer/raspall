package utils

import (
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
	var ipReqChain string

	forwardedHeader := r.Header.Values("X-Forwarded-For")
	if len(forwardedHeader) > 0 {
		ipReqChain = strings.Join(forwardedHeader, ", ")
	} else {
		ipReqChain = r.RemoteAddr
	}

	return ipReqChain
}
