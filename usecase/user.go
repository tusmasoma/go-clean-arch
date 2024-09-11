//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"errors"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/config"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type UserUseCase interface {
	GetUser(ctx context.Context) (*entity.User, error)
	CreateUserAndToken(ctx context.Context, email string, passward string) (string, error)
	UpdateUser(ctx context.Context, name string) error
}

type userUseCase struct {
	ur repository.UserRepository
	tr repository.TransactionRepository
	ar repository.AuthRepository
}

func NewUserUseCase(
	ur repository.UserRepository,
	tr repository.TransactionRepository,
	ar repository.AuthRepository,
) UserUseCase {
	return &userUseCase{
		ur: ur,
		tr: tr,
		ar: ar,
	}
}

func (uuc *userUseCase) GetUser(ctx context.Context) (*entity.User, error) {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		log.Error("User ID not found in request context")
		return nil, errors.New("user name not found in request context")
	}
	user, err := uuc.ur.Get(ctx, userID)
	if err != nil {
		log.Error("Error getting user", log.Fstring("user_id", userID))
		return nil, err
	}
	return user, nil
}

func (uuc *userUseCase) CreateUserAndToken(ctx context.Context, email string, password string) (string, error) {
	var user *entity.User
	if err := uuc.tr.Transaction(ctx, func(ctx context.Context) error {
		exists, err := uuc.ur.LockUserByEmail(ctx, email)
		if err != nil {
			log.Error("Error retrieving user by email", log.Fstring("email", email))
			return err
		}
		if exists {
			log.Info("User with this email already exists", log.Fstring("email", email))
			return errors.New("user with this email already exists")
		}

		user, err = entity.NewUser(email, password) // hash password
		if err != nil {
			log.Error("Error creating new user", log.Fstring("email", email))
			return err
		}

		if err = uuc.ur.Create(ctx, *user); err != nil {
			log.Error("Error creating new user", log.Fstring("email", email))
			return err
		}

		return nil
	}); err != nil {
		return "", err
	}

	jwt, _ := uuc.ar.GenerateToken(user.ID, user.Email)
	return jwt, nil
}

func (uuc *userUseCase) UpdateUser(ctx context.Context, name string) error {
	userIDValue := ctx.Value(config.ContextUserIDKey)
	userID, ok := userIDValue.(string)
	if !ok {
		log.Error("User ID not found in request context")
		return errors.New("user name not found in request context")
	}
	user, err := uuc.ur.Get(ctx, userID)
	if err != nil {
		log.Error("Error getting user", log.Fstring("user_id", userID))
		return err
	}

	// TODO: setter method for user
	user.Name = name
	// user.Email = email
	// user.Password = password

	if err = uuc.ur.Update(ctx, *user); err != nil {
		log.Error("Error updating user", log.Fstring("user_id", userID))
		return err
	}
	return nil
}
