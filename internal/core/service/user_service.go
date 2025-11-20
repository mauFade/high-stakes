package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (s *UserService) Authenticate(email, password string) (string, string, error) {
	// Validate input
	if email == "" {
		return "", "", domain.ErrEmailRequired
	}
	if password == "" {
		return "", "", domain.ErrPasswordRequired
	}

	// Check if user with email exists
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", "", err
	}

	// Check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", domain.ErrInvalidCredentials
	}

	jwtKey := os.Getenv("JWT_SECRET_KEY")
	now := time.Now()

	atClaims := jwt.MapClaims{
		"sub": user.ID,
		"iat": now.Unix(),
		"exp": now.Add(15 * time.Minute).Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(jwtKey))
	if err != nil {
		return "", "", err
	}

	rtClaims := jwt.MapClaims{
		"sub":  user.ID,
		"exp":  now.Add(7 * 24 * time.Hour).Unix(), // 7 days
		"type": "refresh",                          // Claim to differentiate
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	refreshKey := os.Getenv("JWT_REFRESH_SECRET_KEY")
	refreshToken, err := rt.SignedString([]byte(refreshKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}
