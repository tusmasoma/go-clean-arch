package entity

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestEntity_NewUser(t *testing.T) {
	t.Parallel()

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
				password: "password123",
			},
			want: struct {
				user *User
				err  error
			}{
				user: &User{
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
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
				password: "password123",
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
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			getUser, err := NewUser(tt.arg.email, tt.arg.password)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.want.err)
			}

			if d := cmp.Diff(getUser, tt.want.user, cmpopts.IgnoreFields(User{}, "ID")); len(d) != 0 {
				t.Errorf("NewUser() mismatch (-got +want):\n%s", d)
			}
		})
	}
}
