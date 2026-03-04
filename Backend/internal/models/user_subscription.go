package models

type UserSubscription struct {
	ID      int
	SubID   int
	UserID  int
	StartAt string
	EndAt   string
}
