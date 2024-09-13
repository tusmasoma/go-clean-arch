package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) repository.TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (tr *taskRepository) Get(ctx context.Context, id string) (*entity.Task, error) {
	var task entity.Task
	if err := tr.db.WithContext(ctx).First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (tr *taskRepository) List(ctx context.Context, userID string) ([]entity.Task, error) {
	var tasks []entity.Task
	if err := tr.db.WithContext(ctx).Find(&tasks, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (tr *taskRepository) Create(ctx context.Context, task entity.Task) error {
	if err := tr.db.WithContext(ctx).Create(&task).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) Update(ctx context.Context, task entity.Task) error {
	if err := tr.db.WithContext(ctx).Save(&task).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) Delete(ctx context.Context, id string) error {
	if err := tr.db.WithContext(ctx).Delete(&entity.Task{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
