package google

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type GoogleUserInfo struct {
	GoogleID     string
	AccessToken  string
	RefreshToken string
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Error        string `json:"error"`
}

type userInfoResponse struct {
	Sub string `json:"sub"`
}

func ExchangeAuthCode(ctx context.Context, serverAuthCode string) (*GoogleUserInfo, error) {
	tokens, err := exchangeCodeForTokens(ctx, serverAuthCode)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	googleID, err := fetchGoogleID(ctx, tokens.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("fetch google id: %w", err)
	}

	return &GoogleUserInfo{
		GoogleID:     googleID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func exchangeCodeForTokens(ctx context.Context, code string) (*tokenResponse, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("GOOGLE_CLIENT_SECRET"))
	data.Set("redirect_uri", "postmessage")
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx,
		"POST",
		"https://oauth2.googleapis.com/token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &tokenResponse{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	if result.Error != "" {
		return nil, errors.New(result.Error)
	}

	return result, nil
}

func fetchGoogleID(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx,
		"GET",
		"https://www.googleapis.com/oauth2/v3/userinfo",
		nil,
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	userInfo := &userInfoResponse{}
	if err := json.Unmarshal(body, userInfo); err != nil {
		return "", err
	}

	if userInfo.Sub == "" {
		return "", errors.New("empty google id in response")
	}

	return userInfo.Sub, nil
}
