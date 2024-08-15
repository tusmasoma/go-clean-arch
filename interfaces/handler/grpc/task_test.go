package handler

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tusmasoma/go-clean-arch/entity"
	pb "github.com/tusmasoma/go-clean-arch/interfaces/handler/grpc/proto/gateway"
	"github.com/tusmasoma/go-clean-arch/usecase"
	"github.com/tusmasoma/go-clean-arch/usecase/mock"
)

const bufSize = 1024 * 1024

func setupTestServer(t *testing.T, setup func(m *mock.MockTaskUseCase)) (pb.TaskServiceClient, func()) {
	t.Helper()

	ctrl := gomock.NewController(t)
	tuc := mock.NewMockTaskUseCase(ctrl)

	if setup != nil {
		setup(tuc)
	}

	handler := NewTaskHandler(tuc)

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, handler)

	go func() {
		if err := s.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			t.Errorf("failed to serve: %v", err)
		}
	}()

	conn, err := grpc.Dial("", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) { //nolint:staticcheck // ignore deprecation
		return lis.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	client := pb.NewTaskServiceClient(conn)

	cleanup := func() {
		conn.Close()
		s.Stop()
	}

	return client, cleanup
}

func TestHandler_GetTask(t *testing.T) {
	t.Parallel()

	taskID := uuid.New().String()
	dueDate := time.Now().AddDate(0, 0, 1)

	task := &entity.Task{
		ID:          taskID,
		Title:       "title",
		Description: "description",
		DueData:     dueDate,
		Priority:    3,
		CreatedAt:   time.Now(),
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTaskUseCase,
		)
		request    *pb.GetTaskRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockTaskUseCase) {
				tuc.EXPECT().GetTask(
					gomock.Any(),
					taskID,
				).Return(task, nil)
			},
			request:    &pb.GetTaskRequest{Id: taskID},
			wantStatus: codes.OK,
		},
		{
			name:       "Fail: invalid request of id is empty",
			request:    &pb.GetTaskRequest{Id: ""},
			wantStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			resp, err := client.GetTask(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}

			if tt.wantStatus == codes.OK {
				if resp.GetId() != task.ID || resp.GetTitle() != task.Title || resp.GetDescription() != task.Description {
					t.Fatalf("handler returned wrong task data")
				}
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
			DueData:     dueDate,
			Priority:    3,
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			Title:       "title2",
			Description: "description2",
			DueData:     dueDate,
			Priority:    3,
			CreatedAt:   time.Now(),
		},
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTaskUseCase,
		)
		request    *pb.ListTasksRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockTaskUseCase) {
				tuc.EXPECT().ListTasks(
					gomock.Any(),
				).Return(tasks, nil)
			},
			request:    &pb.ListTasksRequest{},
			wantStatus: codes.OK,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			resp, err := client.ListTasks(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}

			if tt.wantStatus == codes.OK {
				if len(resp.GetTasks()) != len(tasks) {
					t.Fatalf("handler returned wrong task data")
				}
			}
		})
	}
}

func TestHandler_CreateTask(t *testing.T) { //nolint:gocognit // ignore cognitive complexity
	t.Parallel()

	dueDate := time.Now().AddDate(0, 0, 1)

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTaskUseCase,
		)
		request    *pb.CreateTaskRequest
		wantStatus codes.Code
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
					if !params.DueData.Equal(dueDate) {
						t.Errorf("unexpected DueData: got %v, want %v", params.DueData, dueDate)
					}
					if params.Priority != 3 {
						t.Errorf("unexpected Priority: got %v, want %v", params.Priority, 3)
					}
				}).Return(nil)
			},
			request: &pb.CreateTaskRequest{
				Title:       "title",
				Description: "description",
				DueDate:     timestamppb.New(dueDate),
				Priority:    3,
			},
			wantStatus: codes.OK,
		},
		{
			name: "Fail: invalid request of title is empty",
			request: &pb.CreateTaskRequest{
				Title:       "",
				Description: "description",
				DueDate:     timestamppb.New(dueDate),
				Priority:    3,
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of description is empty",
			request: &pb.CreateTaskRequest{
				Title:       "title",
				Description: "",
				DueDate:     timestamppb.New(dueDate),
				Priority:    3,
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of due_date is zero",
			request: &pb.CreateTaskRequest{
				Title:       "title",
				Description: "",
				DueDate:     timestamppb.New(time.Time{}),
				Priority:    3,
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of priority is less than 1",
			request: &pb.CreateTaskRequest{
				Title:       "title",
				Description: "description",
				DueDate:     timestamppb.New(dueDate),
				Priority:    0,
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of priority is greater than 5",
			request: &pb.CreateTaskRequest{
				Title:       "title",
				Description: "description",
				DueDate:     timestamppb.New(dueDate),
				Priority:    6,
			},
			wantStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			req, err := client.CreateTask(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}

			if tt.wantStatus == codes.OK {
				if req == nil {
					t.Fatalf("handler returned wrong task data")
				}
			}
		})
	}
}

func TestHandler_UpdateTask(t *testing.T) { //nolint:gocognit // ignore cognitive complexity
	t.Parallel()

	taskID := uuid.New().String()
	dueDate := time.Now().AddDate(0, 0, 1)

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTaskUseCase,
		)
		request    *pb.UpdateTaskRequest
		wantStatus codes.Code
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
					if !params.DueData.Equal(dueDate) {
						t.Errorf("unexpected DueData: got %v, want %v", params.DueData, dueDate)
					}
					if params.Priority != 2 {
						t.Errorf("unexpected Priority: got %v, want %v", params.Priority, 3)
					}
				}).Return(nil)
			},
			request: &pb.UpdateTaskRequest{
				Id:          taskID,
				Title:       "updated title",
				Description: "updated description",
				DueDate:     timestamppb.New(dueDate),
				Priority:    2,
			},
			wantStatus: codes.OK,
		},
		{
			name: "Fail: invalid request of id is empty",
			request: &pb.UpdateTaskRequest{
				Id:          "",
				Title:       "updated title",
				Description: "updated description",
				DueDate:     timestamppb.New(dueDate),
				Priority:    2,
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of title is empty",
			request: &pb.UpdateTaskRequest{
				Id:          taskID,
				Title:       "",
				Description: "updated description",
				DueDate:     timestamppb.New(dueDate),
				Priority:    2,
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of description is empty",
			request: &pb.UpdateTaskRequest{
				Id:          taskID,
				Title:       "updated title",
				Description: "",
				DueDate:     timestamppb.New(dueDate),
				Priority:    2,
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of due_date is zero",
			request: &pb.UpdateTaskRequest{
				Id:          taskID,
				Title:       "updated title",
				Description: "updated description",
				DueDate:     timestamppb.New(time.Time{}),
				Priority:    2,
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of priority is less than 1",
			request: &pb.UpdateTaskRequest{
				Id:          taskID,
				Title:       "updated title",
				Description: "updated description",
				DueDate:     timestamppb.New(dueDate),
				Priority:    0,
			},
			wantStatus: codes.InvalidArgument,
		},
		{
			name: "Fail: invalid request of priority is greater than 5",
			request: &pb.UpdateTaskRequest{
				Id:          taskID,
				Title:       "updated title",
				Description: "updated description",
				DueDate:     timestamppb.New(dueDate),
				Priority:    6,
			},
			wantStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			req, err := client.UpdateTask(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}

			if tt.wantStatus == codes.OK {
				if req == nil {
					t.Fatalf("handler returned wrong task data")
				}
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
		request    *pb.DeleteTaskRequest
		wantStatus codes.Code
	}{
		{
			name: "success",
			setup: func(tuc *mock.MockTaskUseCase) {
				tuc.EXPECT().DeleteTask(
					gomock.Any(),
					taskID,
				).Return(nil)
			},
			request:    &pb.DeleteTaskRequest{Id: taskID},
			wantStatus: codes.OK,
		},
		{
			name:       "Fail: invalid request of id is empty",
			request:    &pb.DeleteTaskRequest{Id: ""},
			wantStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, cleanup := setupTestServer(t, tt.setup)
			defer cleanup()

			req, err := client.DeleteTask(context.Background(), tt.request)
			if status.Code(err) != tt.wantStatus {
				t.Fatalf("handler returned wrong status code: got %v want %v", status.Code(err), tt.wantStatus)
			}

			if tt.wantStatus == codes.OK {
				if req == nil {
					t.Fatalf("handler returned wrong task data")
				}
			}
		})
	}
}
