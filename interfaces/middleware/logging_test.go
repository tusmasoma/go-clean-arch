package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

func Test_Logging(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)

	tests := []struct {
		name            string
		handler         http.HandlerFunc
		expectedStatus  int
		expectAccessLog bool
		expectErrorLog  bool
	}{
		{
			name: "successful request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK")) //nolint:errcheck // ignore error
			},
			expectedStatus:  http.StatusOK,
			expectAccessLog: true,
			expectErrorLog:  false,
		},
		{
			name: "client error request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Bad Request")) //nolint:errcheck // ignore error
			},
			expectedStatus:  http.StatusBadRequest,
			expectAccessLog: true,
			expectErrorLog:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			handlerToTest := Logging(tt.handler)

			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			w := httptest.NewRecorder()

			handlerToTest.ServeHTTP(w, req)

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
