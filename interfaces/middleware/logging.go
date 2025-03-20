package middleware

import (
	"net/http"
	"time"

	"github.com/tusmasoma/go-clean-arch/pkg/log"
	"github.com/tusmasoma/go-clean-arch/pkg/response"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := response.NewCustomResponseWriter(w)
		defer func() {
			// Access log
			log.Info(
				"Access log",
				log.Ftime("Date", time.Now()),
				log.Fstring("URL", r.URL.String()),
				log.Fstring("IP", r.RemoteAddr),
				log.Fint("StatusCode", ww.StatusCode()),
			)
			// Error log if status code is 400 or higher
			if ww.StatusCode() >= http.StatusBadRequest {
				log.Error(
					"Error log",
					log.Ftime("Date", time.Now()),
					log.Fstring("URL", r.URL.String()),
					log.Fint("StatusCode", ww.StatusCode()),
				)
			}
		}()
		next.ServeHTTP(ww, r)
	})
}
