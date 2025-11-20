package service

import (
	"errors"
	"time"

	"github.com/mauFade/high-stakes/internal/core/domain"
	"github.com/mauFade/high-stakes/internal/core/port"
	"github.com/mauFade/high-stakes/internal/core/util"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo port.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo port.UserRepository) port.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(name, email, phone, password string) (*domain.User, error) {
	// Validate input
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if phone == "" {
		return nil, errors.New("phone is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	// Check if user with email already exists
	existing, _ := s.userRepo.GetByEmail(email)
	if existing != nil {
		return nil, errors.New("user with this email already exists")
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

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *userService) GetUserByEmail(email string) (*domain.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates user information
func (s *userService) UpdateUser(id, name, email, phone string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	// Get existing user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if name != "" {
		user.Name = name
	}
	if email != "" {
		// Check if email is already taken by another user
		existing, _ := s.userRepo.GetByEmail(email)
		if existing != nil && existing.ID != id {
			return nil, errors.New("email is already taken")
		}
		user.Email = email
	}
	if phone != "" {
		user.Phone = phone
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(id)
}

// ListUsers retrieves a list of users with pagination
func (s *userService) ListUsers(limit, offset int) ([]*domain.User, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	return s.userRepo.List(limit, offset)
}


