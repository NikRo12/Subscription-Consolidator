package store

import "github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"

type SubRepository interface {
	CreateSub(*models.Subscription) error
	GetSubByName(string) (*models.Subscription, error)
	GetAllSubsForUser(int) ([]*models.Enrty, error)
}
