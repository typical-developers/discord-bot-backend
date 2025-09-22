package handlers

import (
	"fmt"
	"net/http"
	"time"
)

type ResponseStatus struct {
	http.ResponseWriter
	StatusCode int
}

func (r *ResponseStatus) WriteHeader(code int) {
	r.StatusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func RequestLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rs := &ResponseStatus{ResponseWriter: w}
		handler.ServeHTTP(rs, r)

		end := time.Now()
		elapsed := end.Sub(start)

		fmt.Printf("%s | %-4d | %-6s | %-10s | %-15s | %s\n",
			start.Format(time.RFC3339), rs.StatusCode, r.Method, elapsed, r.RemoteAddr, r.URL,
		)
	})
}
