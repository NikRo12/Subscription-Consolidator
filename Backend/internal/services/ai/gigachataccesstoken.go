package ai

// Stores gigachat API http-response
type gigaChatAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}
