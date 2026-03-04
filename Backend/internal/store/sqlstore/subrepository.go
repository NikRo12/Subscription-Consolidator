package sqlstore

import "github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"

type SqlSubRepository struct {
	store *SqlStore
}

func (r *SqlSubRepository) CreateSub(s *models.Subscription) error {
	if err := s.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO subs (name, url, price) VALUES ($1, $2, $3) RETURNING id",
		s.Name,
		s.URL,
		s.Price,
	).Scan(&s.ID)
}

func (r *SqlSubRepository) GetSubByName(name string) (*models.Subscription, error) {
	s := &models.Subscription{}

	if err := r.store.db.QueryRow(
		"SELECT id, name, url, price FROM subs WHERE name = $1",
		name,
	).Scan(&s.Name,
		&s.ID,
		&s.URL,
		&s.Price,
	); err != nil {
		return nil, err
	}

	return s, nil
}

func (r *SqlSubRepository) CreateUserSub(us *models.UserSubscription) error {
	return r.store.db.QueryRow("INSERT INTO user_subs (user_id, sub_id, start_at, end_at) VALUES ($1, $2, $3, $4) RETURNING id",
		us.UserID,
		us.SubID,
		us.StartAt,
		us.EndAt,
	).Scan(&us.ID)
}

func (r *SqlSubRepository) GetAllSubsForUser(userID int) ([]*models.Enrty, error) {
	rows, err := r.store.db.Query(`
		SELECT u.id, u.email, s.name, us.start_at, us.end_at
		FROM subs s 
		JOIN user_subs us ON s.id = us.sub_id 
		JOIN users u ON u.id = us.user_id
		WHERE us.user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	subs := make([]*models.Enrty, 0)

	for rows.Next() {
		s := &models.Enrty{}
		if err := rows.Scan(&s.UserID, &s.Email, &s.SubName, &s.StartAt, &s.EndAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}
