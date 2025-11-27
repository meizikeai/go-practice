// internal/pkg/crypto/manager.go
package crypto

import (
	"errors"
	"fmt"
	"go-practice/internal/config"
	"sync"
)

type Manager struct {
	current *AESGCM // The current key used for encryption (new key)
	old     *AESGCM // The old key, used only for decryption (for backward compatibility with old tokens)
	mu      sync.RWMutex
}

func NewManager(cfg *config.CryptoKeyInstance) (*Manager, error) {
	current, err := NewAESGCM(cfg.CurrentKey)
	if err != nil {
		return nil, fmt.Errorf("load current key failed: %w", err)
	}

	km := &Manager{current: current}

	if cfg.OldKey != "" {
		old, err := NewAESGCM(cfg.OldKey)
		if err != nil {
			return nil, fmt.Errorf("load old key failed: %w", err)
		}
		km.old = old
	}

	return km, nil
}

func (km *Manager) Encrypt(plaintext string) (string, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return km.current.Encrypt(plaintext)
}

func (km *Manager) Decrypt(ciphertext string) (string, error) {
	km.mu.RLock()

	if plain, err := km.current.Decrypt(ciphertext); err == nil {
		km.mu.RUnlock()
		return plain, nil
	}

	if km.old != nil {
		if plain, err := km.old.Decrypt(ciphertext); err == nil {
			km.mu.RUnlock()
			return plain, nil
		}
	}
	km.mu.RUnlock()

	return "", errors.New("invalid or expired token")
}

func (km *Manager) RotateKey(newKeyB64 string) error {
	newAES, err := NewAESGCM(newKeyB64)
	if err != nil {
		return err
	}

	km.mu.Lock()

	km.old = km.current
	km.current = newAES
	km.mu.Unlock()

	return nil
}
