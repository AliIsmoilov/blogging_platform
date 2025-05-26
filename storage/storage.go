package storage

import (
	"github.com/AliIsmoilov/blogging_platform/storage/postgres"
	"github.com/AliIsmoilov/blogging_platform/storage/repo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StorageI interface {
	User() repo.UserStorageI
}

type storage struct {
	userRepo repo.UserStorageI
}

func New(db *pgxpool.Pool) StorageI {
	return &storage{
		userRepo: postgres.NewUser(db),
	}
}

func (s *storage) User() repo.UserStorageI {
	return s.userRepo
}
