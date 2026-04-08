package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"strings"
	"time"

	"github.com/rishik92/velox/auth/model"
	"github.com/rishik92/velox/auth/repository"
)

type APIKeyService struct {
	repo *repository.APIKeyRepository
}

func NewAPIKeyService(repo *repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{repo: repo}
}

const (
	Prefix        = "velox_sk_"
	KeyLength     = 32
	ChecksumLength = 4
)

// GenerateKey creates a new API key for a user.
func (s *APIKeyService) GenerateKey(userID, name string, scopes []string, expiresAt *time.Time) (string, *model.APIKey, error) {
	// 1. Generate 32 bytes of secure random data
	randomBytes := make([]byte, KeyLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", nil, fmt.Errorf("generate random bytes: %w", err)
	}

	// 2. Encode to Base64Url
	encoded := base64.RawURLEncoding.EncodeToString(randomBytes)

	// 3. Compute CRC32 checksum
	checksum := crc32.ChecksumIEEE([]byte(encoded))
	checksumBytes := make([]byte, ChecksumLength)
	checksumBytes[0] = byte(checksum >> 24)
	checksumBytes[1] = byte(checksum >> 16)
	checksumBytes[2] = byte(checksum >> 8)
	checksumBytes[3] = byte(checksum)
	checksumHex := hex.EncodeToString(checksumBytes)

	// 4. Assemble the full key: prefix + encoded + checksum
	fullKey := fmt.Sprintf("%s%s%s", Prefix, encoded, checksumHex)

	// 5. Hash the full key using SHA-256 for storage
	hash := sha256.Sum256([]byte(fullKey))
	keyHash := hex.EncodeToString(hash[:])

	// 6. Create display hint (e.g., "velox_sk_****a1b2")
	displayHint := fmt.Sprintf("%s****%s", Prefix, checksumHex[len(checksumHex)-4:])

	apiKey := &model.APIKey{
		UserID:      userID,
		Name:        name,
		KeyHash:     keyHash,
		DisplayHint: displayHint,
		Scopes:      scopes,
		ExpiresAt:   expiresAt,
	}

	if err := s.repo.CreateKey(apiKey); err != nil {
		return "", nil, err
	}

	return fullKey, apiKey, nil
}

// ValidateKey checks if a plaintext API key is valid.
func (s *APIKeyService) ValidateKey(plaintextKey string) (*model.APIKey, error) {
	// 1. Basic format validation and prefix check
	if !strings.HasPrefix(plaintextKey, Prefix) {
		return nil, fmt.Errorf("invalid key prefix")
	}

	// 2. Extract components
	// Prefix is 9 chars
	// Checksum is 8 chars (4 bytes hex-encoded)
	if len(plaintextKey) < len(Prefix)+8 {
		return nil, fmt.Errorf("key too short")
	}

	checksumHex := plaintextKey[len(plaintextKey)-8:]
	encoded := plaintextKey[len(Prefix) : len(plaintextKey)-8]

	// 3. Early rejection: Verify checksum
	expectedChecksum := crc32.ChecksumIEEE([]byte(encoded))
	actualChecksumBytes, err := hex.DecodeString(checksumHex)
	if err != nil || len(actualChecksumBytes) != 4 {
		return nil, fmt.Errorf("invalid checksum format")
	}
	actualChecksum := uint32(actualChecksumBytes[0])<<24 |
		uint32(actualChecksumBytes[1])<<16 |
		uint32(actualChecksumBytes[2])<<8 |
		uint32(actualChecksumBytes[3])

	if expectedChecksum != actualChecksum {
		return nil, fmt.Errorf("checksum mismatch")
	}

	// 4. Hash and query database
	hash := sha256.Sum256([]byte(plaintextKey))
	keyHash := hex.EncodeToString(hash[:])

	apiKey, err := s.repo.GetKeyByHash(keyHash)
	if err != nil {
		return nil, err
	}

	// 5. Check expiration
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("key expired")
	}

	// 6. Update last_used_at
	_ = s.repo.UpdateLastUsed(apiKey.ID)

	return apiKey, nil
}
