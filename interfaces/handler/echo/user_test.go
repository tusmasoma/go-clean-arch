package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/usecase/mock"
)

func TestUserHandler_GetUser(t *testing.T) {
	t.Parallel()

	user := entity.User{
		ID:    uuid.New().String(),
		Name:  "test",
		Email: "test@gmail.com",
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserUseCase) {
				m.EXPECT().GetUser(gomock.Any()).Return(&user, nil)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/api/user/get", nil)
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: No User ID in context",
			setup: func(m *mock.MockUserUseCase) {
				m.EXPECT().GetUser(gomock.Any()).Return(
					nil,
					fmt.Errorf("user name not found in request context"),
				)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/api/user/get", nil)
				return req
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			uuc := mock.NewMockUserUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(uuc)
			}

			handler := NewUserHandler(uuc)
			e := echo.New()

			e.GET("/api/user/get", handler.GetUser)

			req := tt.in()
			recorder := httptest.NewRecorder()

			e.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestUserHandler_CreateUser(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserUseCase) {
				m.EXPECT().CreateUserAndToken(
					gomock.Any(),
					"test@gmail.com",
					"password123",
				).Return(
					"eyJhbGciOiJIUzI1NiIsI.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0Ijo.SflKxwRJSMeKKF2QT4fwpMeJf36P",
					nil,
				)
			},
			in: func() *http.Request {
				userCreateReq := CreateUserRequest{Email: "test@gmail.com", Password: "password123"}
				reqBody, _ := json.Marshal(userCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request",
			in: func() *http.Request {
				userCreateReq := CreateUserRequest{Email: "test@gmail.com"}
				reqBody, _ := json.Marshal(userCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/user/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			uuc := mock.NewMockUserUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(uuc)
			}

			handler := NewUserHandler(uuc)
			e := echo.New()

			e.POST("/api/user/create", handler.CreateUser)

			req := tt.in()
			recorder := httptest.NewRecorder()

			e.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK {
				if token := recorder.Header().Get("Authorization"); token == "" || strings.TrimPrefix(token, "Bearer ") == "" {
					t.Fatalf("Expected Authorization header to be set")
				}
			}
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserUseCase) {
				m.EXPECT().UpdateUser(gomock.Any(), "updatedTest").Return(
					nil,
				)
			},
			in: func() *http.Request {
				userUpdateReq := UpdateUserRequest{Name: "updatedTest"}
				reqBody, _ := json.Marshal(userUpdateReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/user/update", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request of name",
			in: func() *http.Request {
				userUpdateReq := UpdateUserRequest{Name: ""}
				reqBody, _ := json.Marshal(userUpdateReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/user/update", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			uuc := mock.NewMockUserUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(uuc)
			}

			handler := NewUserHandler(uuc)
			e := echo.New()

			e.PUT("/api/user/update", handler.UpdateUser)

			req := tt.in()
			recorder := httptest.NewRecorder()

			e.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}
