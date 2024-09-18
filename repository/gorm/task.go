package gorm

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type taskModel struct {
	ID          string    `gorm:"type:char(36);primaryKey"`
	UserID      string    `gorm:"column:user_id"`
	Title       string    `gorm:"column:title"`
	Description string    `gorm:"column:description"`
	DueDate     time.Time `gorm:"column:duedate"`
	Priority    int       `gorm:"column:priority"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) repository.TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (tr *taskRepository) Get(ctx context.Context, id string) (*entity.Task, error) {
	executor := tr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	var tm taskModel
	if err := executor.WithContext(ctx).First(&tm, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &entity.Task{
		ID:          tm.ID,
		UserID:      tm.UserID,
		Title:       tm.Title,
		Description: tm.Description,
		DueDate:     tm.DueDate,
		Priority:    tm.Priority,
		CreatedAt:   tm.CreatedAt,
	}, nil
}

func (tr *taskRepository) List(ctx context.Context, userID string) ([]entity.Task, error) {
	executor := tr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	var tms []taskModel
	if err := executor.WithContext(ctx).Find(&tms, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	tasks := make([]entity.Task, len(tms))
	for i, tm := range tms {
		tasks[i] = entity.Task{
			ID:          tm.ID,
			UserID:      tm.UserID,
			Title:       tm.Title,
			Description: tm.Description,
			DueDate:     tm.DueDate,
			Priority:    tm.Priority,
			CreatedAt:   tm.CreatedAt,
		}
	}
	return tasks, nil
}

func (tr *taskRepository) Create(ctx context.Context, task entity.Task) error {
	executor := tr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Create(&taskModel{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) Update(ctx context.Context, task entity.Task) error {
	executor := tr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Model(&taskModel{}).Where("id = ?", task.ID).Updates(&taskModel{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (tr *taskRepository) Delete(ctx context.Context, id string) error {
	executor := tr.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Delete(&taskModel{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
