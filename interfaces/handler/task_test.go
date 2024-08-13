package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/go-clean-arch/usecase"
	"github.com/tusmasoma/go-clean-arch/usecase/mock"
)

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
					&usecase.CreateTaskParams{
						Title:       "title",
						Description: "description",
						DueData:     dueDate,
						Priority:    3,
					},
				).Return(nil)
			},
			in: func() *http.Request {
				taskCreateReq := CreateTaskRequest{
					Title:       "title",
					Description: "description",
					DueData:     dueDate,
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
					DueData:     dueDate,
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
					DueData:     dueDate,
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
					DueData:     time.Time{},
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
					DueData:     dueDate,
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
					DueData:     dueDate,
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
			recorder := httptest.NewRecorder()
			handler.CreateTask(recorder, tt.in())

			if status := recorder.Code; status != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}
