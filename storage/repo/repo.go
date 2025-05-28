package repo

import (
	"context"
)

type MongoBookStorageI interface{}

type PostgresOrderStorageI interface{}

// interfaces for Neo4j features
type Neo4jUserStorageI interface {
	Create(ctx context.Context, req *UserModelRespMongo) (*UserModelRespMongo, error)
	GetAll(ctx context.Context, req *GetAllUserReq) ([]*UserModelRespMongo, error)
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
