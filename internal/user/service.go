package user

import (
	"context"

	"github.com/google/uuid"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, name, email, rawPassword string) (*User, error) {
	passwordHash, err := hashPassword(rawPassword)
	if err != nil {
		return nil, err
	}
	user := NewUser(uuid.New().String(), name, email, passwordHash)
	err = s.repo.Create(ctx, user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Authenticate(ctx context.Context, email, rawPassword string) (*User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !checkPasswordHash(rawPassword, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
