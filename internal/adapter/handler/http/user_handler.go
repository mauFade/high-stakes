package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mauFade/high-stakes/internal/core/dto"
	"github.com/mauFade/high-stakes/internal/core/service"
)

type userHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new HTTP user handler
func NewUserHandler(userService *service.UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}

// CreateUser handles POST /users
func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("failed to decode request body: %s", err.Error())})
		return
	}

	user, err := h.userService.CreateUser(req.Name, req.Email, req.Phone, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("failed to create user: %s", err.Error())})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *userHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthenticateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("failed to decode request body: %s", err.Error())})
		return
	}

	accessToken, refreshToken, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
		return
	}

	isProd := os.Getenv("ENVIRONMENT") == "production"
	sevenDays := 60 * 60 * 24 * 7

	http.SetCookie(w, &http.Cookie{
		Name:     "jid",
		Value:    refreshToken,
		Path:     "/api/auth/refresh",
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   sevenDays,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}
