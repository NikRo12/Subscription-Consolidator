package sqlstore

import (
	"database/sql"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/store"
	_ "github.com/lib/pq"
)

type SqlStore struct {
	db             *sql.DB
	UserRepository *SqlUserRepository
	SubRepository  *SqlSubRepository
}

func New(db *sql.DB) *SqlStore {
	return &SqlStore{
		db: db,
	}
}

func (s *SqlStore) User() store.UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}

	s.UserRepository = &SqlUserRepository{store: s}
	return s.UserRepository
}

func (s *SqlStore) Sub() store.SubRepository {
	if s.SubRepository != nil {
		return s.SubRepository
	}

	s.SubRepository = &SqlSubRepository{store: s}
	return s.SubRepository
}
