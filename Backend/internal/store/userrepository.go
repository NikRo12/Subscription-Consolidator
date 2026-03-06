package store

import "github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"

type UserRepository interface {
	CreateUser(*models.User) error
	FindOrCreateUser(*models.User) error
	GetUserByEmail(string) (*models.User, error)
}
