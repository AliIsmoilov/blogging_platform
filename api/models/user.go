package models

type UserModelResp struct {
	Id          int64   `json:"id"`
	FullName    *string `json:"full_name"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Balance     float64 `json:"balance"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
}

type UpdateUserReq struct {
	Id          int64   `json:"id"`
	PhoneNumber *string `json:"phone_number"`
	FullName    *string `json:"full_name"`
}

type GetAllUsersResp struct {
	Users []*UserModelResp `json:"users"`
	Count int64            `json:"count"`
}

type CreateUserReq struct {
	FullName    *string `json:"full_name"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Password    string  `json:"password"`
}
