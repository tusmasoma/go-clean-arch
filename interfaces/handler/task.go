package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/usecase"
)

type TaskHandler interface {
	// GetTask(w http.ResponseWriter, r *http.Request)
	// ListTasks(w http.ResponseWriter, r *http.Request)
	CreateTask(w http.ResponseWriter, r *http.Request)
	// UpdateTask(w http.ResponseWriter, r *http.Request)
	// DeleteTask(w http.ResponseWriter, r *http.Request)
}

type taskHandler struct {
	tuc usecase.TaskUseCase
}

func NewTaskHandler(tuc usecase.TaskUseCase) TaskHandler {
	return &taskHandler{
		tuc: tuc,
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
