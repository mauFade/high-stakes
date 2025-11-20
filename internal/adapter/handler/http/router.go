package http

import (
	"net/http"

	"github.com/mauFade/high-stakes/internal/core/port"
)

// Router sets up HTTP routes
type Router struct {
	userHandler *userHandler
}

// NewRouter creates a new router with all handlers
func NewRouter(userService port.UserService) *Router {
	return &Router{
		userHandler: NewUserHandler(userService),
	}
}

// SetupRoutes configures all HTTP routes
func (r *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// User routes
	mux.HandleFunc("/users", r.handleUsers)
	mux.HandleFunc("/users/", r.handleUserByID)

	// Health check
	mux.HandleFunc("/health", r.handleHealth)

	return mux
}

// handleUsers routes requests to /users
func (r *Router) handleUsers(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.userHandler.ListUsers(w, req)
	case http.MethodPost:
		r.userHandler.CreateUser(w, req)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleUserByID routes requests to /users/:id
func (r *Router) handleUserByID(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.userHandler.GetUser(w, req)
	case http.MethodPut:
		r.userHandler.UpdateUser(w, req)
	case http.MethodDelete:
		r.userHandler.DeleteUser(w, req)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleHealth handles health check requests
func (r *Router) handleHealth(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}


