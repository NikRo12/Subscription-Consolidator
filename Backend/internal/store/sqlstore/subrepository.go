package sqlstore

import "github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"

type SqlSubRepository struct {
	store *SqlStore
}

func (r *SqlSubRepository) CreateSub(s *models.Subscription) error {
	return r.store.db.QueryRow(
		`INSERT INTO subs (
			title,
			currency,
			category,
			icon_url,
			brand_color,
			description
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (title) DO UPDATE
			SET
				currency    = EXCLUDED.currency,
				category    = EXCLUDED.category,
				icon_url    = EXCLUDED.icon_url,
				brand_color = EXCLUDED.brand_color,
				description = EXCLUDED.description
		RETURNING id;`,
		s.Title,
		s.Currency,
		s.Category,
		s.IconURL,
		s.BrandColor,
		s.Description,
	).Scan(&s.ID)
}

func (r *SqlSubRepository) DeleteAllUserSubs(userID int) error {
	_, err := r.store.db.Exec(
		`DELETE FROM user_subs WHERE user_id = $1`,
		userID,
	)
	return err
}

func (r *SqlSubRepository) CreateUserSub(us *models.UserSubscription) error {
	_, err := r.store.db.Exec(
		`INSERT INTO user_subs (
			user_id,
			sub_id, 
			price,
			period, 
			next_payment_date, 
			is_active
		) VALUES (
		 	$1, 
			$2, 
			$3, 
			$4,
			$5, 
			$6
		)`,
		us.UserID,
		us.SubID,
		us.Price,
		us.Period,
		us.NextPaymentDate,
		us.IsActive,
	)

	return err
}

func (r *SqlSubRepository) GetAllSubsForUser(userID int, category string) ([]*models.Entry, error) {
	query := `
		SELECT 
			s.id, 
			s.title, 
			us.price, 
			s.currency, 
			us.period, 
			s.category,
			us.next_payment_date,
			s.icon_url,
			s.brand_color,
			us.is_active,
			s.description
		FROM subs s 
		JOIN user_subs us ON s.id = us.sub_id 
		JOIN users u ON u.id = us.user_id
		WHERE us.user_id = $1
		`

	args := []any{userID}

	if category != "" {
		query += " AND s.category = $2"
		args = append(args, category)
	}

	rows, err := r.store.db.Query(query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	subs := make([]*models.Entry, 0)

	for rows.Next() {
		entry := &models.Entry{}
		if err := rows.Scan(
			&entry.SubID,
			&entry.Title,
			&entry.Price,
			&entry.Currency,
			&entry.Period,
			&entry.Category,
			&entry.NextPaymentDate,
			&entry.IconURL,
			&entry.BrandColor,
			&entry.IsActive,
			&entry.Description,
		); err != nil {
			return nil, err
		}
		subs = append(subs, entry)
	}
	return subs, nil
}
