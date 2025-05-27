package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Postgres PostgresData
	Mongo    MongoData
}

type PostgresData struct {
	Host     string
	Port     string
	Username string
	Password string
	DB       string
}

type MongoData struct {
	URI string
	DB  string
}

func NewConfig(path string) Config {
	godotenv.Load(path + "/.env")

	conf := viper.New()
	conf.AutomaticEnv()
	return Config{
		Postgres: PostgresData{
			Host:     conf.GetString("POSTGRES_HOST"),
			Port:     conf.GetString("POSTGRES_PORT"),
			Username: conf.GetString("POSTGRES_USER"),
			Password: conf.GetString("POSTGRES_PASSWORD"),
			DB:       conf.GetString("POSTGRES_DB"),
		},
		Mongo: MongoData{
			URI: conf.GetString("MONGO_URI"),
			DB:  conf.GetString("MONGO_DB_NAME"),
		},
	}
}
