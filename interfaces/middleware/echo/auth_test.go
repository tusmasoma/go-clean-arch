package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/tusmasoma/go-clean-arch/config"

	"github.com/tusmasoma/go-clean-arch/repository/mock"
)

func dummyTestHandler(c echo.Context) error {
	ctx := c.Request().Context()
	userIDValue, exists := ctx.Value(config.ContextUserIDKey).(string)
	if !exists || userIDValue == "" {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "userID not found"})
	}
	return c.NoContent(http.StatusOK)
}

func TestAuthMiddleware_Authenticate(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	email := "test@gmail.com"
	jwt := "eyJhbGciOiJIUzI1NiIsI.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0Ijo.SflKxwRJSMeKKF2QT4fwpMeJf36P"

	patterns := []struct {
		name       string
		setup      func(m *mock.MockAuthRepository)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockAuthRepository) {
				m.EXPECT().ValidateAccessToken(jwt).Return(nil)
				m.EXPECT().GetPayloadFromToken(jwt).Return(
					map[string]string{
						"userId": userID,
						"email":  email,
					}, nil,
				)
			},
			in: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer "+jwt)
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: No Auth Header",
			in: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				return req
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail: Invalid Auth Header Format",
			in: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", jwt)
				return req
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail: Invalid Token",
			setup: func(m *mock.MockAuthRepository) {
				m.EXPECT().ValidateAccessToken("invalidToken").Return(errors.New("invalid token"))
			},
			in: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer invalidToken")
				return req
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ar := mock.NewMockAuthRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ar)
			}

			am := NewAuthMiddleware(ar)

			e := echo.New()
			e.Use(am.Authenticate)

			e.GET("/", dummyTestHandler)

			req := tt.in()
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			if status := rec.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}
