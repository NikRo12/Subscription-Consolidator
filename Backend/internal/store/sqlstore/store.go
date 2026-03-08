package sqlstore

import (
	"database/sql"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/store"
	_ "github.com/lib/pq"
)

type SqlStore struct {
	db             *sql.DB
	userRepository *SqlUserRepository
	subRepository  *SqlSubRepository
}

func NewSqlStore(db *sql.DB) *SqlStore {
	return &SqlStore{
		db: db,
	}
}

func (s *SqlStore) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &SqlUserRepository{store: s}
	return s.userRepository
}

func (s *SqlStore) Sub() store.SubRepository {
	if s.subRepository != nil {
		return s.subRepository
	}

	s.subRepository = &SqlSubRepository{store: s}
	return s.subRepository
}
