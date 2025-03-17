package middleware

import (
	"net/http"
	"time"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w}

		next.ServeHTTP(lrw, r)
		// Access log
		log.Info(
			"Access log",
			log.Ftime("Date", time.Now()),
			log.Fstring("URL", r.URL.String()),
			log.Fstring("IP", r.RemoteAddr),
			log.Fint("StatusCode", lrw.statusCode),
		)

		// Error log if status code is 400 or higher
		if lrw.statusCode >= http.StatusBadRequest {
			log.Error(
				"Error log",
				log.Ftime("Date", time.Now()),
				log.Fstring("URL", r.URL.String()),
				log.Fint("StatusCode", lrw.statusCode),
			)
		}
	})
}
