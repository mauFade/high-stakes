package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mauFade/high-stakes/internal/core/service"
)

// Server represents the HTTP server
type Server struct {
	router      *chi.Mux
	server      *http.Server
	userService *service.UserService
	port        string
}

// NewServer creates a new HTTP server
func NewServer(us *service.UserService, p string) *Server {
	return &Server{
		router:      chi.NewRouter(),
		userService: us,
		port:        p,
	}
}

// SetupRoutes configures all HTTP routes
func (s *Server) SetupRoutes() {
	userHandler := NewUserHandler(s.userService)

	s.router.Post("/api/users", userHandler.CreateUser)
	s.router.Post("/api/auth/login", userHandler.Authenticate)
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}
	log.Println("HTTP server running at " + s.port)
	return s.server.ListenAndServe()
}
