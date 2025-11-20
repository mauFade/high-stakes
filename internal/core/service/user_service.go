package service

import (
	"time"

	"github.com/mauFade/high-stakes/internal/core/domain"
	"github.com/mauFade/high-stakes/internal/core/port"
	"github.com/mauFade/high-stakes/internal/core/util"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo port.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo port.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(name, email, phone, password string) (*domain.User, error) {
	// Validate input
	if name == "" {
		return nil, domain.ErrNameRequired
	}
	if email == "" {
		return nil, domain.ErrEmailRequired
	}
	if phone == "" {
		return nil, domain.ErrPhoneRequired
	}
	if password == "" {
		return nil, domain.ErrPasswordRequired
	}

	// Validate phone number
	if !util.ValidatePhone(phone) {
		return nil, domain.ErrInvalidPhone
	}

	// Check if user with email already exists
	existing, _ := s.userRepo.GetByEmail(email)
	if existing != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Generate KSUID
	id, err := util.GenerateKSUID()
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		ID:        id,
		Name:      name,
		Email:     email,
		Phone:     phone,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
