package handler

import (
	"context"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/tusmasoma/go-clean-arch/proto/gateway"
	"github.com/tusmasoma/go-clean-arch/usecase"
)

type TaskHandler interface {
	GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error)
	ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error)
	CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error)
	UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error)
	DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error)
}

type taskHandler struct {
	tuc usecase.TaskUseCase
	pb.UnimplementedTaskServiceServer
}

func NewTaskHandler(tuc usecase.TaskUseCase) *taskHandler { //nolint:revive // This function is used in the test
	return &taskHandler{
		tuc: tuc,
	}
}

func (th *taskHandler) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	id := req.GetId()
	if id == "" {
		log.Warn("ID is required")
		return nil, status.Errorf(codes.InvalidArgument, "ID is required")
	}

	task, err := th.tuc.GetTask(ctx, id)
	if err != nil {
		log.Error("Failed to get task", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to get task")
	}

	return &pb.GetTaskResponse{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     timestamppb.New(task.DueData),
		Priority:    int32(task.Priority),
		CreatedAt:   timestamppb.New(task.CreatedAt),
	}, nil
}

func (th *taskHandler) ListTasks(ctx context.Context, _ *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	tasks, err := th.tuc.ListTasks(ctx)
	if err != nil {
		log.Error("Failed to list tasks", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to list tasks")
	}

	var res []*pb.Task
	for _, task := range tasks {
		res = append(res, &pb.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     timestamppb.New(task.DueData),
			Priority:    int32(task.Priority),
			CreatedAt:   timestamppb.New(task.CreatedAt),
		})
	}

	return &pb.ListTasksResponse{Tasks: res}, nil
}

func (th *taskHandler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	if !th.isValidCreateTasksRequest(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}
	params := th.convertCreateTaskReqeuestToParams(req)
	if err := th.tuc.CreateTask(ctx, params); err != nil {
		log.Error("Failed to create task", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to create task")
	}

	return &pb.CreateTaskResponse{}, nil
}

func (th *taskHandler) isValidCreateTasksRequest(req *pb.CreateTaskRequest) bool {
	if req.GetTitle() == "" ||
		req.GetDescription() == "" ||
		req.GetDueDate().AsTime().IsZero() ||
		int(req.GetPriority()) < 1 ||
		int(req.GetPriority()) > 5 {
		log.Warn(
			"Invalid request",
			log.Fstring("title", req.GetTitle()),
			log.Fstring("description", req.GetDescription()),
			log.Ftime("due_date", req.GetDueDate().AsTime()),
			log.Fint("priority", int(req.GetPriority())),
		)
		return false
	}
	return true
}

func (th *taskHandler) convertCreateTaskReqeuestToParams(req *pb.CreateTaskRequest) *usecase.CreateTaskParams {
	return &usecase.CreateTaskParams{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		DueData:     req.GetDueDate().AsTime(),
		Priority:    int(req.GetPriority()),
	}
}

func (th *taskHandler) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	if !th.isValidUpdateTasksRequest(req) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}
	params := th.convertUpdateTaskRequestToParams(req)
	if err := th.tuc.UpdateTask(ctx, params); err != nil {
		log.Error("Failed to update task", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to update task")
	}

	return &pb.UpdateTaskResponse{}, nil
}

func (th *taskHandler) isValidUpdateTasksRequest(req *pb.UpdateTaskRequest) bool {
	if req.GetId() == "" ||
		req.GetTitle() == "" ||
		req.GetDescription() == "" ||
		req.GetDueDate().AsTime().IsZero() ||
		int(req.GetPriority()) < 1 ||
		int(req.GetPriority()) > 5 {
		log.Warn(
			"Invalid request",
			log.Fstring("id", req.GetId()),
			log.Fstring("title", req.GetTitle()),
			log.Fstring("description", req.GetDescription()),
			log.Ftime("due_date", req.GetDueDate().AsTime()),
			log.Fint("priority", int(req.GetPriority())),
		)
		return false
	}
	return true
}

func (th *taskHandler) convertUpdateTaskRequestToParams(req *pb.UpdateTaskRequest) *usecase.UpdateTaskParams {
	return &usecase.UpdateTaskParams{
		ID:          req.GetId(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		DueData:     req.GetDueDate().AsTime(),
		Priority:    int(req.GetPriority()),
	}
}

func (th *taskHandler) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	id := req.GetId()
	if id == "" {
		log.Warn("ID is required")
		return nil, status.Errorf(codes.InvalidArgument, "ID is required")
	}

	if err := th.tuc.DeleteTask(ctx, id); err != nil {
		log.Error("Failed to delete task", log.Ferror(err))
		return nil, status.Errorf(codes.Internal, "Failed to delete task")
	}

	return &pb.DeleteTaskResponse{}, nil
}
