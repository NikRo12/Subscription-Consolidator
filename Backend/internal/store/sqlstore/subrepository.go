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
		) VALUES (
    		$1,  -- title
    		$2,  -- currency
    		$3,  -- category
    		$4,  -- icon_url
    		$5,  -- brand_color
    		$6   -- description
		)
		RETURNING id;`,
		s.Title,
		s.Currency,
		s.Category,
		s.IconURL,
		s.BrandColor,
		s.Description,
	).Scan(&s.ID)
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

func (r *SqlSubRepository) GetAllSubsForUser(userID int) ([]*models.Entry, error) {
	rows, err := r.store.db.Query(`
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
		WHERE us.user_id = $1`,
		userID,
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
