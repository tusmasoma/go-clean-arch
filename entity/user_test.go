package entity

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestEntity_NewUser(t *testing.T) {
	t.Parallel()
	id := uuid.New().String()
	hashedPasswordByte, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	hashedPassword := string(hashedPasswordByte)
	patterns := []struct {
		name string
		arg  struct {
			id           string
			name         string
			email        string
			passwordHash string
		}
		want struct {
			user *User
			err  error
		}
	}{
		{
			name: "success",
			arg: struct {
				id           string
				name         string
				email        string
				passwordHash string
			}{
				id:           id,
				name:         "test",
				email:        "test@gmail.com",
				passwordHash: hashedPassword,
			},
			want: struct {
				user *User
				err  error
			}{
				user: &User{
					ID:           id,
					Name:         "test",
					Email:        "test@gmail.com",
					PasswordHash: hashedPassword,
				},
				err: nil,
			},
		},
		{
			name: "Fail: id is required",
			arg: struct {
				id           string
				name         string
				email        string
				passwordHash string
			}{
				id:           "",
				name:         "test",
				email:        "test@gmail.com",
				passwordHash: hashedPassword,
			},
			want: struct {
				user *User
				err  error
			}{
				user: nil,
				err:  errors.New("id is required"),
			},
		},
		{
			name: "Fail: name is required",
			arg: struct {
				id           string
				name         string
				email        string
				passwordHash string
			}{
				id:           id,
				name:         "",
				email:        "test@gmail.com",
				passwordHash: hashedPassword,
			},
			want: struct {
				user *User
				err  error
			}{
				user: nil,
				err:  errors.New("name is required"),
			},
		},
		{
			name: "Fail: email is required",
			arg: struct {
				id           string
				name         string
				email        string
				passwordHash string
			}{
				id:           id,
				name:         "test",
				email:        "",
				passwordHash: hashedPassword,
			},
			want: struct {
				user *User
				err  error
			}{
				user: nil,
				err:  errors.New("email is required"),
			},
		},
		{
			name: "Fail: password is required",
			arg: struct {
				id           string
				name         string
				email        string
				passwordHash string
			}{
				id:           id,
				name:         "test",
				email:        "test@gmail.com",
				passwordHash: "",
			},
			want: struct {
				user *User
				err  error
			}{
				user: nil,
				err:  errors.New("hash password is required"),
			},
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			getUser, err := NewUser(tt.arg.id, tt.arg.name, tt.arg.email, tt.arg.passwordHash)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.want.err)
			}
			if !cmp.Equal(getUser, tt.want.user) {
				t.Errorf("NewUser() mismatch")
			}
		})
	}
}

func TestEntity_CreateUser(t *testing.T) {
	t.Parallel()
	password := "password123"
	patterns := []struct {
		name string
		arg  struct {
			email    string
			password string
		}
		want struct {
			user *User
			err  error
		}
	}{
		{
			name: "success",
			arg: struct {
				email    string
				password string
			}{
				email:    "test@gmail.com",
				password: password,
			},
			want: struct {
				user *User
				err  error
			}{
				user: &User{
					Name:  "test",
					Email: "test@gmail.com",
				},
				err: nil,
			},
		},
		{
			name: "Fail: email is required",
			arg: struct {
				email    string
				password string
			}{
				email:    "",
				password: password,
			},
			want: struct {
				user *User
				err  error
			}{
				user: nil,
				err:  errors.New("email is required"),
			},
		},
		{
			name: "Fail: password is required",
			arg: struct {
				email    string
				password string
			}{
				email:    "test@gmail.com",
				password: "",
			},
			want: struct {
				user *User
				err  error
			}{
				user: nil,
				err:  errors.New("password is required"),
			},
		},
		{
			name: "Fail: get name from email",
			arg: struct {
				email    string
				password string
			}{
				email:    "gmail.com",
				password: password,
			},
			want: struct {
				user *User
				err  error
			}{
				user: nil,
				err:  errors.New("email is required"),
			},
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			getUser, err := CreateUser(tt.arg.email, tt.arg.password)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.want.err)
			}
			if d := cmp.Diff(getUser, tt.want.user, cmpopts.IgnoreFields(User{}, "ID", "PasswordHash")); len(d) != 0 {
				t.Errorf("CreateUser() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestEntity_UpdateUser(t *testing.T) {
	t.Parallel()
	id := uuid.New().String()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	patterns := []struct {
		name     string
		initUser *User
		arg      struct {
			name  string
			email string
		}
		want struct {
			user *User
			err  error
		}
	}{
		{
			name: "success",
			initUser: &User{
				ID:           id,
				Name:         "old_name",
				Email:        "old@example.com",
				PasswordHash: string(hashedPassword),
			},
			arg: struct {
				name  string
				email string
			}{
				name:  "new_name",
				email: "new@example.com",
			},
			want: struct {
				user *User
				err  error
			}{
				user: &User{
					ID:           id,
					Name:         "new_name",
					Email:        "new@example.com",
					PasswordHash: string(hashedPassword),
				},
				err: nil,
			},
		},
		{
			name: "Fail: name is required",
			initUser: &User{
				ID:           id,
				Name:         "old_name",
				Email:        "old@example.com",
				PasswordHash: string(hashedPassword),
			},
			arg: struct {
				name  string
				email string
			}{
				name:  "",
				email: "new@example.com",
			},
			want: struct {
				user *User
				err  error
			}{
				user: &User{
					ID:           id,
					Name:         "old_name",
					Email:        "old@example.com",
					PasswordHash: string(hashedPassword),
				},
				err: errors.New("name is required"),
			},
		},
		{
			name: "Fail: email is required",
			initUser: &User{
				ID:           id,
				Name:         "old_name",
				Email:        "old@example.com",
				PasswordHash: string(hashedPassword),
			},
			arg: struct {
				name  string
				email string
			}{
				name:  "new_name",
				email: "",
			},
			want: struct {
				user *User
				err  error
			}{
				user: &User{
					ID:           id,
					Name:         "old_name",
					Email:        "old@example.com",
					PasswordHash: string(hashedPassword),
				},
				err: errors.New("email is required"),
			},
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.initUser.UpdateUser(tt.arg.name, tt.arg.email)
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.want.err)
			}
			if d := cmp.Diff(tt.initUser, tt.want.user); len(d) != 0 {
				t.Errorf("UpdateUser() mismatch (-got +want):\n%s", d)
			}
		})
	}
}
