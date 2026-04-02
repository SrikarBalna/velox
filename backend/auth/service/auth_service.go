package service

import (
	"errors"
	"fmt"
	"net/mail"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rishik92/velox/auth/model"
	"github.com/rishik92/velox/auth/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLen = 8
	tokenTTL       = 24 * time.Hour
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Signup validates input, hashes password, and persists the user.
func (s *AuthService) Signup(email, password string) (*model.User, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user, err := s.repo.CreateUser(email, string(hash))
	if errors.Is(err, repository.ErrEmailExists) {
		return nil, ErrEmailTaken
	}
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	return user, nil
}

// Login verifies credentials and returns a signed JWT token.
func (s *AuthService) Login(email, password string) (string, error) {
	if err := validateEmail(email); err != nil {
		return "", err
	}
	if err := validatePassword(password); err != nil {
		return "", err
	}

	user, err := s.repo.GetUserByEmail(email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", fmt.Errorf("fetching user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		return "", fmt.Errorf("generating token: %w", err)
	}
	return token, nil
}

// --- JWT helpers ---

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func generateJWT(userID string) (string, error) {
	secret := jwtSecret()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT parses and validates a JWT token string.
func ValidateJWT(tokenStr string) (*Claims, error) {
	secret := jwtSecret()
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func jwtSecret() string {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		// Defaulting to supersecurepassword123 for local testing parity
		return "supersecurepassword123"
	}
	return s
}

// --- Validation helpers ---

func validateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < minPasswordLen {
		return ErrPasswordTooShort
	}
	return nil
}

// --- Sentinel errors ---

var (
	ErrEmailTaken        = errors.New("email is already taken")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrPasswordTooShort  = fmt.Errorf("password must be at least %d characters", minPasswordLen)
	ErrInvalidToken      = errors.New("invalid or expired token")
)
