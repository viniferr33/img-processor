package restful

import (
	"context"
	"net/http"

	"github.com/viniferr33/img-processor/internal/constants"
	"github.com/viniferr33/img-processor/internal/jwt"
	"github.com/viniferr33/img-processor/internal/user"
	"github.com/viniferr33/img-processor/internal/utils"
)

type authHandler struct {
	userService user.UserService
	jwtService  jwt.JwtService
}

func NewAuthHandler(userService user.UserService, jwtService jwt.JwtService) *authHandler {
	return &authHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (h *authHandler) handleUserRegistration(w http.ResponseWriter, r *http.Request) {
	var reqBody UserRegistrationRequest
	if err := utils.ParseJsonBody(r, &reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reqBody.Name == "" || reqBody.Email == "" || reqBody.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(r.Context(), reqBody.Name, reqBody.Email, reqBody.Password)
	if err != nil {
		http.Error(w, "Error registering user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJsonResponse(w, http.StatusCreated, NewUserRegistrationResponse(user))
}

func (h *authHandler) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	var reqBody UserLoginRequest
	if err := utils.ParseJsonBody(r, &reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reqBody.Email == "" || reqBody.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Authenticate(r.Context(), reqBody.Email, reqBody.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := h.jwtService.SignToken(user.ID)
	if err != nil {
		http.Error(w, "Error generating token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, TokenResponse{Token: token})
}

func (h *authHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		token, ok := utils.SplitBearerToken(authHeader)
		if !ok {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := h.jwtService.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), constants.ContextKeyUserID, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
