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
	id := uuid.New().String()
	userID := uuid.New().String()
	dueDate := time.Now().AddDate(0, 0, 1).UTC().Truncate(time.Second)
	createdAt := time.Now().UTC().Truncate(time.Second)
	patterns := []struct {
		name string
		arg  struct {
			id          string
			userID      string
			title       string
			description string
			dueDate     time.Time
			priority    int
			createdAt   time.Time
		}
		want struct {
			task *Task
			err  error
		}
	}{
		{
			name: "success",
			arg: struct {
				id          string
				userID      string
				title       string
				description string
				dueDate     time.Time
				priority    int
				createdAt   time.Time
			}{
				id:          id,
				userID:      userID,
				title:       "title",
				description: "description",
				dueDate:     dueDate,
				priority:    int(Medium),
				createdAt:   createdAt,
			},
			want: struct {
				task *Task
				err  error
			}{
				task: &Task{
					ID:          id,
					UserID:      userID,
					Title:       "title",
					Description: "description",
					DueDate:     dueDate,
					CreatedAt:   createdAt,
					Priority:    Medium,
					IsOverdue:   false,
					IsDueSoon:   true,
				},
				err: nil,
			},
		},
		{
			name: "Fail: id is empty",
			arg: struct {
				id          string
				userID      string
				title       string
				description string
				dueDate     time.Time
				priority    int
				createdAt   time.Time
			}{
				id:          "",
				userID:      userID,
				title:       "title",
				description: "description",
				dueDate:     dueDate,
				priority:    int(Medium),
				createdAt:   createdAt,
			},
			want: struct {
				task *Task
				err  error
			}{
				task: nil,
				err:  errors.New("id is required"),
			},
		},
		{
			name: "Fail: userID is empty",
			arg: struct {
				id          string
				userID      string
				title       string
				description string
				dueDate     time.Time
				priority    int
				createdAt   time.Time
			}{
				id:          id,
				userID:      "",
				title:       "title",
				description: "description",
				dueDate:     dueDate,
				priority:    int(Medium),
				createdAt:   createdAt,
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
				id          string
				userID      string
				title       string
				description string
				dueDate     time.Time
				priority    int
				createdAt   time.Time
			}{
				id:          id,
				userID:      userID,
				title:       "",
				description: "description",
				dueDate:     dueDate,
				priority:    int(Medium),
				createdAt:   createdAt,
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
				id          string
				userID      string
				title       string
				description string
				dueDate     time.Time
				priority    int
				createdAt   time.Time
			}{
				id:          id,
				userID:      userID,
				title:       "title",
				description: "",
				dueDate:     dueDate,
				priority:    int(Medium),
				createdAt:   createdAt,
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
				id          string
				userID      string
				title       string
				description string
				dueDate     time.Time
				priority    int
				createdAt   time.Time
			}{
				id:          id,
				userID:      userID,
				title:       "title",
				description: "description",
				dueDate:     dueDate,
				priority:    0,
				createdAt:   createdAt,
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
				id          string
				userID      string
				title       string
				description string
				dueDate     time.Time
				priority    int
				createdAt   time.Time
			}{
				id:          id,
				userID:      userID,
				title:       "title",
				description: "description",
				dueDate:     dueDate,
				priority:    6,
				createdAt:   createdAt,
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
			task, err := NewTask(tt.arg.id, tt.arg.userID, tt.arg.title, tt.arg.description, tt.arg.dueDate, tt.arg.priority, tt.arg.createdAt)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewTask() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewTask() error = %v, wantErr %v", err, tt.want.err)
			}
			if !cmp.Equal(task, tt.want.task) {
				t.Errorf("NewTask() mismatch")
			}
		})
	}
}

func TestEntity_CreateTask(t *testing.T) {
	t.Parallel()
	userID := uuid.New().String()
	dueDate := time.Now().AddDate(0, 0, 1).UTC().Truncate(time.Second)
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
				priority:    int(Medium),
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
					IsDueSoon:   true,
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
				priority:    int(Medium),
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
				priority:    int(Medium),
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
				priority:    int(Medium),
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
			task, err := CreateTask(tt.arg.userID, tt.arg.title, tt.arg.description, tt.arg.dueDate, tt.arg.priority)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.want.err)
			}
			if d := cmp.Diff(task, tt.want.task, cmpopts.IgnoreFields(Task{}, "ID", "CreatedAt")); len(d) != 0 {
				t.Errorf("CreateTask() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestEntity_UpdateTask(t *testing.T) {
	t.Parallel()
	id := uuid.New().String()
	userID := uuid.New().String()
	createdAt := time.Now().UTC().Truncate(time.Second)
	oldDueDate := time.Now().AddDate(0, 0, 1).UTC().Truncate(time.Second)
	newDueDate := time.Now().AddDate(0, 0, 3).UTC().Truncate(time.Second)
	patterns := []struct {
		name     string
		initTask *Task
		arg      struct {
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
			initTask: &Task{
				ID:          id,
				UserID:      userID,
				Title:       "old_title",
				Description: "old_description",
				DueDate:     oldDueDate,
				CreatedAt:   createdAt,
				Priority:    Medium,
			},
			arg: struct {
				title       string
				description string
				dueDate     time.Time
				priority    int
			}{
				title:       "new_title",
				description: "new_description",
				dueDate:     newDueDate,
				priority:    int(High),
			},
			want: struct {
				task *Task
				err  error
			}{
				task: &Task{
					ID:          id,
					UserID:      userID,
					Title:       "new_title",
					Description: "new_description",
					DueDate:     newDueDate,
					CreatedAt:   createdAt,
					Priority:    High,
					IsOverdue:   false,
					IsDueSoon:   false,
				},
				err: nil,
			},
		},
		{
			name: "Fail: title is required",
			initTask: &Task{
				ID:          id,
				UserID:      userID,
				Title:       "old_title",
				Description: "old_description",
				DueDate:     oldDueDate,
				CreatedAt:   createdAt,
				Priority:    Medium,
			},
			arg: struct {
				title       string
				description string
				dueDate     time.Time
				priority    int
			}{
				title:       "",
				description: "new_description",
				dueDate:     newDueDate,
				priority:    int(High),
			},
			want: struct {
				task *Task
				err  error
			}{
				task: &Task{
					ID:          id,
					UserID:      userID,
					Title:       "old_title",
					Description: "old_description",
					DueDate:     oldDueDate,
					CreatedAt:   createdAt,
					Priority:    Medium,
				},
				err: errors.New("title is required"),
			},
		},
		{
			name: "Fail: description is required",
			initTask: &Task{
				ID:          id,
				UserID:      userID,
				Title:       "old_title",
				Description: "old_description",
				DueDate:     oldDueDate,
				CreatedAt:   createdAt,
				Priority:    Medium,
			},
			arg: struct {
				title       string
				description string
				dueDate     time.Time
				priority    int
			}{
				title:       "new_title",
				description: "",
				dueDate:     newDueDate,
				priority:    int(High),
			},
			want: struct {
				task *Task
				err  error
			}{
				task: &Task{
					ID:          id,
					UserID:      userID,
					Title:       "old_title",
					Description: "old_description",
					DueDate:     oldDueDate,
					CreatedAt:   createdAt,
					Priority:    Medium,
				},
				err: errors.New("description is required"),
			},
		},
		{
			name: "Fail: invalid priority (less than 1)",
			initTask: &Task{
				ID:          id,
				UserID:      userID,
				Title:       "old_title",
				Description: "old_description",
				DueDate:     oldDueDate,
				CreatedAt:   createdAt,
				Priority:    Medium,
			},
			arg: struct {
				title       string
				description string
				dueDate     time.Time
				priority    int
			}{
				title:       "new_title",
				description: "new_description",
				dueDate:     newDueDate,
				priority:    0,
			},
			want: struct {
				task *Task
				err  error
			}{
				task: &Task{
					ID:          id,
					UserID:      userID,
					Title:       "old_title",
					Description: "old_description",
					DueDate:     oldDueDate,
					CreatedAt:   createdAt,
					Priority:    Medium,
				},
				err: errors.New("priority must be between 1 and 5"),
			},
		},
		{
			name: "Fail: invalid priority (greater than 5)",
			initTask: &Task{
				ID:          id,
				UserID:      userID,
				Title:       "old_title",
				Description: "old_description",
				DueDate:     oldDueDate,
				CreatedAt:   createdAt,
				Priority:    Medium,
			},
			arg: struct {
				title       string
				description string
				dueDate     time.Time
				priority    int
			}{
				title:       "new_title",
				description: "new_description",
				dueDate:     newDueDate,
				priority:    6,
			},
			want: struct {
				task *Task
				err  error
			}{
				task: &Task{
					ID:          id,
					UserID:      userID,
					Title:       "old_title",
					Description: "old_description",
					DueDate:     oldDueDate,
					CreatedAt:   createdAt,
					Priority:    Medium,
				},
				err: errors.New("priority must be between 1 and 5"),
			},
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.initTask.UpdateTask(tt.arg.title, tt.arg.description, tt.arg.dueDate, tt.arg.priority)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("UpdateTask() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("UpdateTask() error = %v, wantErr %v", err, tt.want.err)
			}
			if d := cmp.Diff(tt.initTask, tt.want.task); len(d) != 0 && errors.Is(err, nil) {
				t.Errorf("UpdateTask() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestEntity_Task_SetOverdue(t *testing.T) {
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
			task.SetOverdue()
			if task.IsOverdue != tt.want {
				t.Errorf("IsOverdue() = %v, want %v", task.IsOverdue, tt.want)
			}
		})
	}
}

func TestEntity_Task_SetkDueSoon(t *testing.T) {
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
			task.SetDueSoon()
			if task.IsDueSoon != tt.want {
				t.Errorf("CheckDueSoon() = %v, want %v", task.IsDueSoon, tt.want)
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
			priority Priority
			err      error
		}
	}{
		{
			name: "success",
			arg:  int(Medium),
			want: struct {
				priority Priority
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
				priority Priority
				err      error
			}{
				priority: Priority(0),
				err:      errors.New("priority must be between 1 and 5"),
			},
		},
		{
			name: "Fail: priority is greater than 5",
			arg:  6,
			want: struct {
				priority Priority
				err      error
			}{
				priority: Priority(6),
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
