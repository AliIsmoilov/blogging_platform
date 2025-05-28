package postgres

import (
	"blogging_platform/storage/repo"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Make sure this struct implements repo.PostgresStorageI
type postgresStorage struct {
	userRepo repo.PostgresUserStorageI
}

// New creates a new PostgresStorageI instance
func New(db *pgxpool.Pool) repo.PostgresStorageI {
	return &postgresStorage{
		userRepo: NewPostgresUser(db), // Defined in user.go
	}
}

// User returns the user repository
func (s *postgresStorage) User() repo.PostgresUserStorageI {
	return s.userRepo
}
