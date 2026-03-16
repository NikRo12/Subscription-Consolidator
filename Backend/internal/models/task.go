package models

type Task struct {
	TaskID       string `json:task-id`
	UserID       string `json:user-id`
	AccessToken  string `json:access-token`
	RefreshToken string `json:refresh-token`
	MessageID    string `json:message-id`
}
