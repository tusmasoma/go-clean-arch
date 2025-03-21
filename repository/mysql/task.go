package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/pkg/log"
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
	db DB
}

func NewTaskRepository(db *sql.DB) repository.TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (tr *taskRepository) Get(ctx context.Context, id string) (*entity.Task, error) {
	query := `SELECT *
	FROM Tasks
	WHERE id = ?
	LIMIT 1
	`
	row := tr.db.QueryRowContext(ctx, query, id)
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
	task, err := entity.NewTask(
		tm.ID,
		tm.UserID,
		tm.Title,
		tm.Description,
		tm.DueDate,
		tm.Priority,
		tm.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (tr *taskRepository) List(ctx context.Context, userID string) ([]entity.Task, error) {
	query := `SELECT *
	FROM Tasks
	WHERE user_id = ?
	`
	rows, err := tr.db.QueryContext(ctx, query, userID)
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
		var task *entity.Task
		task, err = entity.NewTask(
			tm.ID,
			tm.UserID,
			tm.Title,
			tm.Description,
			tm.DueDate,
			tm.Priority,
			tm.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks[i] = *task
	}
	return tasks, nil
}

func (tr *taskRepository) Create(ctx context.Context, task entity.Task) error {
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
		Priority:    int(task.Priority),
		CreatedAt:   task.CreatedAt,
	}
	if _, err := tr.db.ExecContext(
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

func (tr *taskRepository) Update(ctx context.Context, task entity.Task) error {
	query := `UPDATE Tasks
	SET title = ?, description = ?, duedate = ?, priority = ?
	WHERE id = ?
	`
	tm := taskModel{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Priority:    int(task.Priority),
	}
	if _, err := tr.db.ExecContext(
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

func (tr *taskRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM Tasks
	WHERE id = ?
	`
	if _, err := tr.db.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) transaction(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error { //nolint: unused // ignore unused
	tx, err := tr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil || err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error("Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()
	if err = fn(ctx, tx); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
