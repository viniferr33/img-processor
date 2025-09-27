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

func (s *UserService) Register(ctx context.Context, cmd UserRegistrationRequest) (*User, error) {
	passwordHash, err := hashPassword(cmd.Password)
	if err != nil {
		return nil, err
	}
	user := NewUser(uuid.New().String(), cmd.Name, cmd.Email, passwordHash)
	err = s.repo.Create(ctx, user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Authenticate(ctx context.Context, cmd UserLoginRequest) (*User, error) {
	user, err := s.repo.GetByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}

	if !checkPasswordHash(cmd.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
