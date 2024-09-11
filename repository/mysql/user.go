package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type userRepository struct {
	db SQLExecutor
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) Get(ctx context.Context, id string) (*entity.User, error) {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `SELECT *
	FROM Users
	WHERE id = ?
	LIMIT 1`

	row := executor.QueryRowContext(ctx, query, id)

	var user entity.User
	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) Create(ctx context.Context, user entity.User) error {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `INSERT INTO Users (
	id, name, email, password
	)
	VALUES (?, ?, ?, ?)
	`

	if _, err := executor.ExecContext(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
	); err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) Update(ctx context.Context, user entity.User) error {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `UPDATE Users
	SET name = ?, email = ?, password = ?
	WHERE id = ?
	`

	if _, err := executor.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.ID,
	); err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) Delete(ctx context.Context, id string) error {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `DELETE FROM Users
	WHERE id = ?
	`

	if _, err := executor.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) LockUserByEmail(ctx context.Context, email string) (bool, error) {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `SELECT id
	FROM Users
	WHERE email = ?
	FOR UPDATE
	`

	row := executor.QueryRowContext(ctx, query, email)

	var id string
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("No user found with the provided email", log.Fstring("email", email))
			return false, nil
		}
		log.Error("Failed to scan row", log.Ferror(err))
		return false, err
	}
	return true, nil
}
