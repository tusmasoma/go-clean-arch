package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

// The introduction of a custom Priority type was considered,
// but it requires additional implementation for JSON decoding.
// Therefore, for now, the standard int type is being used.
// Introducing a Priority type in the future could improve code readability and safety,
// so this decision should be revisited.
//
// type Priority int
//
// const (
// 	Low Priority = iota + 1
// 	MediumLow
// 	Medium
// 	MediumHigh
// 	High
// )
//
// var ValidPriorities = map[Priority]bool{
// 	Low:        true,
// 	MediumLow:  true,
// 	Medium:     true,
// 	MediumHigh: true,
// 	High:       true,
// }

const (
	Low int = iota + 1
	MediumLow
	Medium
	MediumHigh
	High
)

var ValidPriorities = map[int]bool{
	Low:        true,
	MediumLow:  true,
	Medium:     true,
	MediumHigh: true,
	High:       true,
}

type Task struct {
	ID          string    `json:"id" db:"id" bson:"_id,omitempty"`
	UserID      string    `json:"user_id" db:"user_id" bson:"user_id"`
	Title       string    `json:"title" db:"title" bson:"title"`
	Description string    `json:"description" db:"description" bson:"description"`
	DueDate     time.Time `json:"due_date" db:"duedate" bson:"duedate"`
	Priority    int       `json:"priority" db:"priority" bson:"priority"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" bson:"created_at"`
}

func NewTask(userID, title, description string, dueDate time.Time, priority int) (*Task, error) {
	if userID == "" {
		log.Error("userID is required")
		return nil, errors.New("userID is required")
	}
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
		UserID:      userID,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Priority:    priority,
		CreatedAt:   time.Now(),
	}, nil
}
