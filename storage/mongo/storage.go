package mongo

import (
	"blogging_platform/storage/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoStorage struct {
	userRepo repo.MongoUserStorageI
	postRepo repo.MongoPostStorageI
}

func New(db *mongo.Database) repo.MongoStorageI {
	return &mongoStorage{
		userRepo: NewMongoUser(db),
		postRepo: NewMongoPost(db),
	}
}

func (s *mongoStorage) User() repo.MongoUserStorageI {
	return s.userRepo
}

// User returns the user repository
func (s *mongoStorage) Post() repo.MongoPostStorageI {
	return s.postRepo
}
