package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"blogging_platform/api"
	"blogging_platform/config"
	"blogging_platform/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	cfg := config.NewConfig(".")
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DB,
	)

	dbPool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	connection, err := dbPool.Acquire(ctx)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error while acquiring connection from the database pool!!")
	}
	defer connection.Release()

	err = connection.Ping(ctx)
	if err != nil {
		log.Fatal("Could not ping database")
	}
	fmt.Println("PostgreSQL connected successfully")

	// ------------------------------------------
	// MongoDB connection
	clientOpts := options.Client().ApplyURI(cfg.Mongo.URI)
	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Ping to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connected successfully...")

	// ------------------------------------------
	// Neo4j connection
	neo4jDriver, err := neo4j.NewDriverWithContext(
		cfg.Neo4j.URI,
		neo4j.BasicAuth(cfg.Neo4j.User, cfg.Neo4j.Password, ""))
	if err != nil {
		panic(err)
	}
	defer neo4jDriver.Close(ctx)

	err = neo4jDriver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Neo4j connected successfully...")

	strg := storage.New(dbPool, client.Database(cfg.Mongo.DB), neo4jDriver)

	engine := api.New(&api.Handler{
		Strg: strg,
		Cfg:  &cfg,
	})

	if err = engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
