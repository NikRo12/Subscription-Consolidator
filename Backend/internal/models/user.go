package models

type User struct {
	ID           int    `json:"id"`
	GoogleID     string `json:"google_id"`
	RefreshToken string `json:"refreshToken,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
}

func (u *User) Sanitize() {
	u.RefreshToken = ""
	u.AccessToken = ""
}
