package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/tusmasoma/go-clean-arch/entity"
)

func Test_TaskRepository(t *testing.T) {
	ctx := context.Background()

	if client == nil {
		t.Skip("MongoDB is not available")
	}

	var cli Client
	cli.cli = client
	cli.db = "goCleanArcTestDB"
	repo := NewTaskRepository(&cli)

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
	if d := cmp.Diff(task1, gottask, cmpopts.IgnoreFields(entity.Task{}, "DueDate", "CreatedAt")); len(d) != 0 {
		t.Errorf("differs: (-want +got)\n%s", d)
	}

	// List
	gottasks, err := repo.List(ctx)
	ValidateErr(t, err, nil)
	if len(gottasks) != 2 {
		t.Errorf("want: %v, got: %v", 2, len(gottasks))
	}

	// Update
	gottask.Title = "Updated First Task"
	err = repo.Update(ctx, *gottask)
	ValidateErr(t, err, nil)

	updatedtask, err := repo.Get(ctx, task1.ID)
	ValidateErr(t, err, nil)
	if d := cmp.Diff(gottask, updatedtask, cmpopts.IgnoreFields(entity.Task{}, "CreatedAt")); len(d) != 0 {
		t.Errorf("differs: (-want +got)\n%s", d)
	}

	// Delete
	err = repo.Delete(ctx, task1.ID)
	ValidateErr(t, err, nil)

	_, err = repo.Get(ctx, task1.ID)
	if err == nil {
		t.Errorf("want: %v, got: %v", nil, err)
	}
}
