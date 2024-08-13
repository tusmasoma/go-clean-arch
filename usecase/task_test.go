package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository/mock"
)

func TestUseCase_GetTask(t *testing.T) {
	t.Parallel()

	taskID := uuid.New().String()
	dueDate := time.Now().AddDate(0, 0, 1)

	task := &entity.Task{
		ID:          taskID,
		Title:       "title",
		Description: "description",
		DueData:     dueDate,
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
				ctx: context.Background(),
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
				t.Errorf("GetTask() got = %v, want %v", task, tt.want.task)
			}
		})
	}
}

func TestUseCase_CreateTask(t *testing.T) {
	t.Parallel()

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
					if task.Title != "title" {
						t.Errorf("unexpected Title: got %v, want %v", task.Title, "title")
					}
					if task.Description != "description" {
						t.Errorf("unexpected Description: got %v, want %v", task.Description, "description")
					}
					if !task.DueData.Equal(dueDate) {
						t.Errorf("unexpected DueData: got %v, want %v", task.DueData, dueDate)
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
				ctx: context.Background(),
				params: &CreateTaskParams{
					Title:       "title",
					Description: "description",
					DueData:     dueDate,
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
