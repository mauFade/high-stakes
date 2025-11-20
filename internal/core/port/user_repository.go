package port

import "github.com/mauFade/high-stakes/internal/core/domain"

// UserRepository defines the interface for user data persistence
type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
	List(limit, offset int) ([]*domain.User, error)
}


