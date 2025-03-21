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
	ID          string
	UserID      string
	Title       string
	Description string
	DueDate     time.Time
	Priority    Priority
	CreatedAt   time.Time
	IsOverdue   bool
	IsDueSoon   bool
}

func (t *Task) checkOverdue() bool {
	return time.Now().After(t.DueDate)
}

func (t *Task) checkDueSoon() bool {
	now := time.Now()
	return now.Before(t.DueDate) && time.Now().After(t.DueDate.AddDate(0, 0, -1))
}

func (t *Task) SetOverdue() {
	if t.checkOverdue() {
		t.IsOverdue = true
		return
	}
	t.IsOverdue = false
}

func (t *Task) SetDueSoon() {
	if t.checkDueSoon() {
		t.IsDueSoon = true
		return
	}
	t.IsDueSoon = false
}

func (t *Task) SetPriority(priority int) error {
	if !ValidPriorities[Priority(priority)] {
		return errors.New("priority must be between 1 and 5")
	}
	t.Priority = Priority(priority)
	return nil
}

func NewTask(id, userID, title, description string, dueDate time.Time, priority int, createdAt time.Time) (*Task, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if userID == "" {
		return nil, errors.New("userID is required")
	}
	if title == "" {
		return nil, errors.New("title is required")
	}
	if description == "" {
		return nil, errors.New("description is required")
	}
	task := &Task{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: description,
		DueDate:     dueDate.UTC().Truncate(time.Second),
		CreatedAt:   createdAt.UTC().Truncate(time.Second),
	}
	if err := task.SetPriority(priority); err != nil {
		return nil, err
	}
	task.SetOverdue()
	task.SetDueSoon()
	return task, nil
}

func CreateTask(userID, title, description string, dueDate time.Time, priority int) (*Task, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}
	if title == "" {
		return nil, errors.New("title is required")
	}
	if description == "" {
		return nil, errors.New("description is required")
	}
	task := &Task{
		ID:          uuid.New().String(),
		UserID:      userID,
		Title:       title,
		Description: description,
		DueDate:     dueDate.UTC().Truncate(time.Second),
		CreatedAt:   time.Now().UTC().Truncate(time.Second),
	}
	if err := task.SetPriority(priority); err != nil {
		return nil, err
	}
	task.SetOverdue()
	task.SetDueSoon()
	return task, nil
}

func (t *Task) UpdateTask(title, description string, dueDate time.Time, priority int) error {
	if title == "" {
		return errors.New("title is required")
	}
	if description == "" {
		return errors.New("description is required")
	}
	t.Title = title
	t.Description = description
	t.DueDate = dueDate.UTC().Truncate(time.Second)
	if err := t.SetPriority(priority); err != nil {
		return err
	}
	t.SetOverdue()
	t.SetDueSoon()
	return nil
}
