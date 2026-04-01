package crypto

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
)

// KeyFromHex decodes a 64-character hex string into a 32-byte AES-256 key.
// Example env value: ENCRYPTION_KEY=000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f
func KeyFromHex(h string) ([]byte, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return nil, fmt.Errorf("hex decode key: %w", err)
	}
	if len(b) != 32 {
		return nil, ErrInvalidKeySize
	}
	return b, nil
}

// KeyFromBase64 decodes a base64url-encoded string into a 32-byte AES-256 key.
// Example env value: ENCRYPTION_KEY=AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8=
func KeyFromBase64(s string) ([]byte, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		// Try standard encoding as fallback (some tools emit padding).
		var err2 error
		b, err2 = base64.StdEncoding.DecodeString(s)
		if err2 != nil {
			return nil, fmt.Errorf("base64 decode key: %w", err)
		}
	}
	if len(b) != 32 {
		return nil, ErrInvalidKeySize
	}
	return b, nil
}

func MustKeyFromEnv(envVar string) []byte {
	val := os.Getenv(envVar)
	if val == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", envVar))
	}

	if len(val) == 64 {
		key, err := KeyFromHex(val)
		if err == nil {
			return key
		}
	}

	key, err := KeyFromBase64(val)
	if err != nil {
		panic(fmt.Sprintf(
			"environment variable %q must be a 64-char hex or base64url-encoded 32-byte key: %v",
			envVar, err,
		))
	}
	return key
}
