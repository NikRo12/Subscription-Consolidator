package store

import "github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"

type UserRepository struct {
	store *Store
}

func (r *UserRepository) CreateUser(u *models.User) (*models.User, error) {
	if err := r.store.db.QueryRow(
		"INSERT INTO users (email, refresh_token) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.RefreshToken,
	).Scan(&u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
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
