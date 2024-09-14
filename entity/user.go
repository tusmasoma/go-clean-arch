package entity

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
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
	name := extractNameFromEmail(email)
	return &User{
		ID:       uuid.New().String(),
		Name:     name,
		Email:    email,
		Password: password,
	}, nil
}

func extractNameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}
