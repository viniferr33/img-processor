package restful

import "github.com/viniferr33/img-processor/internal/user"

type UserRegistrationResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func NewUserRegistrationResponse(user *user.User) *UserRegistrationResponse {
	return &UserRegistrationResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
