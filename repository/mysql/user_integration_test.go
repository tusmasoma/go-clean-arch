package mysql

import (
	"context"
	"reflect"
	"testing"

	"github.com/tusmasoma/go-clean-arch/entity"
)

func Test_UserRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository(db)

	user, err := entity.NewUser(
		"test@gmail.com",
		"password",
	)
	ValidateErr(t, err, nil)

	// Create
	err = repo.Create(ctx, *user)
	ValidateErr(t, err, nil)

	// Get
	gotUser, err := repo.Get(ctx, user.ID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(user, gotUser) {
		t.Errorf("want: %v, got: %v", user, gotUser)
	}

	// LockUserByEmail
	exists, err := repo.LockUserByEmail(ctx, "test@gmail.com")
	ValidateErr(t, err, nil)
	if !exists {
		t.Fatalf("Failed to get user by email")
	}

	// Update
	gotUser.Name = "updatedName"
	err = repo.Update(ctx, *gotUser)
	ValidateErr(t, err, nil)

	updatedUser, err := repo.Get(ctx, user.ID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(gotUser, updatedUser) {
		t.Errorf("want: %v, got: %v", gotUser, updatedUser)
	}

	// Delete
	err = repo.Delete(ctx, user.ID)
	ValidateErr(t, err, nil)

	_, err = repo.Get(ctx, user.ID)
	if err == nil {
		t.Errorf("want: %v, got: %v", nil, err)
	}
}
