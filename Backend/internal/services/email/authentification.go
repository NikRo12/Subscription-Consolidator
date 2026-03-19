package email

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2GoogleV2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type GoogleUserInfo struct {
	GoogleID     string
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	config *oauth2.Config
}

func NewAuthService(clientID, clientSecret string) *AuthService {
	return &AuthService{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  "postmessage",
			Scopes: []string{
				"https://www.googleapis.com/auth/gmail.readonly",
				oauth2GoogleV2.UserinfoProfileScope,
			},
		},
	}
}

func (s *AuthService) ExchangeAuthCode(ctx context.Context, serverAuthCode string, redirectURI string) (*GoogleUserInfo, error) {
	cfg := *s.config
	if redirectURI != "" {
		cfg.RedirectURL = redirectURI
	}

	token, err := cfg.Exchange(ctx, serverAuthCode)
	if err != nil {
		return nil, fmt.Errorf("exchange code failed: %w", err)
	}

	oauth2Service, err := oauth2GoogleV2.NewService(ctx, option.WithTokenSource(cfg.TokenSource(ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 service: %w", err)
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	return &GoogleUserInfo{
		GoogleID:     userInfo.Id,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
