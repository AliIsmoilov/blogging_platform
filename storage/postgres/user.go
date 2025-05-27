package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/AliIsmoilov/blogging_platform/storage/repo"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	db    *pgxpool.Pool   // PostgreSQL DB
	mongo *mongo.Database // MongoDB DB
}

func NewUser(db *pgxpool.Pool, mongo *mongo.Database) repo.UserStorageI {
	return &userRepo{
		db:    db,
		mongo: mongo,
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

func (u *userRepo) GetAllMongo(ctx context.Context, req *repo.GetAllUserReq) (*repo.GetAllUserResp, error) {
	collection := u.mongo.Collection("users")

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

func (u *userRepo) CreateMongo(ctx context.Context, req *repo.UserModelRespMongo) (*repo.UserModelRespMongo, error) {
	collection := u.mongo.Collection("users")

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
