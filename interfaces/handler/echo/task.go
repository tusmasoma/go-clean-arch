package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/usecase"
)

type TaskHandler interface {
	GetTask(c echo.Context) error
	ListTasks(c echo.Context) error
	CreateTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	DeleteTask(c echo.Context) error
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
	DueDate     time.Time `json:"due_date"`
	Priority    int       `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
}

func (th *taskHandler) GetTask(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.QueryParam("id")
	if id == "" {
		log.Warn("ID is required")
		return c.NoContent(http.StatusBadRequest)
	}

	task, err := th.tuc.GetTask(ctx, id)
	if err != nil {
		log.Error("Failed to get task", log.Ferror(err))
		return c.NoContent(http.StatusInternalServerError)
	}
	response := GetTaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt,
	}
	return c.JSON(http.StatusOK, response)
}

type ListTasksResponse struct {
	Tasks []struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		DueDate     time.Time `json:"due_date"`
		Priority    int       `json:"priority"`
		CreatedAt   time.Time `json:"created_at"`
	} `json:"tasks"`
}

func (th *taskHandler) ListTasks(c echo.Context) error {
	ctx := c.Request().Context()

	tasks, err := th.tuc.ListTasks(ctx)
	if err != nil {
		log.Error("Failed to list tasks", log.Ferror(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	response := th.convertTasksToListTasksResponse(tasks)
	return c.JSON(http.StatusOK, response)
}

func (th *taskHandler) convertTasksToListTasksResponse(tasks []entity.Task) ListTasksResponse {
	var tasksResponse []struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		DueDate     time.Time `json:"due_date"`
		Priority    int       `json:"priority"`
		CreatedAt   time.Time `json:"created_at"`
	}
	for _, task := range tasks {
		tasksResponse = append(tasksResponse, struct {
			ID          string    `json:"id"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			DueDate     time.Time `json:"due_date"`
			Priority    int       `json:"priority"`
			CreatedAt   time.Time `json:"created_at"`
		}{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
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
	DueDate     time.Time `json:"due_date"`
	Priority    int       `json:"priority"`
}

func (th *taskHandler) CreateTask(c echo.Context) error {
	ctx := c.Request().Context()

	var requestBody CreateTaskRequest
	if err := c.Bind(&requestBody); err != nil {
		log.Error("Failed to decode request body", log.Ferror(err))
		return c.NoContent(http.StatusBadRequest)
	}
	if !th.isValidCreateTasksRequest(&requestBody) {
		return c.NoContent(http.StatusBadRequest)
	}

	params := th.convertCreateTaskReqeuestToParams(requestBody)
	if err := th.tuc.CreateTask(ctx, params); err != nil {
		log.Error("Failed to create task", log.Ferror(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (th *taskHandler) isValidCreateTasksRequest(requestBody *CreateTaskRequest) bool {
	if requestBody.Title == "" ||
		requestBody.Description == "" ||
		requestBody.DueDate.IsZero() ||
		!entity.ValidPriorities[requestBody.Priority] {
		log.Warn("Invalid request body: %v", requestBody)
		return false
	}
	return true
}

func (th *taskHandler) convertCreateTaskReqeuestToParams(req CreateTaskRequest) *usecase.CreateTaskParams {
	return &usecase.CreateTaskParams{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Priority:    req.Priority,
	}
}

type UpdateTaskRequest struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Priority    int       `json:"priority"`
}

func (th *taskHandler) UpdateTask(c echo.Context) error {
	ctx := c.Request().Context()

	var requestBody UpdateTaskRequest
	if err := c.Bind(&requestBody); err != nil {
		log.Error("Failed to decode request body", log.Ferror(err))
		return c.NoContent(http.StatusBadRequest)
	}
	if !th.isValidUpdateTasksRequest(&requestBody) {
		return c.NoContent(http.StatusBadRequest)
	}

	params := th.convertUpdateTaskReqeuestToParams(requestBody)
	if err := th.tuc.UpdateTask(ctx, params); err != nil {
		log.Error("Failed to update task", log.Ferror(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (th *taskHandler) isValidUpdateTasksRequest(requestBody *UpdateTaskRequest) bool {
	if requestBody.ID == "" ||
		requestBody.Title == "" ||
		requestBody.Description == "" ||
		requestBody.DueDate.IsZero() ||
		!entity.ValidPriorities[requestBody.Priority] {
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
		DueDate:     req.DueDate,
		Priority:    req.Priority,
	}
}

func (th *taskHandler) DeleteTask(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.QueryParam("id")
	if id == "" {
		log.Warn("ID is required")
		return c.NoContent(http.StatusBadRequest)
	}

	if err := th.tuc.DeleteTask(ctx, id); err != nil {
		log.Error("Failed to delete task", log.Ferror(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
