package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type taskRepository struct {
	client *Client
	table  string
}

func NewTaskRepository(client *Client) repository.TaskRepository {
	return &taskRepository{
		client: client,
		table:  "Tasks",
	}
}

func (tr *taskRepository) Get(ctx context.Context, id string) (*entity.Task, error) {
	filter := bson.M{"_id": id}

	collection := tr.client.cli.Database(tr.client.db).Collection(tr.table)

	var task entity.Task
	if err := collection.FindOne(ctx, filter).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (tr *taskRepository) List(ctx context.Context, userID string) ([]entity.Task, error) {
	collection := tr.client.cli.Database(tr.client.db).Collection(tr.table)
	filter := bson.M{
		"user_id": userID,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []entity.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (tr *taskRepository) Create(ctx context.Context, task entity.Task) error {
	collection := tr.client.cli.Database(tr.client.db).Collection(tr.table)

	if _, err := collection.InsertOne(ctx, task); err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) Update(ctx context.Context, task entity.Task) error {
	collection := tr.client.cli.Database(tr.client.db).Collection(tr.table)

	filter := bson.M{"_id": task.ID}

	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"duedate":     task.DueDate,
			"priority":    task.Priority,
			"created_at":  task.CreatedAt,
		},
	}

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) Delete(ctx context.Context, id string) error {
	collection := tr.client.cli.Database(tr.client.db).Collection(tr.table)

	filter := bson.M{"_id": id}

	if _, err := collection.DeleteOne(ctx, filter); err != nil {
		return err
	}
	return nil
}
