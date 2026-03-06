package sqlstore

import (
	"database/sql"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
)

type SqlUserRepository struct {
	store *SqlStore
}

func (r *SqlUserRepository) CreateUser(u *models.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (email, refresh_token) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.RefreshToken,
	).Scan(&u.ID)
}

func (r *SqlUserRepository) FindOrCreateUser(u *models.User) error {
	err := r.store.db.QueryRow(
		"SELECT id, email, refresh_token FROM users WHERE email = $1",
		u.Email,
	).Scan(&u.ID, &u.Email, &u.RefreshToken)

	if err == nil {
		return nil
	}

	if err != sql.ErrNoRows {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (email, refresh_token) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.RefreshToken,
	).Scan(&u.ID)
}

func (r *SqlUserRepository) GetUserByEmail(email string) (*models.User, error) {
	u := &models.User{}

	if err := r.store.db.QueryRow(
		"SELECT id, email, refresh_token FROM users WHERE email = $1",
		email,
	).Scan(&u.ID,
		&u.Email,
		&u.RefreshToken,
	); err != nil {
		return nil, err
	}

	return u, nil
}
