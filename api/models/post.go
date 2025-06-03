package models

type CreatePostReq struct {
	UserId  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostModelResp struct {
	Id        int64   `json:"id"`
	UserId    int     `json:"user_id"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

type GetAllPostsResp struct {
	Posts []*PostModelResp `json:"posts"`
	Count int64            `json:"count"`
}
