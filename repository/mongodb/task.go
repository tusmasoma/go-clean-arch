package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type taskModel struct {
	ID          string    `bson:"_id,omitempty"`
	UserID      string    `bson:"user_id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	DueDate     time.Time `bson:"duedate"`
	Priority    int       `bson:"priority"`
	CreatedAt   time.Time `bson:"created_at"`
}

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

	var tm taskModel
	if err := collection.FindOne(ctx, filter).Decode(&tm); err != nil {
		return nil, err
	}
	return &entity.Task{
		ID:          tm.ID,
		UserID:      tm.UserID,
		Title:       tm.Title,
		Description: tm.Description,
		DueDate:     tm.DueDate,
		Priority:    tm.Priority,
		CreatedAt:   tm.CreatedAt,
	}, nil
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

	var tms []taskModel
	if err = cursor.All(ctx, &tms); err != nil {
		return nil, err
	}

	tasks := make([]entity.Task, len(tms))
	for i, tm := range tms {
		tasks[i] = entity.Task{
			ID:          tm.ID,
			UserID:      tm.UserID,
			Title:       tm.Title,
			Description: tm.Description,
			DueDate:     tm.DueDate,
			Priority:    tm.Priority,
			CreatedAt:   tm.CreatedAt,
		}
	}
	return tasks, nil
}

func (tr *taskRepository) Create(ctx context.Context, task entity.Task) error {
	collection := tr.client.cli.Database(tr.client.db).Collection(tr.table)

	tm := taskModel{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt,
	}

	if _, err := collection.InsertOne(ctx, tm); err != nil {
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
