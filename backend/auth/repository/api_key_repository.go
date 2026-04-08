package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/rishik92/velox/auth/model"
)

var ErrKeyNotFound = errors.New("api key not found")

type APIKeyRepository struct {
	db *sql.DB
}

func NewAPIKeyRepository(db *sql.DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

// CreateKey inserts a new API key into the database.
func (r *APIKeyRepository) CreateKey(key *model.APIKey) error {
	query := `
		INSERT INTO api_keys (user_id, name, key_hash, display_hint, scopes, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	scopesJSON, err := json.Marshal(key.Scopes)
	if err != nil {
		return fmt.Errorf("marshal scopes: %w", err)
	}

	err = r.db.QueryRow(
		query,
		key.UserID,
		key.Name,
		key.KeyHash,
		key.DisplayHint,
		scopesJSON,
		key.ExpiresAt,
	).Scan(&key.ID, &key.CreatedAt)

	if err != nil {
		return fmt.Errorf("create api key: %w", err)
	}

	return nil
}

// GetKeyByHash fetches an API key by its SHA-256 hash.
func (r *APIKeyRepository) GetKeyByHash(hash string) (*model.APIKey, error) {
	query := `
		SELECT id, user_id, name, key_hash, display_hint, scopes, expires_at, last_used_at, created_at
		FROM api_keys
		WHERE key_hash = $1
	`
	key := &model.APIKey{}
	var scopesJSON []byte

	err := r.db.QueryRow(query, hash).Scan(
		&key.ID,
		&key.UserID,
		&key.Name,
		&key.KeyHash,
		&key.DisplayHint,
		&scopesJSON,
		&key.ExpiresAt,
		&key.LastUsedAt,
		&key.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get key by hash: %w", err)
	}

	if err := json.Unmarshal(scopesJSON, &key.Scopes); err != nil {
		return nil, fmt.Errorf("unmarshal scopes: %w", err)
	}

	return key, nil
}

// UpdateLastUsed updates the last_used_at timestamp of an API key.
func (r *APIKeyRepository) UpdateLastUsed(id string) error {
	query := `UPDATE api_keys SET last_used_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}
