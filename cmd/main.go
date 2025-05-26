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
	fmt.Println("Connection Successfully")

	strg := storage.New(dbPool)

	engine := api.New(&api.Handler{
		Strg: strg,
		Cfg:  &cfg,
	})

	if err = engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
