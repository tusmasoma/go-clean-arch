package entity

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	pe "github.com/tusmasoma/go-clean-arch/pkg/email"
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
}

func NewUser(id, name, email, passwordHash string) (*User, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if passwordHash == "" {
		return nil, errors.New("hash password is required")
	}
	return &User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}

func CreateUser(email, password string) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	name, err := pe.GetAddressPart(email)
	if err != nil {
		return nil, errors.New("email is required")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	return &User{
		ID:           uuid.New().String(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}, nil
}

func (u *User) UpdateUser(name, email string) error {
	if name == "" {
		return errors.New("name is required")
	}
	if email == "" {
		return errors.New("email is required")
	}
	u.Name = name
	u.Email = email
	return nil
}
