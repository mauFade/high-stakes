package port

import "github.com/mauFade/high-stakes/internal/core/domain"

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(name, email, phone, password string) (*domain.User, error)
	GetUserByID(id string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(id, name, email, phone string) (*domain.User, error)
	DeleteUser(id string) error
	ListUsers(limit, offset int) ([]*domain.User, error)
}


