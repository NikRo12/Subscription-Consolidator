package models

type Task struct {
	UserID       int    `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	MessageID    int    `json:"message_id"`
}
