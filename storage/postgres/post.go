package postgres

import (
	"blogging_platform/storage/repo"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresPostRepo struct {
	db *pgxpool.Pool
}

func NewPostgresPost(db *pgxpool.Pool) repo.PostgresPostStorageI {
	return &postgresPostRepo{db: db}
}

func (u *postgresPostRepo) Create(ctx context.Context, req *repo.CreatePostReq) (*repo.PostModelResp, error) {
	query := `
	    INSERT INTO posts (
	        user_id,
	        title,
	        content
	    ) VALUES ($1, $2, $3) RETURNING id, user_id, title, content, created_at
	`
	var post repo.PostModelResp
	err := u.db.QueryRow(ctx, query, req.UserId, req.Title, req.Content).Scan(
		&post.Id,
		&post.UserId,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (u *postgresPostRepo) GetAll(ctx context.Context, req *repo.GetAllUserReq) (*repo.GetAllPostsResp, error) {
	query := `
        SELECT
            id,
            user_id,
            title,
            content,
            created_at,
            updated_at
        FROM posts WHERE deleted_at IS NULL
    `
	var filter string
	order := " ORDER BY created_at DESC "

	if req.Limit != "" {
		order += fmt.Sprintf(" LIMIT %v ", req.Limit)
	}
	if req.Offset != "" {
		order += fmt.Sprintf(" OFFSET %v ", req.Offset)
	}
	if req.Query != "" {
		filter += fmt.Sprintf(`
            AND (
                title ILIKE '%%%[1]v%%'
            )
        `, req.Query)
	}

	rows, err := u.db.Query(ctx, query+filter+order)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data := repo.GetAllPostsResp{}
	for rows.Next() {
		var resp repo.PostModelResp
		err = rows.Scan(
			&resp.Id,
			&resp.UserId,
			&resp.Title,
			&resp.Content,
			&resp.CreatedAt,
			&resp.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		data.Posts = append(data.Posts, &resp)
	}

	err = u.db.QueryRow(ctx, "SELECT count(1) FROM posts WHERE deleted_at IS NULL "+filter).Scan(&data.Count)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
