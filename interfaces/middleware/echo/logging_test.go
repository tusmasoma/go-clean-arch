package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

func Test_Logging(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)

	tests := []struct {
		name            string
		handler         echo.HandlerFunc
		expectedStatus  int
		expectAccessLog bool
		expectErrorLog  bool
	}{
		{
			name: "successful request",
			handler: func(c echo.Context) error {
				return c.String(http.StatusOK, "OK")
			},
			expectedStatus:  http.StatusOK,
			expectAccessLog: true,
			expectErrorLog:  false,
		},
		{
			name: "client error request",
			handler: func(c echo.Context) error {
				return c.String(http.StatusBadRequest, "Bad Request")
			},
			expectedStatus:  http.StatusBadRequest,
			expectAccessLog: true,
			expectErrorLog:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Use(Logging)
			e.GET("/foo", tt.handler)

			req := httptest.NewRequest(http.MethodGet, "/foo", nil)
			w := httptest.NewRecorder()

			e.ServeHTTP(w, req)

			if status := w.Result().StatusCode; status != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, status)
			}

			logOutput := logBuf.String()
			if tt.expectAccessLog && !strings.Contains(logOutput, "Access log") {
				t.Errorf("expected access log, got %s", logOutput)
			}

			if tt.expectErrorLog && !strings.Contains(logOutput, "Error log") {
				t.Errorf("expected error log, got %s", logOutput)
			}

			logBuf.Reset()
		})
	}
}
