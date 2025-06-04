package mongo

import (
	"blogging_platform/storage/repo"
	"context"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoPostRepo struct {
	db *mongo.Database
}

func NewMongoPost(db *mongo.Database) repo.MongoPostStorageI {
	return &mongoPostRepo{db: db}
}

func (u *mongoPostRepo) Create(ctx context.Context, req *repo.PostModelRespMongo) (*repo.PostModelRespMongo, error) {
	collection := u.db.Collection("posts")

	now := time.Now()
	id := now.UnixNano()

	postDoc := bson.M{
		"id":         id,
		"user_id":    req.UserId,
		"title":      req.Title,
		"content":    req.Content,
		"created_at": now,
		"updated_at": nil,
		"deleted_at": nil,
	}

	_, err := collection.InsertOne(ctx, postDoc)
	if err != nil {
		return nil, err
	}

	return &repo.PostModelRespMongo{
		Id:        id,
		UserId:    req.UserId,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: now,
	}, nil
}

func (u *mongoPostRepo) GetAll(ctx context.Context, req *repo.GetAllUserReq) (*repo.GetAllPostsResp, error) {
	collection := u.db.Collection("posts")

	// Build base filter: exclude deleted posts
	filter := bson.M{
		"$or": []bson.M{
			{"deleted_at": bson.M{"$exists": false}},
			{"deleted_at": nil},
		},
	}

	// Add search on title if provided
	if req.Query != "" {
		filter["title"] = bson.M{"$regex": req.Query, "$options": "i"}
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", -1}})

	// Pagination
	if req.Limit != "" {
		limit, err := strconv.ParseInt(req.Limit, 10, 64)
		if err == nil {
			findOptions.SetLimit(limit)
		}
	}
	if req.Offset != "" {
		offset, err := strconv.ParseInt(req.Offset, 10, 64)
		if err == nil {
			findOptions.SetSkip(offset)
		}
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*repo.PostModelResp
	for cursor.Next(ctx) {
		var post repo.PostModelResp
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &repo.GetAllPostsResp{
		Posts: posts,
		Count: count,
	}, nil
}
