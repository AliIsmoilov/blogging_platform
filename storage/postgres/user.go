package postgres

import (
	"context"
	"fmt"

	"github.com/AliIsmoilov/blogging_platform/storage/repo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUser(db *pgxpool.Pool) repo.UserStorageI {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Create(ctx context.Context, req *repo.CreateUserReq) (*repo.UserModelResp, error) {
	query := `
        INSERT INTO users(
            full_name,
            email,
            password,
            phone_number,
            balance
        ) VALUES ($1, $2, $3, $4, 0) RETURNING id, full_name, email, password, phone_number, balance, created_at
    `
	var user repo.UserModelResp
	err := u.db.QueryRow(ctx, query, req.FullName, req.Email, req.Password, req.PhoneNumber).Scan(
		&user.Id,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.PhoneNumber,
		&user.Balance,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) Update(ctx context.Context, req *repo.UpdateUserReq) (*repo.UserModelResp, error) {
	query := `
        UPDATE users SET
        full_name = $1,
        phone_number = $2,
        updated_at = CURRENT_TIMESTAMP
        WHERE id = $3 
        RETURNING id, full_name, email, password, phone_number, balance, created_at, updated_at
  `
	var user repo.UserModelResp
	err := u.db.QueryRow(ctx, query, req.FullName, req.PhoneNumber, req.Id).Scan(
		&user.Id,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.PhoneNumber,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) GetById(ctx context.Context, id int64) (*repo.UserModelResp, error) {
	var resp repo.UserModelResp
	query := `
        SELECT
            id,
            full_name,
            email,
            password,
            phone_number,
            balance,
            created_at,
            updated_at
        FROM users WHERE id = $1
    `
	err := u.db.QueryRow(ctx, query, id).Scan(
		&resp.Id,
		&resp.FullName,
		&resp.Email,
		&resp.Password,
		&resp.PhoneNumber,
		&resp.Balance,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (u *userRepo) GetByEmail(ctx context.Context, email string) (*repo.UserModelResp, error) {
	var resp repo.UserModelResp
	query := `
        SELECT
            id,
            full_name,
            email,
            password,
            phone_number,
            balance,
            created_at,
            updated_at
        FROM users WHERE email = $1
    `
	err := u.db.QueryRow(ctx, query, email).Scan(
		&resp.Id,
		&resp.FullName,
		&resp.Email,
		&resp.Password,
		&resp.PhoneNumber,
		&resp.Balance,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (u *userRepo) Delete(ctx context.Context, userId int64) error {
	_, err := u.db.Exec(ctx, "DELETE FROM users WHERE id = $1", userId)
	return err
}

func (u *userRepo) GetAll(ctx context.Context, req *repo.GetAllUserReq) (*repo.GetAllUserResp, error) {
	query := `
        SELECT
            id,
            full_name,
            email,
            password,
            phone_number,
            balance,
            created_at,
            updated_at
        FROM users WHERE deleted_at IS NULL
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
                full_name ILIKE '%%%[1]v%%'
                OR phone_number ILIKE '%%%[1]v%%'
            )
        `, req.Query)
	}

	rows, err := u.db.Query(ctx, query+filter+order)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data := repo.GetAllUserResp{}
	for rows.Next() {
		var resp repo.UserModelResp
		err = rows.Scan(
			&resp.Id,
			&resp.FullName,
			&resp.Email,
			&resp.Password,
			&resp.PhoneNumber,
			&resp.Balance,
			&resp.CreatedAt,
			&resp.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		data.Users = append(data.Users, &resp)
	}

	err = u.db.QueryRow(ctx, "SELECT count(1) FROM users WHERE deleted_at IS NULL "+filter).Scan(&data.Count)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
