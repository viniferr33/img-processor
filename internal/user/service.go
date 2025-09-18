package user

import "github.com/google/uuid"

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(name, email, rawPassword string) (*User, error) {
	passwordHash, err := hashPassword(rawPassword)
	if err != nil {
		return nil, err
	}
	user := NewUser(uuid.New().String(), name, email, passwordHash)
	err = s.repo.Create(user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Authenticate(email, rawPassword string) (*User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if !checkPasswordHash(rawPassword, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
