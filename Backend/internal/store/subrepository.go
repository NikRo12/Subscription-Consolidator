package store

import "github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"

type SubRepository interface {
	CreateSub(*models.Subscription) error
	CreateUserSub(*models.UserSubscription) error
	GetAllSubsForUser(int, string) ([]*models.Entry, error)
}
