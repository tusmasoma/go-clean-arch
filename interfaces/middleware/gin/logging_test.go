package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

func Test_Logging(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)

	tests := []struct {
		name            string
		handler         gin.HandlerFunc
		expectedStatus  int
		expectAccessLog bool
		expectErrorLog  bool
	}{
		{
			name: "successful request",
			handler: func(c *gin.Context) {
				c.String(http.StatusOK, "OK")
			},
			expectedStatus:  http.StatusOK,
			expectAccessLog: true,
			expectErrorLog:  false,
		},
		{
			name: "client error request",
			handler: func(c *gin.Context) {
				c.String(http.StatusBadRequest, "Bad Request")
			},
			expectedStatus:  http.StatusBadRequest,
			expectAccessLog: true,
			expectErrorLog:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Ginのルーターを作成
			router := gin.New()
			router.Use(Logging())
			router.GET("/foo", tt.handler)

			// リクエストを作成して実行
			req := httptest.NewRequest(http.MethodGet, "/foo", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// ステータスコードのチェック
			if status := w.Result().StatusCode; status != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, status)
			}

			// ログのチェック
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
