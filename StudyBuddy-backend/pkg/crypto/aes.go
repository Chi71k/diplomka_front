package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var (
	ErrInvalidKeySize     = errors.New("encryption key must be exactly 32 bytes (AES-256)")
	ErrCiphertextTooShort = errors.New("ciphertext is too short to be valid")
	ErrDecryptFailed      = errors.New("decryption failed: ciphertext is corrupt or key is wrong")
)

// Encrypt encrypts plaintext using AES-256-GCM with a random nonce.
func Encrypt(key []byte, plaintext string) (string, error) {
	if len(key) != 32 {
		return "", ErrInvalidKeySize
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create aes cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize()) // 12 bytes
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}

	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.RawURLEncoding.EncodeToString(sealed), nil
}

// Decrypt decrypts a value produced by Encrypt.
func Decrypt(key []byte, encoded string) (string, error) {
	if len(key) != 32 {
		return "", ErrInvalidKeySize
	}

	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create aes cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create gcm: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize+gcm.Overhead() {
		return "", ErrCiphertextTooShort
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// Don't wrap the underlying error — it leaks implementation details.
		return "", ErrDecryptFailed
	}

	return string(plaintext), nil
}
