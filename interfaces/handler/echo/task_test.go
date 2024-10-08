package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/usecase"
	"github.com/tusmasoma/go-clean-arch/usecase/mock"
)

func TestHandler_GetTask(t *testing.T) {
	t.Parallel()

	taskID := uuid.New().String()
	dueDate := time.Now().AddDate(0, 0, 1)

	task := &entity.Task{
		ID:          taskID,
		Title:       "title",
		Description: "description",
		DueDate:     dueDate,
		Priority:    3,
		CreatedAt:   time.Now(),
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTaskUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockTaskUseCase) {
				tuc.EXPECT().GetTask(
					gomock.Any(),
					taskID,
				).Return(task, nil)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/task/get?id=%s", taskID), nil)
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request of id is empty",
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/api/task/get", nil)
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
			tuc := mock.NewMockTaskUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(tuc)
			}

			handler := NewTaskHandler(tuc)
			e := echo.New()

			e.GET("/api/task/get", handler.GetTask)

			req := tt.in()
			recorder := httptest.NewRecorder()

			e.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestHandler_ListTasks(t *testing.T) {
	t.Parallel()

	dueDate := time.Now().AddDate(0, 0, 1)

	tasks := []entity.Task{
		{
			ID:          uuid.New().String(),
			Title:       "title1",
			Description: "description1",
			DueDate:     dueDate,
			Priority:    3,
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			Title:       "title2",
			Description: "description2",
			DueDate:     dueDate,
			Priority:    3,
			CreatedAt:   time.Now(),
		},
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTaskUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockTaskUseCase) {
				tuc.EXPECT().ListTasks(
					gomock.Any(),
				).Return(tasks, nil)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/api/task/list", nil)
				return req
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tuc := mock.NewMockTaskUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(tuc)
			}

			handler := NewTaskHandler(tuc)
			e := echo.New()

			e.GET("/api/task/list", handler.ListTasks)

			req := tt.in()
			recorder := httptest.NewRecorder()

			e.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestHandler_CreateTask(t *testing.T) {
	t.Parallel()

	dueDate := time.Now().AddDate(0, 0, 1)

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTaskUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockTaskUseCase) {
				tuc.EXPECT().CreateTask(
					gomock.Any(),
					gomock.Any(),
				).Do(func(_ context.Context, params *usecase.CreateTaskParams) {
					if params.Title != "title" {
						t.Errorf("unexpected Title: got %v, want %v", params.Title, "title")
					}
					if params.Description != "description" {
						t.Errorf("unexpected Description: got %v, want %v", params.Description, "description")
					}
					if !params.DueDate.Equal(dueDate) {
						t.Errorf("unexpected DueDate: got %v, want %v", params.DueDate, dueDate)
					}
					if params.Priority != 3 {
						t.Errorf("unexpected Priority: got %v, want %v", params.Priority, 3)
					}
				}).Return(nil)
			},
			in: func() *http.Request {
				taskCreateReq := CreateTaskRequest{
					Title:       "title",
					Description: "description",
					DueDate:     dueDate,
					Priority:    3,
				}
				reqBody, _ := json.Marshal(taskCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/task/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request of title is empty",
			in: func() *http.Request {
				taskCreateReq := CreateTaskRequest{
					Title:       "",
					Description: "description",
					DueDate:     dueDate,
					Priority:    3,
				}
				reqBody, _ := json.Marshal(taskCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/task/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: invalid request of description is empty",
			in: func() *http.Request {
				taskCreateReq := CreateTaskRequest{
					Title:       "title",
					Description: "",
					DueDate:     dueDate,
					Priority:    3,
				}
				reqBody, _ := json.Marshal(taskCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/task/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: invalid request of due_date is zero",
			in: func() *http.Request {
				taskCreateReq := CreateTaskRequest{
					Title:       "title",
					Description: "description",
					DueDate:     time.Time{},
					Priority:    3,
				}
				reqBody, _ := json.Marshal(taskCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/task/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: invalid request of priority is less than 1",
			in: func() *http.Request {
				taskCreateReq := CreateTaskRequest{
					Title:       "title",
					Description: "description",
					DueDate:     dueDate,
					Priority:    0,
				}
				reqBody, _ := json.Marshal(taskCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/task/create", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: invalid request of priority is greater than 5",
			in: func() *http.Request {
				taskCreateReq := CreateTaskRequest{
					Title:       "title",
					Description: "description",
					DueDate:     dueDate,
					Priority:    6,
				}
				reqBody, _ := json.Marshal(taskCreateReq)
				req, _ := http.NewRequest(http.MethodPost, "/api/task/create", bytes.NewBuffer(reqBody))
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
			tuc := mock.NewMockTaskUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(tuc)
			}

			handler := NewTaskHandler(tuc)
			e := echo.New()

			e.POST("/api/task/create", handler.CreateTask)

			req := tt.in()
			recorder := httptest.NewRecorder()

			e.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestHandler_UpdateTask(t *testing.T) {
	t.Parallel()

	taskID := uuid.New().String()
	dueDate := time.Now().AddDate(0, 0, 1)

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTaskUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockTaskUseCase) {
				tuc.EXPECT().UpdateTask(
					gomock.Any(),
					gomock.Any(),
				).Do(func(_ context.Context, params *usecase.UpdateTaskParams) {
					if params.Title != "updated title" {
						t.Errorf("unexpected Title: got %v, want %v", params.Title, "title")
					}
					if params.Description != "updated description" {
						t.Errorf("unexpected Description: got %v, want %v", params.Description, "description")
					}
					if !params.DueDate.Equal(dueDate) {
						t.Errorf("unexpected DueDate: got %v, want %v", params.DueDate, dueDate)
					}
					if params.Priority != 2 {
						t.Errorf("unexpected Priority: got %v, want %v", params.Priority, 3)
					}
				}).Return(nil)
			},
			in: func() *http.Request {
				taskUpdateReq := UpdateTaskRequest{
					ID:          taskID,
					Title:       "updated title",
					Description: "updated description",
					DueDate:     dueDate,
					Priority:    2,
				}
				reqBody, _ := json.Marshal(taskUpdateReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/task/update", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request of id is empty",
			in: func() *http.Request {
				taskUpdateReq := UpdateTaskRequest{
					ID:          "",
					Title:       "updated title",
					Description: "updated description",
					DueDate:     dueDate,
					Priority:    2,
				}
				reqBody, _ := json.Marshal(taskUpdateReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/task/update", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: invalid request of title is empty",
			in: func() *http.Request {
				taskUpdateReq := UpdateTaskRequest{
					ID:          taskID,
					Title:       "",
					Description: "updated description",
					DueDate:     dueDate,
					Priority:    2,
				}
				reqBody, _ := json.Marshal(taskUpdateReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/task/update", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: invalid request of description is empty",
			in: func() *http.Request {
				taskUpdateReq := UpdateTaskRequest{
					ID:          taskID,
					Title:       "updated title",
					Description: "",
					DueDate:     dueDate,
					Priority:    2,
				}
				reqBody, _ := json.Marshal(taskUpdateReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/task/update", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: invalid request of due_date is zero",
			in: func() *http.Request {
				taskUpdateReq := UpdateTaskRequest{
					ID:          taskID,
					Title:       "updated title",
					Description: "updated description",
					DueDate:     time.Time{},
					Priority:    2,
				}
				reqBody, _ := json.Marshal(taskUpdateReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/task/update", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: invalid request of priority is less than 1",
			in: func() *http.Request {
				taskUpdateReq := UpdateTaskRequest{
					ID:          taskID,
					Title:       "updated title",
					Description: "updated description",
					DueDate:     dueDate,
					Priority:    0,
				}
				reqBody, _ := json.Marshal(taskUpdateReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/task/update", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail: invalid request of priority is greater than 5",
			in: func() *http.Request {
				taskUpdateReq := UpdateTaskRequest{
					ID:          taskID,
					Title:       "updated title",
					Description: "updated description",
					DueDate:     dueDate,
					Priority:    6,
				}
				reqBody, _ := json.Marshal(taskUpdateReq)
				req, _ := http.NewRequest(http.MethodPut, "/api/task/update", bytes.NewBuffer(reqBody))
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
			tuc := mock.NewMockTaskUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(tuc)
			}

			handler := NewTaskHandler(tuc)
			e := echo.New()

			e.PUT("/api/task/update", handler.UpdateTask)

			req := tt.in()
			recorder := httptest.NewRecorder()

			e.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestHandler_DeleteTask(t *testing.T) {
	t.Parallel()

	taskID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTaskUseCase,
		)
		in         func() *http.Request
		wantStatus int
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockTaskUseCase) {
				tuc.EXPECT().DeleteTask(
					gomock.Any(),
					taskID,
				).Return(nil)
			},
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/task/delete?id=%s", taskID), nil)
				return req
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: invalid request of id is empty",
			in: func() *http.Request {
				req, _ := http.NewRequest(http.MethodDelete, "/api/task/delete", nil)
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
			tuc := mock.NewMockTaskUseCase(ctrl)

			if tt.setup != nil {
				tt.setup(tuc)
			}

			handler := NewTaskHandler(tuc)
			e := echo.New()

			e.DELETE("/api/task/delete", handler.DeleteTask)

			req := tt.in()
			recorder := httptest.NewRecorder()

			e.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}
