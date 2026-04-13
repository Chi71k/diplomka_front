package gcal

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"

	"studybuddy/backend/services/availability/domain"
)

func (p *Provider) ExchangeCode(ctx context.Context, code string) (*domain.GCalConnection, error) {
	token, err := p.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("oauth2 code exchange: %w", err)
	}

	if !token.Valid() {
		return nil, fmt.Errorf("oauth2 exchange returned an invalid token")
	}

	if token.RefreshToken == "" {
		return nil, fmt.Errorf(
			"no refresh token in oauth2 response — ensure ApprovalForce is set and " +
				"the user granted offline access",
		)
	}

	return tokenToConnection(token), nil
}

func (p *Provider) RefreshToken(ctx context.Context, conn *domain.GCalConnection) (*domain.GCalConnection, error) {
	existing := &oauth2.Token{
		AccessToken:  conn.AccessToken,
		RefreshToken: conn.RefreshToken,
		Expiry:       conn.TokenExpiry,
		TokenType:    "Bearer",
	}

	tokenSource := p.oauthCfg.TokenSource(ctx, existing)

	refreshed, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("refresh oauth2 token: %w", err)
	}

	updated := tokenToConnection(refreshed)
	updated.UserID = conn.UserID
	updated.CalendarID = conn.CalendarID
	updated.SyncEnabled = conn.SyncEnabled
	updated.LastSyncedAt = conn.LastSyncedAt

	if updated.RefreshToken == "" {
		updated.RefreshToken = conn.RefreshToken
	}

	return updated, nil
}

func tokenToConnection(t *oauth2.Token) *domain.GCalConnection {
	return &domain.GCalConnection{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		TokenExpiry:  expiryOrFallback(t),
		CalendarID:   "primary",
		SyncEnabled:  true,
	}
}

func expiryOrFallback(t *oauth2.Token) time.Time {
	if !t.Expiry.IsZero() {
		return t.Expiry
	}
	return time.Now().Add(time.Hour)
}
