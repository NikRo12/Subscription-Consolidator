package models

type Task struct {
	TaskID       string `json:"task_id"`
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	MessageID    string `json:"message_id"`
}
