package sqlstore

import (
	"database/sql"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
)

type SqlUserRepository struct {
	store *SqlStore
}

func (r *SqlUserRepository) CreateUser(u *models.User) error {
	return r.store.db.QueryRow(
		"INSERT INTO users (refresh_token) VALUES ($1) RETURNING id",
		u.RefreshToken,
	).Scan(&u.ID)
}

func (r *SqlUserRepository) FindOrCreateUser(u *models.User) error {
	err := r.store.db.QueryRow(
		"SELECT id, refresh_token FROM users WHERE id = $1",
		u.ID,
	).Scan(&u.ID, &u.RefreshToken)

	if err == nil {
		return nil
	}

	if err != sql.ErrNoRows {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (refresh_token) VALUES ($1) RETURNING id",
		u.RefreshToken,
	).Scan(&u.ID)
}

func (r *SqlUserRepository) GetUserByID(id int) (*models.User, error) {
	u := &models.User{}

	if err := r.store.db.QueryRow(
		"SELECT id, refresh_token FROM users WHERE id = $1",
		id,
	).Scan(&u.ID,
		&u.RefreshToken,
	); err != nil {
		return nil, err
	}

	return u, nil
}
