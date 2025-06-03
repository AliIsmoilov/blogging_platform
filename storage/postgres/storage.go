package postgres

import (
	"blogging_platform/storage/repo"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Make sure this struct implements repo.PostgresStorageI
type postgresStorage struct {
	userRepo repo.PostgresUserStorageI
	postRepo repo.PostgresPostStorageI
}

// New creates a new PostgresStorageI instance
func New(db *pgxpool.Pool) repo.PostgresStorageI {
	return &postgresStorage{
		userRepo: NewPostgresUser(db),
		postRepo: NewPostgresPost(db),
	}
}

// User returns the user repository
func (s *postgresStorage) User() repo.PostgresUserStorageI {
	return s.userRepo
}

// User returns the user repository
func (s *postgresStorage) Post() repo.PostgresPostStorageI {
	return s.postRepo
}
