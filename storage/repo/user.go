package repo

import (
	"context"
	"database/sql"
	"time"
)

type UserStorageI interface {
	Create(context.Context, *CreateUserReq) (*UserModelResp, error)
	Update(context.Context, *UpdateUserReq) (*UserModelResp, error)
	GetById(context.Context, int64) (*UserModelResp, error)
	GetByEmail(context.Context, string) (*UserModelResp, error)
	Delete(context.Context, int64) error
	GetAll(context.Context, *GetAllUserReq) (*GetAllUserResp, error)
}

type GetAllUserReq struct {
	Limit  string
	Offset string
	Query  string
}

type GetAllUserResp struct {
	Users []*UserModelResp
	Count int64
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
