package user

import "time"

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    int64
	UpdatedAt    int64
}

type UserRepository interface {
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id string) error
}

func NewUser(id, name, email, passwordHash string) *User {
	return &User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
}
