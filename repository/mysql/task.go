package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type taskModel struct {
	ID          string    `db:"id"`
	UserID      string    `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	DueDate     time.Time `db:"duedate"`
	Priority    int       `db:"priority"`
	CreatedAt   time.Time `db:"created_at"`
}

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
	WHERE id = ?
	LIMIT 1
	`

	row := executor.QueryRowContext(ctx, query, id)

	var tm taskModel
	if err := row.Scan(
		&tm.ID,
		&tm.UserID,
		&tm.Title,
		&tm.Description,
		&tm.DueDate,
		&tm.Priority,
		&tm.CreatedAt,
	); err != nil {
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

func (ur *taskRepository) List(ctx context.Context, userID string) ([]entity.Task, error) {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `SELECT *
	FROM Tasks
	WHERE user_id = ?
	`

	rows, err := executor.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tms []taskModel
	for rows.Next() {
		var tm taskModel
		if err = rows.Scan(
			&tm.ID,
			&tm.UserID,
			&tm.Title,
			&tm.Description,
			&tm.DueDate,
			&tm.Priority,
			&tm.CreatedAt,
		); err != nil {
			return nil, err
		}
		tms = append(tms, tm)
	}
	if err = rows.Err(); err != nil {
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

func (ur *taskRepository) Create(ctx context.Context, task entity.Task) error {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `INSERT INTO Tasks (
	id, user_id, title, description, duedate, priority, created_at
	)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	tm := taskModel{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt,
	}

	if _, err := executor.ExecContext(
		ctx,
		query,
		tm.ID,
		tm.UserID,
		tm.Title,
		tm.Description,
		tm.DueDate,
		tm.Priority,
		tm.CreatedAt,
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
	SET title = ?, description = ?, duedate = ?, priority = ?
	WHERE id = ?
	`

	tm := taskModel{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Priority:    task.Priority,
	}

	if _, err := executor.ExecContext(
		ctx,
		query,
		tm.Title,
		tm.Description,
		tm.DueDate,
		tm.Priority,
		tm.ID,
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
	WHERE id = ?
	`

	if _, err := executor.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}
