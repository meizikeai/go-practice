// internal/pkg/jwt/jwt.go
package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"go-practice/internal/config"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Envoy / APISIX

// // generate jwt key
// openssl genrsa -out private_2025_01.pem 2048
// openssl rsa -in private_2025_01.pem -pubout -out public_2025_01.pem

var (
	ErrTokenExpired = errors.New("token has expired")
	ErrTokenInvalid = errors.New("invalid token")
)

type Manager struct {
	currentKey *rsa.PrivateKey
	currentID  string
	oldKey     *rsa.PrivateKey
	oldID      string
	mu         sync.RWMutex
}

func NewManager(cfg *config.JwtKeyInstance) (*Manager, error) {
	current, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.CurrentKey))
	if err != nil {
		return nil, err
	}

	m := &Manager{
		currentKey: current,
		currentID:  cfg.CurrentKeyID,
	}

	if cfg.OldKey != "" {
		old, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.OldKey))
		if err != nil {
			return nil, err
		}
		m.oldKey = old
		m.oldID = cfg.OldKeyID
	}

	return m, nil
}

func (m *Manager) GenerateToken(uid int64, duration time.Duration) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now()
	claims := jwt.MapClaims{
		"iss": "fuck",                     // issuer
		"sub": strconv.FormatInt(uid, 10), // subject
		"iat": now.Unix(),                 // issued at
		"exp": now.Add(duration).Unix(),   // expiration
		"nbf": now.Unix(),                 // not before
		"jti": generateJTI(),              // token id
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = m.currentID

	return token.SignedString(m.currentKey)
}

func (m *Manager) ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, jwt.ErrTokenMalformed
		}

		m.mu.RLock()
		defer m.mu.RUnlock()

		if kid == m.currentID {
			return &m.currentKey.PublicKey, nil
		}
		if m.oldKey != nil && kid == m.oldID {
			return &m.oldKey.PublicKey, nil
		}

		return nil, errors.New("unknown key id")
	}, jwt.WithValidMethods([]string{"RS256"}))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

func (m *Manager) RotateKey(newKeyPEM string, newKeyID string) error {
	newKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(newKeyPEM))
	if err != nil {
		return err
	}

	m.mu.Lock()
	m.oldKey = m.currentKey
	m.oldID = m.currentID
	m.currentKey = newKey
	m.currentID = newKeyID
	m.mu.Unlock()

	return nil
}

func generateJTI() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
