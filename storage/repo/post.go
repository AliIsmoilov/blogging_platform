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
	Posts []*PostModelResp
	// UsersMongo []*UserModelRespMongo
	Count int64
}
