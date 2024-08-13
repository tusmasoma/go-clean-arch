//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"time"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type TaskUseCase interface {
	GetTask(ctx context.Context, id string) (*entity.Task, error)
	CreateTask(ctx context.Context, params *CreateTaskParams) error
	UpdateTask(ctx context.Context, params *UpdateTaskParams) error
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
	task, err := tuc.tr.Get(ctx, id)
	if err != nil {
		log.Error("Failed to get task", log.Ferror(err))
		return nil, err
	}
	return task, nil
}

type CreateTaskParams struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueData     time.Time `json:"due_date"`
	Priority    int       `json:"priority"`
}

func (tuc *taskUseCase) CreateTask(ctx context.Context, params *CreateTaskParams) error {
	task, err := entity.NewTask(params.Title, params.Description, params.DueData, params.Priority)
	if err != nil {
		log.Error("Failed to create task", log.Ferror(err))
		return err
	}
	if err = tuc.tr.Create(ctx, *task); err != nil {
		log.Error("Failed to create task", log.Ferror(err))
		return err
	}
	return nil
}

type UpdateTaskParams struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueData     time.Time `json:"due_date"`
	Priority    int       `json:"priority"`
}

func (tuc *taskUseCase) UpdateTask(ctx context.Context, params *UpdateTaskParams) error {
	task, err := tuc.tr.Get(ctx, params.ID)
	if err != nil {
		log.Error("Failed to get task", log.Ferror(err))
		return err
	}

	task.Title = params.Title
	task.Description = params.Description
	task.DueData = params.DueData
	task.Priority = params.Priority

	if err = tuc.tr.Update(ctx, *task); err != nil {
		log.Error("Failed to update task", log.Ferror(err))
		return err
	}
	return nil
}
