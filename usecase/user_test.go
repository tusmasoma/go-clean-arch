package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tusmasoma/go-clean-arch/config"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository/mock"
)

func TestUserUseCase_GetUser(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	ctx := context.WithValue(context.Background(), config.ContextUserIDKey, userID)

	user := entity.User{
		ID:    userID,
		Name:  "test",
		Email: "test@gmail.com",
	}

	patterns := []struct {
		name  string
		ctx   context.Context
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockTransactionRepository,
		)
		wantErr error
	}{
		{
			name: "success",
			ctx:  ctx,
			setup: func(m *mock.MockUserRepository, m1 *mock.MockTransactionRepository) {
				m.EXPECT().Get(
					ctx,
					userID,
				).Return(&user, nil)
			},
			wantErr: nil,
		},
		{
			name:    "Fail: User ID not found in request context",
			ctx:     context.Background(),
			wantErr: errors.New("user name not found in request context"),
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)
			ar := mock.NewMockAuthRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, tr)
			}

			usecase := NewUserUseCase(ur, tr, ar)
			_, err := usecase.GetUser(tt.ctx)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserUseCase_CreateUserAndToken(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockTransactionRepository,
			m2 *mock.MockAuthRepository,
		)
		arg struct {
			ctx      context.Context
			email    string
			password string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockTransactionRepository, m2 *mock.MockAuthRepository) {
				m1.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				m.EXPECT().LockUserByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(false, nil)
				m.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
				).Do(func(_ context.Context, user entity.User) {
					if user.Email != "test@gmail.com" {
						t.Errorf("unexpected Email: got %v, want %v", user.Email, "test@gmail.com")
					}
					if user.Name != "test" {
						t.Errorf("unexpected Name: got %v, want %v", user.Name, "test")
					}
					// TODO: check password hash
				}).Return(nil)
				m2.EXPECT().GenerateToken(
					gomock.Any(),
					"test@gmail.com",
				).Return("jwt", "jti")
			},
			arg: struct {
				ctx      context.Context
				email    string
				password string
			}{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				password: "password123",
			},
			wantErr: nil,
		},
		{
			name: "Fail: Username already exists",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockTransactionRepository, m2 *mock.MockAuthRepository) {
				m1.EXPECT().Transaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				m.EXPECT().LockUserByEmail(
					gomock.Any(),
					"test@gmail.com",
				).Return(true, nil)
			},
			arg: struct {
				ctx      context.Context
				email    string
				password string
			}{
				ctx:      context.Background(),
				email:    "test@gmail.com",
				password: "password123",
			},
			wantErr: errors.New("user with this email already exists"),
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)
			ar := mock.NewMockAuthRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, tr, ar)
			}

			usecase := NewUserUseCase(ur, tr, ar)
			jwt, err := usecase.CreateUserAndToken(tt.arg.ctx, tt.arg.email, tt.arg.password)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CreateUserAndToken() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreateUserAndToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && jwt == "" {
				t.Error("Failed to generate token")
			}
		})
	}
}

func TestUserUseCase_UpdateUser(t *testing.T) {
	t.Parallel()

	userID := uuid.New().String()
	ctx := context.WithValue(context.Background(), config.ContextUserIDKey, userID)

	user := entity.User{
		ID:    userID,
		Name:  "test",
		Email: "test@gmail.com",
	}

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockUserRepository,
			m1 *mock.MockTransactionRepository,
		)
		arg struct {
			ctx  context.Context
			name string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockUserRepository, m1 *mock.MockTransactionRepository) {
				m.EXPECT().Get(
					ctx,
					userID,
				).Return(&user, nil)
				user.Name = "updatedName"
				m.EXPECT().Update(
					gomock.Any(),
					user,
				).Return(nil)
			},
			arg: struct {
				ctx  context.Context
				name string
			}{
				ctx:  ctx,
				name: "updatedName",
			},
			wantErr: nil,
		},
		{
			name: "Fail: User ID not found in request context",
			arg: struct {
				ctx  context.Context
				name string
			}{
				ctx:  context.Background(),
				name: "updatedName",
			},
			wantErr: errors.New("user name not found in request context"),
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ur := mock.NewMockUserRepository(ctrl)
			tr := mock.NewMockTransactionRepository(ctrl)
			ar := mock.NewMockAuthRepository(ctrl)

			if tt.setup != nil {
				tt.setup(ur, tr)
			}

			usecase := NewUserUseCase(ur, tr, ar)
			err := usecase.UpdateUser(tt.arg.ctx, tt.arg.name)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
