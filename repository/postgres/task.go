package postgres

import (
	"context"
	"database/sql"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type taskRepository struct {
	db SQLExecutor
}

func NewTaskRepository(db *sql.DB) repository.TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (ur *taskRepository) Get(ctx context.Context, id string) (*entity.Task, error) {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `SELECT *
	FROM Tasks
	WHERE id = $1
	LIMIT 1
	`

	row := executor.QueryRowContext(ctx, query, id)

	var task entity.Task
	if err := row.Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Priority,
		&task.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &task, nil
}

func (ur *taskRepository) List(ctx context.Context, userID string) ([]entity.Task, error) {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `SELECT *
	FROM Tasks
	WHERE user_id = $1
	`

	rows, err := executor.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		if err = rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.Priority,
			&task.CreatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (ur *taskRepository) Create(ctx context.Context, task entity.Task) error {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `INSERT INTO Tasks (
	id, user_id, title, description, duedate, priority, created_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := executor.ExecContext(
		ctx,
		query,
		task.ID,
		task.UserID,
		task.Title,
		task.Description,
		task.DueDate,
		task.Priority,
		task.CreatedAt,
	); err != nil {
		return err
	}
	return nil
}

func (ur *taskRepository) Update(ctx context.Context, task entity.Task) error {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `UPDATE Tasks
	SET title = $1, description = $2, duedate = $3, priority = $4
	WHERE id = $5
	`

	if _, err := executor.ExecContext(
		ctx,
		query,
		task.Title,
		task.Description,
		task.DueDate,
		task.Priority,
		task.ID,
	); err != nil {
		return err
	}
	return nil
}

func (ur *taskRepository) Delete(ctx context.Context, id string) error {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `DELETE FROM Tasks
	WHERE id = $1
	`

	if _, err := executor.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}
