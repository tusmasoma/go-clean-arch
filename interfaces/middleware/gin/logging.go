package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type loggingResponseWriter struct {
	gin.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		lrw := &loggingResponseWriter{ResponseWriter: c.Writer}
		c.Writer = lrw

		c.Next()

		log.Info(
			"Access log",
			log.Ftime("Date", time.Now()),
			log.Fstring("URL", c.Request.URL.String()),
			log.Fstring("IP", c.ClientIP()),
			log.Fint("StatusCode", lrw.statusCode),
		)

		if lrw.statusCode >= http.StatusBadRequest {
			log.Error(
				"Error log",
				log.Ftime("Date", time.Now()),
				log.Fstring("URL", c.Request.URL.String()),
				log.Fint("StatusCode", lrw.statusCode),
			)
		}
	}
}
