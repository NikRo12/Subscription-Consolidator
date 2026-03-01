package store_test

import (
	"testing"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s, tearDown := store.TestStore(t, driver, databaseURL)
	defer tearDown("users")

	u, err := s.User().CreateUser(&models.User{
		Email:        "example@gmail.com",
		RefreshToken: "123456789",
	})

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	s, tearDown := store.TestStore(t, driver, databaseURL)
	defer tearDown("users")

	u, err := s.User().CreateUser(&models.User{
		Email:        "example@gmail.com",
		RefreshToken: "123456789",
	})

	assert.NoError(t, err)
	assert.NotNil(t, u)

	u2, err := s.User().GetUserByEmail("example@gmail.com")
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
