//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package repository

import (
	"context"

	"github.com/tusmasoma/go-clean-arch/entity"
)

type UserRepository interface {
	Get(ctx context.Context, id string) (*entity.User, error)
	Create(ctx context.Context, user entity.User) error
	Update(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id string) error
	LockUserByEmail(ctx context.Context, email string) (bool, error)
}
