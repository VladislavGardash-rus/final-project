package middleware

import (
	"fmt"
	"github.com/gardashvs/final-project/internal/logger"
	"net/http"
	"time"
)

func Logging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		handler.ServeHTTP(w, r)

		logger.UseLogger().Info(
			fmt.Sprintf("%s %s %s %s %d %d %s",
				r.RemoteAddr,
				startTime.Format(time.DateTime),
				r.Method,
				r.RequestURI,
				r.Response.StatusCode,
				time.Since(startTime).Milliseconds(),
				r.Header["User-Agent"],
			))
	})
}
