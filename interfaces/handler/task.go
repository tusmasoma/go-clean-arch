package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/usecase"
)

type TaskHandler interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	ListTasks(w http.ResponseWriter, r *http.Request)
	CreateTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
}

type taskHandler struct {
	tuc usecase.TaskUseCase
}

func NewTaskHandler(tuc usecase.TaskUseCase) TaskHandler {
	return &taskHandler{
		tuc: tuc,
	}
}

type GetTaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueData     time.Time `json:"due_date"`
	Priority    int       `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
}

func (th *taskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Warn("ID is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task, err := th.tuc.GetTask(ctx, id)
	if err != nil {
		log.Error("Failed to get task", log.Ferror(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(GetTaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueData:     task.DueData,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt,
	}); err != nil {
		http.Error(w, "Failed to encode task to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type ListTasksResponse struct {
	Tasks []struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		DueData     time.Time `json:"due_date"`
		Priority    int       `json:"priority"`
		CreatedAt   time.Time `json:"created_at"`
	} `json:"tasks"`
}

func (th *taskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := th.tuc.ListTasks(ctx)
	if err != nil {
		log.Error("Failed to list tasks", log.Ferror(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := th.convertTasksToListTasksResponse(tasks)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode tasks to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (th *taskHandler) convertTasksToListTasksResponse(tasks []entity.Task) ListTasksResponse {
	var tasksResponse []struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		DueData     time.Time `json:"due_date"`
		Priority    int       `json:"priority"`
		CreatedAt   time.Time `json:"created_at"`
	}
	for _, task := range tasks {
		tasksResponse = append(tasksResponse, struct {
			ID          string    `json:"id"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			DueData     time.Time `json:"due_date"`
			Priority    int       `json:"priority"`
			CreatedAt   time.Time `json:"created_at"`
		}{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			DueData:     task.DueData,
			Priority:    task.Priority,
			CreatedAt:   task.CreatedAt,
		})
	}
	return ListTasksResponse{
		Tasks: tasksResponse,
	}
}

type CreateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueData     time.Time `json:"due_date"`
	Priority    int       `json:"priority"`
}

func (th *taskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody CreateTaskRequest
	defer r.Body.Close()
	if !th.isValidCreateTasksRequest(r.Body, &requestBody) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := th.convertCreateTaskReqeuestToParams(requestBody)
	if err := th.tuc.CreateTask(ctx, params); err != nil {
		log.Error("Failed to create task", log.Ferror(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (th *taskHandler) isValidCreateTasksRequest(body io.ReadCloser, requestBody *CreateTaskRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Failed to decode request body", log.Ferror(err))
		return false
	}
	if requestBody.Title == "" ||
		requestBody.Description == "" ||
		requestBody.DueData.IsZero() ||
		requestBody.Priority < 1 ||
		requestBody.Priority > 5 {
		log.Warn("Invalid request body: %v", requestBody)
		return false
	}
	return true
}

func (th *taskHandler) convertCreateTaskReqeuestToParams(req CreateTaskRequest) *usecase.CreateTaskParams {
	return &usecase.CreateTaskParams{
		Title:       req.Title,
		Description: req.Description,
		DueData:     req.DueData,
		Priority:    req.Priority,
	}
}

type UpdateTaskRequest struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueData     time.Time `json:"due_date"`
	Priority    int       `json:"priority"`
}

func (th *taskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody UpdateTaskRequest
	defer r.Body.Close()
	if !th.isValidUpdateTasksRequest(r.Body, &requestBody) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := th.convertUpdateTaskReqeuestToParams(requestBody)
	if err := th.tuc.UpdateTask(ctx, params); err != nil {
		log.Error("Failed to update task", log.Ferror(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (th *taskHandler) isValidUpdateTasksRequest(body io.ReadCloser, requestBody *UpdateTaskRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Failed to decode request body", log.Ferror(err))
		return false
	}
	if requestBody.ID == "" ||
		requestBody.Title == "" ||
		requestBody.Description == "" ||
		requestBody.DueData.IsZero() ||
		requestBody.Priority < 1 ||
		requestBody.Priority > 5 {
		log.Warn("Invalid request body: %v", requestBody)
		return false
	}
	return true
}

func (th *taskHandler) convertUpdateTaskReqeuestToParams(req UpdateTaskRequest) *usecase.UpdateTaskParams {
	return &usecase.UpdateTaskParams{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		DueData:     req.DueData,
		Priority:    req.Priority,
	}
}

func (th *taskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Warn("ID is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := th.tuc.DeleteTask(ctx, id); err != nil {
		log.Error("Failed to delete task", log.Ferror(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
