package storage

import (
	"github.com/AliIsmoilov/blogging_platform/storage/postgres"
	"github.com/AliIsmoilov/blogging_platform/storage/repo"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

type StorageI interface {
	User() repo.UserStorageI
}

type storage struct {
	userRepo repo.UserStorageI
}

func New(db *pgxpool.Pool, mongo *mongo.Database) StorageI {
	return &storage{
		userRepo: postgres.NewUser(db, mongo),
	}
}

func (s *storage) User() repo.UserStorageI {
	return s.userRepo
}
