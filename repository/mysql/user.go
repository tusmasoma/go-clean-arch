package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/pkg/log"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type userModel struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type userRepository struct {
	db DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) Get(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT *
	FROM Users
	WHERE id = ?
	LIMIT 1`
	row := ur.db.QueryRowContext(ctx, query, id)
	var um userModel
	if err := row.Scan(
		&um.ID,
		&um.Name,
		&um.Email,
		&um.Password,
	); err != nil {
		return nil, err
	}
	user, err := entity.NewUser(
		um.ID,
		um.Name,
		um.Email,
		um.Password,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) Create(ctx context.Context, user entity.User) error {
	if err := ur.transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		exists, err := ur.lockUserByEmail(ctx, tx, user.Email)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("user with this email already exists")
		}
		if err = ur.create(ctx, tx, user); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) create(ctx context.Context, tx *sql.Tx, user entity.User) error {
	query := `INSERT INTO Users (
	id, name, email, password
	)
	VALUES (?, ?, ?, ?)
	`
	um := userModel{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.PasswordHash,
	}
	if _, err := tx.ExecContext(
		ctx,
		query,
		um.ID,
		um.Name,
		um.Email,
		um.Password,
	); err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) lockUserByEmail(ctx context.Context, tx *sql.Tx, email string) (bool, error) {
	query := `SELECT id
	FROM Users
	WHERE email = ?
	FOR UPDATE
	`

	row := tx.QueryRowContext(ctx, query, email)

	var id string
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (ur *userRepository) Update(ctx context.Context, user entity.User) error {
	query := `UPDATE Users
	SET name = ?, email = ?, password = ?
	WHERE id = ?
	`
	um := userModel{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.PasswordHash,
	}
	if _, err := ur.db.ExecContext(
		ctx,
		query,
		um.Name,
		um.Email,
		um.Password,
		um.ID,
	); err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM Users
	WHERE id = ?
	`
	if _, err := ur.db.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) transaction(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := ur.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
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
