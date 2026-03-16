package sqlstore_test

import (
	"testing"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestCreateSub(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseDriver, databaseURL)
	defer teardown("subs")

	s := sqlstore.NewSqlStore(db)
	sub := models.TestSubscription(t)
	assert.NoError(t, s.Sub().CreateSub(sub))
}

func TestCreateUserSub(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseDriver, databaseURL)
	defer teardown("users", "subs", "user_subs")

	s := sqlstore.NewSqlStore(db)
	user := models.TestUser(t)
	sub := models.TestSubscription(t)
	user_sub := models.TestUserSubscription(t)

	assert.NoError(t, s.User().CreateUser(user))
	assert.NoError(t, s.Sub().CreateSub(sub))
	assert.NoError(t, s.Sub().CreateUserSub(user_sub))
}

func TestGetAllSubsForUser(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseDriver, databaseURL)
	defer teardown("users", "subs", "user_subs")

	s := sqlstore.NewSqlStore(db)
	user := models.TestUser(t)
	sub := models.TestSubscription(t)
	user_sub := models.TestUserSubscription(t)

	assert.NoError(t, s.User().CreateUser(user))
	assert.NoError(t, s.Sub().CreateSub(sub))
	assert.NoError(t, s.Sub().CreateUserSub(user_sub))

	_, err := s.Sub().GetAllSubsForUser(1)
	assert.NoError(t, err)
}
