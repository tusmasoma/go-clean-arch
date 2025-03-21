package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
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
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Priority    Priority  `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	IsOverdue   bool      `json:"is_overdue"`
	IsDueSoon   bool      `json:"is_due_soon"`
}

func (t *Task) CheckOverdue() bool {
	return time.Now().After(t.DueDate)
}

func (t *Task) CheckDueSoon() bool {
	now := time.Now()
	return now.Before(t.DueDate) && time.Now().After(t.DueDate.AddDate(0, 0, -1))
}

func (t *Task) SetPriority(priority int) error {
	if !ValidPriorities[Priority(priority)] {
		return errors.New("priority must be between 1 and 5")
	}
	t.Priority = Priority(priority)
	return nil
}

func NewTask(userID, title, description string, dueDate time.Time, priority int) (*Task, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}
	if title == "" {
		return nil, errors.New("title is required")
	}
	if description == "" {
		return nil, errors.New("description is required")
	}
	// TODO: Check if dueDate is in the future
	if !ValidPriorities[Priority(priority)] {
		return nil, errors.New("priority must be between 1 and 5")
	}
	return &Task{
		ID:          uuid.New().String(),
		UserID:      userID,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Priority:    Priority(priority),
		CreatedAt:   time.Now(),
	}, nil
}
