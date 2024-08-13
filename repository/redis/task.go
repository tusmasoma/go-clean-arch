package redis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/tusmasoma/go-tech-dojo/config"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type taskRepository struct {
	client *redis.Client
}

func NewTaskRepository(client *redis.Client) repository.TaskRepository {
	return &taskRepository{
		client: client,
	}
}

func (tr *taskRepository) Get(ctx context.Context, id string) (*entity.Task, error) {
	val, err := tr.client.Get(ctx, id).Result()
	if errors.Is(err, redis.Nil) {
		log.Warn("Cache miss", log.Fstring("key", id))
		return nil, config.ErrCacheMiss
	} else if err != nil {
		log.Error("Failed to get cache", log.Ferror(err))
		return nil, err
	}
	tasks, err := tr.deserialize(val)
	if err != nil {
		log.Error("Failed to deserialize task", log.Ferror(err))
		return nil, err
	}
	log.Info("Cache hit", log.Fstring("key", id))
	return tasks, nil
}

func (tr *taskRepository) List(ctx context.Context) ([]entity.Task, error) {
	ids, err := tr.client.Keys(ctx, "*").Result()
	if err != nil {
		log.Error("Failed to get keys", log.Ferror(err))
		return nil, err
	}
	var tasks []entity.Task
	for _, id := range ids {
		val, err := tr.client.Get(ctx, id).Result() //nolint: govet // This is a false positive
		if err != nil {
			log.Error("Failed to get cache", log.Ferror(err))
			return nil, err
		}
		task, err := tr.deserialize(val)
		if err != nil {
			log.Error("Failed to deserialize task", log.Ferror(err))
			return nil, err
		}
		tasks = append(tasks, *task)
		log.Info("Cache hit", log.Fstring("key", id))
	}
	return tasks, nil
}

func (tr *taskRepository) Create(ctx context.Context, task entity.Task) error {
	serializeTask, err := tr.serialize(task)
	if err != nil {
		log.Error("Failed to serialize Task", log.Ferror(err))
		return err
	}
	if err = tr.client.Set(ctx, task.ID, serializeTask, 0).Err(); err != nil {
		log.Error("Failed to set cache", log.Ferror(err))
		return err
	}
	log.Info("Cache set successfully", log.Fstring("key", task.ID))
	return nil
}

func (tr *taskRepository) Update(ctx context.Context, task entity.Task) error {
	serializeTask, err := tr.serialize(task)
	if err != nil {
		log.Error("Failed to serialize Task", log.Ferror(err))
		return err
	}
	if err = tr.client.Set(ctx, task.ID, serializeTask, 0).Err(); err != nil {
		log.Error("Failed to set cache", log.Ferror(err))
		return err
	}
	log.Info("Cache updated successfully", log.Fstring("key", task.ID))
	return nil
}

func (tr *taskRepository) Delete(ctx context.Context, id string) error {
	if err := tr.client.Del(ctx, id).Err(); err != nil {
		log.Error("Failed to delete cache", log.Ferror(err))
		return err
	}
	log.Info("Cache deleted successfully", log.Fstring("key", id))
	return nil
}

func (tr *taskRepository) serialize(tasks entity.Task) (string, error) {
	data, err := json.Marshal(tasks)
	if err != nil {
		log.Error("Failed to serialize tasks", log.Ferror(err))
		return "", err
	}
	return string(data), nil
}

func (tr *taskRepository) deserialize(data string) (*entity.Task, error) {
	var task entity.Task
	if err := json.Unmarshal([]byte(data), &task); err != nil {
		log.Error("Failed to deserialize tasks", log.Ferror(err))
		return nil, err
	}
	return &task, nil
}
