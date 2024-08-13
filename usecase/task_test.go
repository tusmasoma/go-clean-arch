package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/go-clean-arch/repository/mock"
)

func TestUseCase_CreateTask(t *testing.T) {
	t.Parallel()

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
				).Return(nil)
			},
			arg: struct {
				ctx    context.Context
				params *CreateTaskParams
			}{
				ctx: context.Background(),
				params: &CreateTaskParams{
					Title:       "title",
					Description: "description",
					DueData:     time.Now().AddDate(0, 0, 1),
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
