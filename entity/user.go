package entity

import (
	"errors"

	"github.com/google/uuid"
	pe "github.com/tusmasoma/go-clean-arch/pkg/email"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(email, password string) (*User, error) {
	if email == "" {
		log.Error("email is required")
		return nil, errors.New("email is required")
	}
	if password == "" {
		log.Error("password is required")
		return nil, errors.New("password is required")
	}
	name, err := pe.GetAddressPart(email)
	if err != nil {
		return nil, errors.New("email is required")
	}
	return &User{
		ID:       uuid.New().String(),
		Name:     name,
		Email:    email,
		Password: password,
	}, nil
}
