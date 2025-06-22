package main

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIpFromRequest(r *http.Request) (clientIp string) {
	clientIp = r.Header.Get("X-Original-Forwarded-For")

	if clientIp == "" {
		clientIp = r.Header.Get("X-Real-IP")
	}

	if clientIp == "" {
		clientIp = r.Header.Get("X-Forwarded-For")
	}

	if clientIp == "" {
		clientIp, _, err := net.SplitHostPort(r.RemoteAddr)

		if err != nil {
			return ""
		}

		return clientIp
	}

	clientIpChunks := strings.Split(clientIp, ",")

	return strings.Trim(clientIpChunks[len(clientIpChunks)-1], " ")
}

func main() {
	hostAddress := "0.0.0.0:6868"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		clientIp := GetClientIpFromRequest(r)
		w.Write([]byte(clientIp))
	})
	http.ListenAndServe(hostAddress, nil)
}
