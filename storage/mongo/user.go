package mongo

import (
	"context"
	"time"

	"blogging_platform/storage/repo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepo struct {
	db *mongo.Database
}

func NewMongoUser(db *mongo.Database) repo.MongoUserStorageI {
	return &mongoUserRepo{db: db}
}

func (u *mongoUserRepo) GetAll(ctx context.Context, req *repo.GetAllUserReq) (*repo.GetAllUserResp, error) {
	collection := u.db.Collection("users")

	// Filter: deleted_at == null OR deleted_at does not exist
	filter := bson.M{
		"$or": []bson.M{
			{"deleted_at": bson.M{"$exists": false}},
			{"deleted_at": nil},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	data := repo.GetAllUserResp{}

	for cursor.Next(ctx) {
		var user repo.UserModelRespMongo
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		data.UsersMongo = append(data.UsersMongo, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// Count documents matching the same filter
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	data.Count = count

	return &data, nil
}

func (u *mongoUserRepo) Create(ctx context.Context, req *repo.UserModelRespMongo) (*repo.UserModelRespMongo, error) {
	collection := u.db.Collection("users")

	now := time.Now()
	if req.Id == 0 {
		req.Id = now.UnixNano()
	}

	// Prepare the document to insert
	userDoc := bson.M{
		"id":           req.Id, // If you want to use a specific ID, otherwise remove this line and let MongoDB generate it
		"full_name":    req.FullName,
		"email":        req.Email,
		"password":     req.Password,
		"phone_number": req.PhoneNumber,
		"balance":      0.0,
		"created_at":   now,
		"updated_at":   nil, // if you use it
		"deleted_at":   nil, // if you use soft delete
	}

	_, err := collection.InsertOne(ctx, userDoc)
	if err != nil {
		return nil, err
	}

	// Prepare response
	user := &repo.UserModelRespMongo{
		Id:          req.Id,
		FullName:    req.FullName,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Balance:     req.Balance,
	}

	return user, nil
}
