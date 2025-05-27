package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AliIsmoilov/blogging_platform/api"
	"github.com/AliIsmoilov/blogging_platform/config"
	"github.com/AliIsmoilov/blogging_platform/storage"
	"github.com/jackc/pgx/v5/pgxpool"
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

	strg := storage.New(dbPool, client.Database(cfg.Mongo.DB))

	engine := api.New(&api.Handler{
		Strg: strg,
		Cfg:  &cfg,
	})

	if err = engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
