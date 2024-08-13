package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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

func Logging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lrw := &loggingResponseWriter{ResponseWriter: c.Response().Writer}
		c.Response().Writer = lrw

		err := next(c)

		log.Info(
			"Access log",
			log.Ftime("Date", time.Now()),
			log.Fstring("URL", c.Request().URL.String()),
			log.Fstring("IP", c.RealIP()),
			log.Fint("StatusCode", lrw.statusCode),
		)

		if lrw.statusCode >= http.StatusBadRequest {
			log.Error(
				"Error log",
				log.Ftime("Date", time.Now()),
				log.Fstring("URL", c.Request().URL.String()),
				log.Fint("StatusCode", lrw.statusCode),
			)
		}

		return err
	}
}
