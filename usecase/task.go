//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/tusmasoma/go-clean-arch/pkg/config"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type TaskUseCase interface {
	GetTask(ctx context.Context, id string) (*entity.Task, error)
	ListTasks(ctx context.Context) ([]entity.Task, error)
	CreateTask(ctx context.Context, params *CreateTaskParams) error
	UpdateTask(ctx context.Context, params *UpdateTaskParams) error
	DeleteTask(ctx context.Context, id string) error
}

type taskUseCase struct {
	tr repository.TaskRepository
}

func NewTaskUseCase(tr repository.TaskRepository) TaskUseCase {
	return &taskUseCase{
		tr: tr,
	}
}

func (tuc *taskUseCase) GetTask(ctx context.Context, id string) (*entity.Task, error) {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		return nil, errors.New("user name not found in request context")
	}
	task, err := tuc.tr.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if task.UserID != userID {
		return nil, errors.New("task does not belong to the user")
	}
	return task, nil
}

func (tuc *taskUseCase) ListTasks(ctx context.Context) ([]entity.Task, error) {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		return nil, errors.New("user name not found in request context")
	}
	tasks, err := tuc.tr.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

type CreateTaskParams struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Priority    int       `json:"priority"`
}

func (tuc *taskUseCase) CreateTask(ctx context.Context, params *CreateTaskParams) error {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		return errors.New("user name not found in request context")
	}
	task, err := entity.CreateTask(userID, params.Title, params.Description, params.DueDate, params.Priority)
	if err != nil {
		return err
	}
	if err = tuc.tr.Create(ctx, *task); err != nil {
		return err
	}
	return nil
}

type UpdateTaskParams struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Priority    int       `json:"priority"`
}

func (tuc *taskUseCase) UpdateTask(ctx context.Context, params *UpdateTaskParams) error {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		return errors.New("user name not found in request context")
	}
	task, err := tuc.tr.Get(ctx, params.ID)
	if err != nil {
		return err
	}
	if task.UserID != userID {
		return errors.New("task does not belong to the user")
	}
	if err = task.UpdateTask(params.Title, params.Description, params.DueDate, params.Priority); err != nil {
		return err
	}
	if err = tuc.tr.Update(ctx, *task); err != nil {
		return err
	}
	return nil
}

func (tuc *taskUseCase) DeleteTask(ctx context.Context, id string) error {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		return errors.New("user name not found in request context")
	}
	task, err := tuc.tr.Get(ctx, id)
	if err != nil {
		return err
	}
	if task.UserID != userID {
		return errors.New("task does not belong to the user")
	}
	if err = tuc.tr.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
