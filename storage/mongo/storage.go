package mongo

import (
	"blogging_platform/storage/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoStorage struct {
	userRepo repo.MongoUserStorageI
}



func New(db *mongo.Database) repo.MongoStorageI {
	return &mongoStorage{
		userRepo: NewMongoUser(db),
	}
}

func (s *mongoStorage) User() repo.MongoUserStorageI {
	return s.userRepo
}
