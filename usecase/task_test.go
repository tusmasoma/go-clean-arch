package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/go-tech-dojo/config"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository/mock"
)

func TestUseCase_GetTask(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	ctx := context.WithValue(context.Background(), config.ContextUserIDKey, userID)

	taskID := uuid.New().String()
	dueDate := time.Now().AddDate(0, 0, 1)

	task := &entity.Task{
		ID:          taskID,
		UserID:      userID,
		Title:       "title",
		Description: "description",
		DueDate:     dueDate,
		Priority:    3,
		CreatedAt:   time.Now(),
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTaskRepository,
		)
		arg struct {
			ctx context.Context
			id  string
		}
		want struct {
			task *entity.Task
			err  error
		}
	}{
		{
			name: "success",
			setup: func(tr *mock.MockTaskRepository) {
				tr.EXPECT().Get(gomock.Any(), taskID).Return(task, nil)
			},
			arg: struct {
				ctx context.Context
				id  string
			}{
				ctx: ctx,
				id:  taskID,
			},
			want: struct {
				task *entity.Task
				err  error
			}{
				task: task,
				err:  nil,
			},
		},
		{
			name: "Fail: Task does not belong to the user",
			setup: func(tr *mock.MockTaskRepository) {
				tr.EXPECT().Get(gomock.Any(), taskID).Return(&entity.Task{
					ID:     taskID,
					UserID: uuid.New().String(),
				}, nil)
			},
			arg: struct {
				ctx context.Context
				id  string
			}{
				ctx: ctx,
				id:  taskID,
			},
			want: struct {
				task *entity.Task
				err  error
			}{
				task: nil,
				err:  errors.New("task does not belong to the user"),
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tr := mock.NewMockTaskRepository(ctrl)

			if tt.setup != nil {
				tt.setup(tr)
			}

			tuc := NewTaskUseCase(tr)

			getTask, err := tuc.GetTask(tt.arg.ctx, tt.arg.id)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("GetTask() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("GetTask() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(getTask, tt.want.task) {
				t.Errorf("GetTask() got = %v, want %v", getTask, tt.want.task)
			}
		})
	}
}

func TestUseCase_ListTasks(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	ctx := context.WithValue(context.Background(), config.ContextUserIDKey, userID)
	dueDate := time.Now().AddDate(0, 0, 1)

	tasks := []entity.Task{
		{
			ID:          uuid.New().String(),
			UserID:      userID,
			Title:       "title",
			Description: "description",
			DueDate:     dueDate,
			Priority:    3,
			CreatedAt:   time.Now(),
		},
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTaskRepository,
		)
		arg struct {
			ctx context.Context
		}
		want struct {
			tasks []entity.Task
			err   error
		}
	}{
		{
			name: "success",
			setup: func(tr *mock.MockTaskRepository) {
				tr.EXPECT().List(
					gomock.Any(),
					userID,
				).Return(tasks, nil)
			},
			arg: struct {
				ctx context.Context
			}{
				ctx: ctx,
			},
			want: struct {
				tasks []entity.Task
				err   error
			}{
				tasks: tasks,
				err:   nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tr := mock.NewMockTaskRepository(ctrl)

			if tt.setup != nil {
				tt.setup(tr)
			}

			tuc := NewTaskUseCase(tr)

			getTasks, err := tuc.ListTasks(tt.arg.ctx)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("ListTasks() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("ListTasks() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(getTasks, tt.want.tasks) {
				t.Errorf("ListTasks() got = %v, want %v", getTasks, tt.want.tasks)
			}
		})
	}
}

func TestUseCase_CreateTask(t *testing.T) { //nolint: gocognit // The complexity is caused by the test patterns
	t.Parallel()

	userID := uuid.New().String()
	ctx := context.WithValue(context.Background(), config.ContextUserIDKey, userID)
	dueDate := time.Now().AddDate(0, 0, 1)

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTaskRepository,
		)
		arg struct {
			ctx    context.Context
			params *CreateTaskParams
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(tr *mock.MockTaskRepository) {
				tr.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
				).Do(func(_ context.Context, task entity.Task) {
					if task.UserID != userID {
						t.Errorf("unexpected UserID: got %v, want %v", task.UserID, userID)
					}
					if task.Title != "title" {
						t.Errorf("unexpected Title: got %v, want %v", task.Title, "title")
					}
					if task.Description != "description" {
						t.Errorf("unexpected Description: got %v, want %v", task.Description, "description")
					}
					if !task.DueDate.Equal(dueDate) {
						t.Errorf("unexpected DueDate: got %v, want %v", task.DueDate, dueDate)
					}
					if task.Priority != 3 {
						t.Errorf("unexpected Priority: got %v, want %v", task.Priority, 3)
					}
				}).Return(nil)
			},
			arg: struct {
				ctx    context.Context
				params *CreateTaskParams
			}{
				ctx: ctx,
				params: &CreateTaskParams{
					Title:       "title",
					Description: "description",
					DueDate:     dueDate,
					Priority:    3,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tr := mock.NewMockTaskRepository(ctrl)

			if tt.setup != nil {
				tt.setup(tr)
			}

			tuc := NewTaskUseCase(tr)

			err := tuc.CreateTask(tt.arg.ctx, tt.arg.params)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestUseCase_UpdateTask(t *testing.T) { //nolint: gocognit // The complexity is caused by the test patterns
	t.Parallel()

	userID := uuid.New().String()
	ctx := context.WithValue(context.Background(), config.ContextUserIDKey, userID)
	taskID := uuid.New().String()
	dueDate := time.Now().AddDate(0, 0, 1)

	task := &entity.Task{
		ID:          taskID,
		UserID:      userID,
		Title:       "title",
		Description: "description",
		DueDate:     dueDate,
		Priority:    3,
		CreatedAt:   time.Now(),
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTaskRepository,
		)
		arg struct {
			ctx    context.Context
			params *UpdateTaskParams
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(tr *mock.MockTaskRepository) {
				tr.EXPECT().Get(
					gomock.Any(),
					taskID,
				).Return(task, nil)
				tr.EXPECT().Update(
					gomock.Any(),
					gomock.Any(),
				).Do(func(_ context.Context, task entity.Task) {
					if task.UserID != userID {
						t.Errorf("unexpected UserID: got %v, want %v", task.UserID, userID)
					}
					if task.Title != "updated title" {
						t.Errorf("unexpected Title: got %v, want %v", task.Title, "updated title")
					}
					if task.Description != "updated description" {
						t.Errorf("unexpected Description: got %v, want %v", task.Description, "updated description")
					}
					if !task.DueDate.Equal(dueDate) {
						t.Errorf("unexpected DueDate: got %v, want %v", task.DueDate, dueDate)
					}
					if task.Priority != 2 {
						t.Errorf("unexpected Priority: got %v, want %v", task.Priority, 2)
					}
				}).Return(nil)
			},
			arg: struct {
				ctx    context.Context
				params *UpdateTaskParams
			}{
				ctx: ctx,
				params: &UpdateTaskParams{
					ID:          taskID,
					Title:       "updated title",
					Description: "updated description",
					DueDate:     dueDate,
					Priority:    2,
				},
			},
			wantErr: nil,
		},
		{
			name: "Fail: Task does not belong to the user",
			setup: func(tr *mock.MockTaskRepository) {
				tr.EXPECT().Get(
					gomock.Any(),
					taskID,
				).Return(&entity.Task{
					ID:     taskID,
					UserID: uuid.New().String(),
				}, nil)
			},
			arg: struct {
				ctx    context.Context
				params *UpdateTaskParams
			}{
				ctx: ctx,
				params: &UpdateTaskParams{
					ID:          taskID,
					Title:       "updated title",
					Description: "updated description",
					DueDate:     dueDate,
					Priority:    2,
				},
			},
			wantErr: errors.New("task does not belong to the user"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tr := mock.NewMockTaskRepository(ctrl)

			if tt.setup != nil {
				tt.setup(tr)
			}

			tuc := NewTaskUseCase(tr)

			err := tuc.UpdateTask(tt.arg.ctx, tt.arg.params)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsaCase_DeleteTask(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	ctx := context.WithValue(context.Background(), config.ContextUserIDKey, userID)
	taskID := uuid.New().String()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockTaskRepository,
		)
		arg struct {
			ctx context.Context
			id  string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(tr *mock.MockTaskRepository) {
				tr.EXPECT().Get(
					gomock.Any(),
					taskID,
				).Return(&entity.Task{
					ID:     taskID,
					UserID: userID,
				}, nil)
				tr.EXPECT().Delete(gomock.Any(), taskID).Return(nil)
			},
			arg: struct {
				ctx context.Context
				id  string
			}{
				ctx: ctx,
				id:  taskID,
			},
			wantErr: nil,
		},
		{
			name: "Fail: Task does not belong to the user",
			setup: func(tr *mock.MockTaskRepository) {
				tr.EXPECT().Get(
					gomock.Any(),
					taskID,
				).Return(&entity.Task{
					ID:     taskID,
					UserID: uuid.New().String(),
				}, nil)
			},
			arg: struct {
				ctx context.Context
				id  string
			}{
				ctx: ctx,
				id:  taskID,
			},
			wantErr: errors.New("task does not belong to the user"),
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tr := mock.NewMockTaskRepository(ctrl)

			if tt.setup != nil {
				tt.setup(tr)
			}

			tuc := NewTaskUseCase(tr)

			err := tuc.DeleteTask(tt.arg.ctx, tt.arg.id)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
