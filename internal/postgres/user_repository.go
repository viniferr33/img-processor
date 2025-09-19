package postgres

import (
	"context"
	"database/sql"

	"github.com/viniferr33/img-processor/internal/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	query := "SELECT id, email, password_hash, updated_at, created_at FROM users WHERE id = $1"

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.UpdatedAt, &u.CreatedAt); err != nil {
			return nil, err
		}
		return &u, nil
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := "SELECT id, email, password_hash, updated_at, created_at FROM users WHERE email = $1"

	rows, err := r.db.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.UpdatedAt, &u.CreatedAt); err != nil {
			return nil, err
		}
		return &u, nil
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	query := "INSERT INTO users (id, email, name, password_hash, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Email, u.Name, u.PasswordHash, u.UpdatedAt, u.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	query := "UPDATE users SET email = $1, name = $2, password_hash = $3, updated_at = $4 WHERE id = $5"
	_, err := r.db.ExecContext(ctx, query, u.Email, u.Name, u.PasswordHash, u.UpdatedAt, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
