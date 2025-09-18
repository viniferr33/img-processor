package postgres

import (
	"database/sql"

	"github.com/viniferr33/img-processor/internal/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(id string) (*user.User, error) {
	query := "SELECT id, email, password_hash, updated_at, created_at FROM users WHERE id = $1"

	rows, err := r.db.Query(query, id)
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

func (r *UserRepository) GetByEmail(email string) (*user.User, error) {
	query := "SELECT id, email, password_hash, updated_at, created_at FROM users WHERE email = $1"

	rows, err := r.db.Query(query, email)
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

func (r *UserRepository) Create(u *user.User) error {
	query := "INSERT INTO users (id, email, name, password_hash, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.db.Exec(query, u.ID, u.Email, u.Name, u.PasswordHash, u.UpdatedAt, u.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Update(u *user.User) error {
	// Implementation goes here
	return nil
}

func (r *UserRepository) Delete(id string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
