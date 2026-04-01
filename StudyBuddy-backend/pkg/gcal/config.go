package gcal

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// Provider implements usecase.GCalProvider.
type Provider struct {
	oauthCfg *oauth2.Config
}

func New(cfg Config) *Provider {
	oauthCfg := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{calendar.CalendarReadonlyScope},
		Endpoint:     google.Endpoint,
	}
	return &Provider{oauthCfg: oauthCfg}
}

func (p *Provider) GetAuthURL(state string) string {
	return p.oauthCfg.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	)
}
