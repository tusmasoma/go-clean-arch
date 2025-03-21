package email

import (
	"errors"
	"testing"
)

func Test_GetAddressPart(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		name string
		arg  struct {
			email string
		}
		want struct {
			address string
			err     error
		}
	}{
		{
			name: "success",
			arg: struct {
				email string
			}{
				email: "test@gmail.com",
			},
			want: struct {
				address string
				err     error
			}{
				address: "test",
				err:     nil,
			},
		},
		{
			name: "Fail: email is required",
			arg: struct {
				email string
			}{
				email: "",
			},
			want: struct {
				address string
				err     error
			}{
				address: "",
				err:     errors.New("invalid email address: missing '@' or multiple '@'"),
			},
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotAddress, err := GetAddressPart(tt.arg.email)
			if gotAddress != tt.want.address {
				t.Errorf("GetAddressPart gotAddress = %v, wantName %v", gotAddress, tt.want.address)
			}
			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("GetAddressPart error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("GetAddressPart error = %v, wantErr %v", err, tt.want.err)
			}
		})
	}
}
