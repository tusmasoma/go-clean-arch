//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"errors"

	"github.com/tusmasoma/go-clean-arch/config"
	"github.com/tusmasoma/go-clean-arch/pkg/jwt"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type UserUseCase interface {
	GetUser(ctx context.Context) (*entity.User, error)
	CreateUserAndToken(ctx context.Context, email string, passward string) (string, error)
	UpdateUser(ctx context.Context, name, email string) error
}

type userUseCase struct {
	ur repository.UserRepository
	jg jwt.Generator
}

func NewUserUseCase(
	ur repository.UserRepository,
	jg jwt.Generator,
) UserUseCase {
	return &userUseCase{
		ur: ur,
		jg: jg,
	}
}

func (uuc *userUseCase) GetUser(ctx context.Context) (*entity.User, error) {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		return nil, errors.New("user name not found in request context")
	}
	user, err := uuc.ur.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uuc *userUseCase) CreateUserAndToken(ctx context.Context, email string, password string) (string, error) {
	user, err := entity.CreateUser(email, password)
	if err != nil {
		return "", err
	}
	if err = uuc.ur.Create(ctx, *user); err != nil {
		return "", err
	}
	jwt, _ := uuc.jg.GenerateToken(user.ID, user.Email)
	return jwt, nil
}

func (uuc *userUseCase) UpdateUser(ctx context.Context, name, email string) error {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		return errors.New("user name not found in request context")
	}
	user, err := uuc.ur.Get(ctx, userID)
	if err != nil {
		return err
	}
	if err = user.UpdateUser(name, email); err != nil {
		return err
	}
	if err = uuc.ur.Update(ctx, *user); err != nil {
		return err
	}
	return nil
}
