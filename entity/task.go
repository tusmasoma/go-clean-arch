package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type Priority int

const (
	Low Priority = iota + 1
	MediumLow
	Medium
	MediumHigh
	High
)

var ValidPriorities = map[Priority]bool{
	Low:        true,
	MediumLow:  true,
	Medium:     true,
	MediumHigh: true,
	High:       true,
}

type Task struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	DueDate     time.Time `json:"due_date" db:"duedate"`
	Priority    Priority  `json:"priority" db:"priority"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

func NewTask(title, description string, dueDate time.Time, priority Priority) (*Task, error) {
	if title == "" {
		log.Error("title is required")
		return nil, errors.New("title is required")
	}
	if description == "" {
		log.Error("description is required")
		return nil, errors.New("description is required")
	}
	// TODO: Check if dueDate is in the future
	if !ValidPriorities[priority] {
		log.Error("priority must be between 1 and 5")
		return nil, errors.New("priority must be between 1 and 5")
	}
	return &Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Priority:    priority,
		CreatedAt:   time.Now(),
	}, nil
}
