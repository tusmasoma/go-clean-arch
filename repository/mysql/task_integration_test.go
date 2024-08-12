package mysql

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/tusmasoma/go-clean-arch/entity"
)

func Test_TaskRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewTaskRepository(db)

	task1, err := entity.NewTask(
		"First Task",
		"First Description",
		time.Now().Add(24*time.Hour),
		3,
	)
	ValidateErr(t, err, nil)
	task2, err := entity.NewTask(
		"Second Task",
		"Second Description",
		time.Now().Add(48*time.Hour),
		4,
	)
	ValidateErr(t, err, nil)

	// Create
	err = repo.Create(ctx, *task1)
	ValidateErr(t, err, nil)
	err = repo.Create(ctx, *task2)
	ValidateErr(t, err, nil)

	// Get
	gottask, err := repo.Get(ctx, task1.ID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(task1, gottask) {
		t.Errorf("want: %v, got: %v", task1, gottask)
	}

	// List
	gottasks, err := repo.List(ctx)
	ValidateErr(t, err, nil)
	if len(gottasks) != 2 {
		t.Errorf("want: %v, got: %v", 2, len(gottasks))
	}

	// Update
	gottask.Title = "Updated First Task"
	gottask.DueData = time.Now().Add(48 * time.Hour)
	err = repo.Update(ctx, *gottask)
	ValidateErr(t, err, nil)

	updatedtask, err := repo.Get(ctx, task1.ID)
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(gottask, updatedtask) {
		t.Errorf("want: %v, got: %v", gottask, updatedtask)
	}

	// Delete
	err = repo.Delete(ctx, task1.ID)
	ValidateErr(t, err, nil)

	_, err = repo.Get(ctx, task1.ID)
	if err == nil {
		t.Errorf("want: %v, got: %v", nil, err)
	}
}
