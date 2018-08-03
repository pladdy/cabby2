package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// WithAcceptType decorates a handle with content type check
func WithAcceptType(h http.HandlerFunc, typeToCheck string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType, _ := splitAcceptHeader(r.Header.Get("Accept"))

		if contentType != typeToCheck {
			unsupportedMediaType(w, fmt.Errorf("Invalid 'Accept' Header: %v", contentType))
			return
		}
		h(w, r)
	}
}

// func withBasicAuth(h http.Handler, us cabby.UserService) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		u, p, ok := r.BasicAuth()
// 		tu, _ := us.Read(u, p)
//
// 		if !ok || !tu.valid() {
// 			unauthorized(w, errors.New("Invalid user/pass combination"))
// 			return
// 		}
//
// 		log.WithFields(log.Fields{"user": u}).Info("User authenticated")
// 		h.ServeHTTP(withHSTS(w), withTaxiiUser(r, tu))
// 	})
// }

func withRequestLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		milliSecondOfNanoSeconds := int64(1000000)

		start := time.Now().In(time.UTC)
		log.WithFields(log.Fields{
			"method":   r.Method,
			"start_ts": start.UnixNano() / milliSecondOfNanoSeconds,
			"url":      r.URL,
		}).Info("Serving request made to server")

		h.ServeHTTP(w, r)

		end := time.Now().In(time.UTC)
		elapsed := time.Since(start)

		log.WithFields(log.Fields{
			"elapsed_ts": float64(elapsed.Nanoseconds()) / float64(milliSecondOfNanoSeconds),
			"method":     r.Method,
			"end_ts":     end.UnixNano() / milliSecondOfNanoSeconds,
			"url":        r.URL,
		}).Info("Served request made to server")
	})
}

/* helpers */

func splitAcceptHeader(h string) (string, string) {
	parts := strings.Split(h, ";")
	first := parts[0]

	var second string
	if len(parts) > 1 {
		second = parts[1]
	}

	return first, second
}
