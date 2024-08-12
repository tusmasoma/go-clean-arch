//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/go-clean-arch/entity"
)

type TaskRepository interface {
	Get(ctx context.Context, id string) (*entity.Task, error)
	List(ctx context.Context) ([]entity.Task, error)
	Create(ctx context.Context, task entity.Task) error
	Update(ctx context.Context, task entity.Task) error
	Delete(ctx context.Context, id string) error
}
