package storage

import (
	mongoStorage "blogging_platform/storage/mongo"
	"blogging_platform/storage/postgres"
	"blogging_platform/storage/repo"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
)

type StorageI interface {
	MongoStorage() repo.MongoStorageI
	Postgres() repo.PostgresStorageI
	Neo4j() repo.Neo4jStorageI
	// User() repo.UserStorageI
}

type storage struct {
	mongo    repo.MongoStorageI
	postgres repo.PostgresStorageI
	neo4j    repo.Neo4jStorageI
	// userRepo repo.UserStorageI
}

// func New(db *pgxpool.Pool, mongo *mongo.Database, neo4j neo4j.DriverWithContext) StorageI {
// 	return &storage{
// userRepo: postgres.NewUser(db, mongo, neo4j),
// 	}
// }

func New(db *pgxpool.Pool, mongoDB *mongo.Database, neo4jDriver neo4j.DriverWithContext) StorageI {
	return &storage{
		mongo:    mongoStorage.New(mongoDB),
		postgres: postgres.New(db),
		// neo4j:    neo4jStore.New(neo4jDriver),
	}
}

// func (s *storage) User() repo.UserStorageI {
// 	return s.userRepo
// }

func (s *storage) MongoStorage() repo.MongoStorageI {
	return s.mongo
}

func (s *storage) Postgres() repo.PostgresStorageI {
	return s.postgres
}

func (s *storage) Neo4j() repo.Neo4jStorageI {
	return s.neo4j
}
