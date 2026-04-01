package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"studybuddy/backend/services/availability/domain"
	"time"
)

type GCalConnect interface {
	// Returns the Google OAuth consent URL. The frontend redirects
	// the user to this URL. The state param encodes user ID securely.
	GetAuthUrl(userID string) (string, error)

	// Called after Google redirects back to our callback endpoint
	// Validates the state, exchanges the code for tokens, and persists
	// the GCalConnection for the user
	HandleCallback(ctx context.Context, state, code string) error
}

type gcalConnect struct {
	gcal     GCalProvider
	gcalRepo GCalRepository
	stateKey []byte // HMAC key for signing state, will be kept in the env
}

func NewGCalConnect(gcal GCalProvider, gcalRepo GCalRepository, stateKey []byte) GCalConnect {
	return &gcalConnect{
		gcal:     gcal,
		gcalRepo: gcalRepo,
		stateKey: stateKey,
	}
}

func (gc *gcalConnect) GetAuthUrl(userID string) (string, error) {
	state, err := gc.buildState(userID)
	if err != nil {
		return "", fmt.Errorf("failed to build state: %w", err)
	}
	return gc.gcal.GetAuthURL(state), nil
}

func (gc *gcalConnect) HandleCallback(ctx context.Context, state, code string) error {
	userID, err := gc.verifyState(state)
	if err != nil {
		return fmt.Errorf("failed to verify state: %w", err)
	}

	conn, err := gc.gcal.ExchangeCode(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to exchange code: %w", err)
	}

	conn.UserID = userID

	if err := gc.gcalRepo.UpsertConnection(conn); err != nil {
		return fmt.Errorf("persist gcal connection: %w", err)
	}
	return nil
}

// helpers

// buildState encodes "<userID>:<expiry_unix>" and appends an HMAC signature
// Format: base64(<userID>:<expiry>).<base64(hmac)>
func (gc *gcalConnect) buildState(userID string) (string, error) {
	expiry := time.Now().Add(10 * time.Minute).Unix()
	payload := fmt.Sprintf("%s:%d", userID, expiry)
	encodedPayload := base64.RawURLEncoding.EncodeToString([]byte(payload))

	mac := hmac.New(sha256.New, gc.stateKey)
	mac.Write([]byte(encodedPayload))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return encodedPayload + "|" + signature, nil
}

// verifyState validates the HMAC signature and expiry, then returns the userID
func (gc *gcalConnect) verifyState(state string) (string, error) {
	parts := strings.SplitN(state, "|", 2)
	if len(parts) != 2 {
		return "", domain.ErrInvalidOAuthState
	}
	encodedPayload, signature := parts[0], parts[1]

	mac := hmac.New(sha256.New, gc.stateKey)
	mac.Write([]byte(encodedPayload))
	expectedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return "", domain.ErrInvalidOAuthState
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return "", domain.ErrInvalidOAuthState
	}
	payload := string(payloadBytes)

	lastColon := strings.LastIndex(payload, ":")
	if lastColon < 0 {
		return "", domain.ErrInvalidOAuthState
	}
	userID := payload[:lastColon]
	expiryStr := payload[lastColon+1:]

	var expiry int64
	if _, err := fmt.Sscanf(expiryStr, "%d", &expiry); err != nil {
		return "", domain.ErrInvalidOAuthState
	}
	if time.Now().Unix() > expiry {
		return "", domain.ErrOAuthStateExpired
	}

	return userID, nil
}
