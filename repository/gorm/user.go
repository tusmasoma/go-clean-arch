package gorm

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) Get(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	if err := ur.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) Create(ctx context.Context, user entity.User) error {
	if err := ur.db.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) Update(ctx context.Context, user entity.User) error {
	if err := ur.db.WithContext(ctx).Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) Delete(ctx context.Context, id string) error {
	if err := ur.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) LockUserByEmail(ctx context.Context, email string) (bool, error) {
	var user entity.User
	if err := ur.db.WithContext(ctx).Set("gorm:query_option", "FOR UPDATE").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("No user found with the provided email", email)
			return false, nil
		}
		log.Println("Failed to scan row", err)
		return false, err
	}

	return true, nil
}
