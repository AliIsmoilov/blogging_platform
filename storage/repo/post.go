package repo

import (
	"context"
	"time"
)

type PostgresPostStorageI interface {
	Create(context.Context, *CreatePostReq) (*PostModelResp, error)
	// Update(context.Context, *UpdateUserReq) (*UserModelResp, error)
	// GetById(context.Context, int64) (*UserModelResp, error)
	// GetByEmail(context.Context, string) (*UserModelResp, error)
	// Delete(context.Context, int64) error
	GetAll(context.Context, *GetAllUserReq) (*GetAllPostsResp, error)

	// CreateUserNeo4j(ctx context.Context, req *UserModelRespMongo) (*UserModelRespMongo, error)
}

type MongoPostStorageI interface {
	Create(ctx context.Context, req *PostModelRespMongo) (*PostModelRespMongo, error)
	GetAll(ctx context.Context, req *GetAllUserReq) (*GetAllPostsResp, error)
}

type CreatePostReq struct {
	UserId  int
	Title   string
	Content string
}

type PostModelResp struct {
	Id        int64
	UserId    int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetAllPostsResp struct {
	Posts      []*PostModelResp
	UsersMongo []*PostModelRespMongo
	Count      int64
}

type PostModelRespMongo struct {
	Id        int64     `bson:"id"`
	UserId    int       `bson:"user_id,omitempty"`
	Title     string    `bson:"title"`
	Content   string    `bson:"content"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
}
