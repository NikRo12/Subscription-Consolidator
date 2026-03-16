package models

type User struct {
	ID           int    `json:"id"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

func (u *User) Sanitize() {
	u.RefreshToken = ""
}
