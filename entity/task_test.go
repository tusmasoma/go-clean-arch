package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestEntity_NewTask(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	dueDate := time.Now().AddDate(0, 0, 1)

	patterns := []struct {
		name string
		arg  struct {
			userID      string
			title       string
			description string
			dueDate     time.Time
			priority    int
		}
		want struct {
			task *Task
			err  error
		}
	}{
		{
			name: "success",
			arg: struct {
				userID      string
				title       string
				description string
				dueDate     time.Time
				priority    int
			}{
				userID:      userID,
				title:       "title",
				description: "description",
				dueDate:     dueDate,
				priority:    Medium,
			},
			want: struct {
				task *Task
				err  error
			}{
				task: &Task{
					UserID:      userID,
					Title:       "title",
					Description: "description",
					DueDate:     dueDate,
					Priority:    Medium,
				},
				err: nil,
			},
		},
		{
			name: "Fail: userID is empty",
			arg: struct {
				userID      string
				title       string
				description string
				dueDate     time.Time
				priority    int
			}{
				userID:      "",
				title:       "title",
				description: "description",
				dueDate:     dueDate,
				priority:    Medium,
			},
			want: struct {
				task *Task
				err  error
			}{
				task: nil,
				err:  errors.New("userID is required"),
			},
		},
		{
			name: "Fail: title is empty",
			arg: struct {
				userID      string
				title       string
				description string
				dueDate     time.Time
				priority    int
			}{
				userID:      userID,
				title:       "",
				description: "description",
				dueDate:     dueDate,
				priority:    Medium,
			},
			want: struct {
				task *Task
				err  error
			}{
				task: nil,
				err:  errors.New("title is required"),
			},
		},
		{
			name: "Fail: description is empty",
			arg: struct {
				userID      string
				title       string
				description string
				dueDate     time.Time
				priority    int
			}{
				userID:      userID,
				title:       "title",
				description: "",
				dueDate:     dueDate,
				priority:    Medium,
			},
			want: struct {
				task *Task
				err  error
			}{
				task: nil,
				err:  errors.New("description is required"),
			},
		},
		{
			name: "Fail: priority is less than 1",
			arg: struct {
				userID      string
				title       string
				description string
				dueDate     time.Time
				priority    int
			}{
				userID:      userID,
				title:       "title",
				description: "description",
				dueDate:     dueDate,
				priority:    0,
			},
			want: struct {
				task *Task
				err  error
			}{
				task: nil,
				err:  errors.New("priority must be between 1 and 5"),
			},
		},
		{
			name: "Fail: priority is greater than 5",
			arg: struct {
				userID      string
				title       string
				description string
				dueDate     time.Time
				priority    int
			}{
				userID:      userID,
				title:       "title",
				description: "description",
				dueDate:     dueDate,
				priority:    6,
			},
			want: struct {
				task *Task
				err  error
			}{
				task: nil,
				err:  errors.New("priority must be between 1 and 5"),
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			task, err := NewTask(tt.arg.userID, tt.arg.title, tt.arg.description, tt.arg.dueDate, tt.arg.priority)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewTask() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewTask() error = %v, wantErr %v", err, tt.want.err)
			}

			if d := cmp.Diff(task, tt.want.task, cmpopts.IgnoreFields(Task{}, "ID", "CreatedAt")); len(d) != 0 {
				t.Errorf("NewTask() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestEntity_Task_CheckOverdue(t *testing.T) {
	t.Parallel()

	now := time.Now()

	patterns := []struct {
		name string
		arg  time.Time
		want bool
	}{
		{
			name: "Success: Not overdue, due tomorrow",
			arg:  now.AddDate(0, 0, 1),
			want: false,
		},
		{
			name: "Success: Already overdue",
			arg:  now.AddDate(0, 0, -1),
			want: true,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			task := &Task{
				DueDate: tt.arg,
			}

			if got := task.CheckOverdue(); got != tt.want {
				t.Errorf("IsOverdue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntity_Task_CheckDueSoon(t *testing.T) {
	t.Parallel()

	now := time.Now()

	patterns := []struct {
		name string
		arg  time.Time
		want bool
	}{
		{
			name: "Success: Due soon within 1 day",
			arg:  now.AddDate(0, 0, 1),
			want: true,
		},
		{
			name: "Success: Not due soon, 2 days left",
			arg:  now.AddDate(0, 0, 2),
			want: false,
		},
		{
			name: "Success: Already overdue",
			arg:  now.AddDate(0, 0, -1),
			want: false,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			task := &Task{
				DueDate: tt.arg,
			}

			if got := task.CheckDueSoon(); got != tt.want {
				t.Errorf("CheckDueSoon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntity_Task_SetPriority(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name string
		arg  int
		want struct {
			priority int
			err      error
		}
	}{
		{
			name: "success",
			arg:  Medium,
			want: struct {
				priority int
				err      error
			}{
				priority: Medium,
				err:      nil,
			},
		},
		{
			name: "Fail: priority is less than 1",
			arg:  0,
			want: struct {
				priority int
				err      error
			}{
				priority: 0,
				err:      errors.New("priority must be between 1 and 5"),
			},
		},
		{
			name: "Fail: priority is greater than 5",
			arg:  6,
			want: struct {
				priority int
				err      error
			}{
				priority: 6,
				err:      errors.New("priority must be between 1 and 5"),
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			task := &Task{}

			err := task.SetPriority(tt.arg)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("SetPriority() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("SetPriority() error = %v, wantErr %v", err, tt.want.err)
			}

			if tt.want.err == nil && task.Priority != tt.want.priority {
				t.Errorf("SetPriority() priority = %v, want %v", task.Priority, tt.want.priority)
			}
		})
	}
}
