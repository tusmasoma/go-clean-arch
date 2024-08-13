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
	CreateTask(ctx context.Context, params *CreateTaskParams) error
}

type taskUseCase struct {
	tr repository.TaskRepository
}

func NewTaskUseCase(tr repository.TaskRepository) TaskUseCase {
	return &taskUseCase{
		tr: tr,
	}
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
