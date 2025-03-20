package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type userModel struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type userRepository struct {
	db SQLExecutor
	tr transactionRepository
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

	var um userModel
	if err := row.Scan(
		&um.ID,
		&um.Name,
		&um.Email,
		&um.Password,
	); err != nil {
		return nil, err
	}
	return &entity.User{
		ID:       um.ID,
		Name:     um.Name,
		Email:    um.Email,
		Password: um.Password,
	}, nil
}

func (ur *userRepository) Create(ctx context.Context, user entity.User) error {
	if err := ur.tr.Transaction(ctx, func(ctx context.Context) error {
		exists, err := ur.lockUserByEmail(ctx, user.Email)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("user with this email already exists")
		}
		if err = ur.create(ctx, user); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) create(ctx context.Context, user entity.User) error {
	executor := ur.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query := `INSERT INTO Users (
	id, name, email, password
	)
	VALUES (?, ?, ?, ?)
	`

	um := userModel{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	if _, err := executor.ExecContext(
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

func (ur *userRepository) lockUserByEmail(ctx context.Context, email string) (bool, error) {
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
			return false, nil
		}
		return false, err
	}
	return true, nil
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

	um := userModel{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	if _, err := executor.ExecContext(
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
