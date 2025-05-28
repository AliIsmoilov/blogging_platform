package repo

import (
	"context"
	"database/sql"
	"time"
)

type PostgresUserStorageI interface {
	Create(context.Context, *CreateUserReq) (*UserModelResp, error)
	Update(context.Context, *UpdateUserReq) (*UserModelResp, error)
	GetById(context.Context, int64) (*UserModelResp, error)
	GetByEmail(context.Context, string) (*UserModelResp, error)
	Delete(context.Context, int64) error
	GetAll(context.Context, *GetAllUserReq) (*GetAllUserResp, error)

	// CreateUserNeo4j(ctx context.Context, req *UserModelRespMongo) (*UserModelRespMongo, error)
}

type MongoUserStorageI interface {
	Create(ctx context.Context, req *UserModelRespMongo) (*UserModelRespMongo, error)
	GetAll(ctx context.Context, req *GetAllUserReq) (*GetAllUserResp, error)
}

type GetAllUserReq struct {
	Limit  string
	Offset string
	Query  string
}

type GetAllUserResp struct {
	Users      []*UserModelResp
	UsersMongo []*UserModelRespMongo
	Count      int64
}

type CreateUserReq struct {
	FullName    *string
	Email       string
	Password    string
	PhoneNumber *string
}

type UpdateUserReq struct {
	Id          int64
	FullName    *string
	PhoneNumber *string
}

type UserModelResp struct {
	Id          int64
	FullName    sql.NullString
	Email       string
	Password    string
	PhoneNumber sql.NullString
	Balance     float64
	CreatedAt   time.Time
	UpdatedAt   sql.NullTime
}

type UserModelRespMongo struct {
	Id          int64     `bson:"id"`
	FullName    string    `bson:"full_name,omitempty"`
	Email       string    `bson:"email"`
	Password    string    `bson:"password"`
	PhoneNumber string    `bson:"phone_number,omitempty"`
	Balance     float64   `bson:"balance"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at,omitempty"`
}
