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
		"INSERT INTO users (google_id, refresh_token, access_token) VALUES ($1, $2, $3) RETURNING id",
		u.GoogleID, u.RefreshToken, u.AccessToken,
	).Scan(&u.ID)
}

func (r *SqlUserRepository) FindOrCreateUser(u *models.User) error {
	err := r.store.db.QueryRow(
		"SELECT id FROM users WHERE google_id = $1",
		u.GoogleID,
	).Scan(&u.ID)

	if err == nil {
		_, err = r.store.db.Exec(
			"UPDATE users SET refresh_token = $1, access_token = $2 WHERE id = $3",
			u.RefreshToken,
			u.AccessToken,
			u.ID,
		)
		return err
	}

	if err != sql.ErrNoRows {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (google_id, refresh_token, access_token) VALUES ($1, $2, $3) RETURNING id",
		u.GoogleID,
		u.RefreshToken,
		u.AccessToken,
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
