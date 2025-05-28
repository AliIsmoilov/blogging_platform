package neo4j

import (
	"blogging_platform/storage/repo"
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type neo4jUserRepo struct {
	driver neo4j.DriverWithContext
}

// Interface this should implement (update with your actual interface name)
func NewNeo4jUser(driver neo4j.DriverWithContext) repo.Neo4jUserStorageI {
	return &neo4jUserRepo{driver: driver}
}

func (u *neo4jUserRepo) CreateUserNeo4j(ctx context.Context, req *repo.UserModelRespMongo) (*repo.UserModelRespMongo, error) {
	session := u.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	// If no Id, generate one (same as Mongo)
	if req.Id == 0 {
		req.Id = time.Now().UnixNano()
	}

	// Cypher query to create user node
	query := `
	CREATE (user:User {
		id: $id,
		full_name: $full_name,
		email: $email,
		password: $password,
		phone_number: $phone_number,
		balance: 0.0,
		created_at: datetime($created_at),
		updated_at: null,
		deleted_at: null
	})
	RETURN user.id AS id, user.full_name AS full_name, user.email AS email, user.password AS password, user.phone_number AS phone_number, user.balance AS balance
	`

	// Parameters for query
	params := map[string]interface{}{
		"id":           req.Id,
		"full_name":    req.FullName,
		"email":        req.Email,
		"password":     req.Password,
		"phone_number": req.PhoneNumber,
		"created_at":   time.Now().Format(time.RFC3339),
	}

	// Run query in write transaction
	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		record, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if record.Next(ctx) {
			rec := record.Record()
			id, _ := rec.Get("id")
			fullName, _ := rec.Get("full_name")
			email, _ := rec.Get("email")
			password, _ := rec.Get("password")
			phoneNumber, _ := rec.Get("phone_number")
			balance, _ := rec.Get("balance")

			// Type assert each value to expected type
			return &repo.UserModelRespMongo{
				Id:          id.(int64),
				FullName:    fullName.(string),
				Email:       email.(string),
				Password:    password.(string),
				PhoneNumber: phoneNumber.(string),
				Balance:     balance.(float64),
			}, nil
		}

		return nil, record.Err()
	})

	if err != nil {
		return nil, err
	}

	return result.(*repo.UserModelRespMongo), nil
}
