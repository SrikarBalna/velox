package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/rishik92/velox/auth/model"
)

var ErrEmailExists = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user into the database.
func (r *UserRepository) CreateUser(name, email, passwordHash string) (*model.User, error) {
	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, password_hash, created_at
	`
	user := &model.User{}
	err := r.db.QueryRow(query, name, email, passwordHash).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		// pq error code 23505 = unique_violation
		if err.Error() != "" && isUniqueViolation(err) {
			return nil, ErrEmailExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

// GetUserByEmail fetches a user record by email.
func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`
	user := &model.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

// GetUserByID fetches a user record by ID.
func (r *UserRepository) GetUserByID(id string) (*model.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`
	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

// isUniqueViolation checks if a postgres error is a unique constraint violation.
func isUniqueViolation(err error) bool {
	return err != nil && len(err.Error()) > 0 &&
		(containsCode(err.Error(), "23505") || containsCode(err.Error(), "unique"))
}

func containsCode(msg, code string) bool {
	for i := 0; i <= len(msg)-len(code); i++ {
		if msg[i:i+len(code)] == code {
			return true
		}
	}
	return false
}
