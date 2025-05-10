package repository

import (
	"context"
	"database/sql"

	"github.com/MaxFando/lms/user-service/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	FindByName(ctx context.Context, name string) (*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
	UpdateRefreshToken(ctx context.Context, userID int64, token string) error

	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	CreateTx(ctx context.Context, tx *sqlx.Tx, user *model.User) (int64, error)
	UpdateRefreshTokenTx(ctx context.Context, tx *sqlx.Tx, userID int64, token string) error
}

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *model.User) (int64, error) {
	var id int64
	err := r.db.QueryRowxContext(
		ctx,
		`INSERT INTO users (name, password, refresh_token, role)
         VALUES ($1, $2, $3, $4) RETURNING id`,
		user.Name, user.Password, user.RefreshToken, user.Role,
	).Scan(&id)
	return id, err
}

func (r *PostgresUserRepository) FindByName(ctx context.Context, name string) (*model.User, error) {
	var u model.User
	err := r.db.GetContext(ctx, &u, "SELECT * FROM users WHERE name=$1", name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var u model.User
	err := r.db.GetContext(ctx, &u, "SELECT * FROM users WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (r *PostgresUserRepository) List(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	err := r.db.SelectContext(ctx, &users, "SELECT * FROM users")
	return users, err
}

func (r *PostgresUserRepository) UpdateRefreshToken(ctx context.Context, userID int64, token string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET refresh_token=$1 WHERE id=$2", token, userID)
	return err
}

func (r *PostgresUserRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

func (r *PostgresUserRepository) CreateTx(ctx context.Context, tx *sqlx.Tx, user *model.User) (int64, error) {
	var id int64
	err := tx.QueryRowxContext(
		ctx,
		`INSERT INTO users (name, password, refresh_token, role)
         VALUES ($1, $2, $3, $4) RETURNING id`,
		user.Name, user.Password, user.RefreshToken, user.Role,
	).Scan(&id)
	return id, err
}

func (r *PostgresUserRepository) UpdateRefreshTokenTx(ctx context.Context, tx *sqlx.Tx, userID int64, token string) error {
	_, err := tx.ExecContext(ctx, "UPDATE users SET refresh_token=$1 WHERE id=$2", token, userID)
	return err
}
