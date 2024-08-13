package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository/mock"
)

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
