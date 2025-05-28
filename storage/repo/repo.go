package repo

import (
	"context"
)

type MongoBookStorageI interface{}

type PostgresOrderStorageI interface{}

// interfaces for Neo4j features
type Neo4jUserStorageI interface {
	CreateUserNeo4j(ctx context.Context, req *UserModelRespMongo) (*UserModelRespMongo, error)
}

type MongoStorageI interface {
	User() MongoUserStorageI
}

type Neo4jStorageI interface {
	User() Neo4jUserStorageI
}

// Top-Level Interfaces Per DB
type PostgresStorageI interface {
	User() PostgresUserStorageI
}
