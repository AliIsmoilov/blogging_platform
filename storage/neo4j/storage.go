package neo4j

import (
	"blogging_platform/storage/repo"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type neo4jStorage struct {
	userRepo repo.Neo4jUserStorageI
}

func New(driver neo4j.DriverWithContext) repo.Neo4jStorageI {
	return &neo4jStorage{
		userRepo: NewNeo4jUser(driver),
	}
}

func (s *neo4jStorage) User() repo.Neo4jUserStorageI {
	return s.userRepo
}
