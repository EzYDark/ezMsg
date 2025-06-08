package db

import (
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/surrealdb/surrealdb.go"
)

type mainDB struct {
	DbPtr *surrealdb.DB

	Host   string
	Port   int
	Prefix string

	User string
	Pass string

	NS       string
	Database string
}

var DbConfig = &mainDB{
	DbPtr: nil,

	Host:   "localhost",
	Port:   8000,
	Prefix: "ws://", // TODO: Make secure WSS connection

	User: "ezy",
	Pass: "1234",

	// NS:     "prod",
	NS:       "dev",
	Database: "ezMsg",
}

func InitDB() (*mainDB, string) {
	if DbConfig.DbPtr != nil {
		log.Error().Msg("DB already initialized")
	}

	var err error
	DbConfig.DbPtr, err = surrealdb.New(DbConfig.Prefix + DbConfig.Host + ":" + strconv.Itoa(DbConfig.Port))
	if err != nil {
		log.Error().Msgf("Failed to connect to DB:\n%v", err)
	}
	log.Debug().Msg("Connected to DB successfully")

	token, err := DbConfig.DbPtr.SignIn(&surrealdb.Auth{
		Username: DbConfig.User,
		Password: DbConfig.Pass,

		Namespace: DbConfig.NS,
		Database:  DbConfig.Database,
	})
	if err != nil {
		log.Error().Msgf("Failed to sign in to DB:\n%v", err)
	}
	log.Debug().Msg("Signed in to DB successfully")

	return DbConfig, token
}

func GetDB() *mainDB {
	if DbConfig.DbPtr == nil {
		log.Fatal().Msg("DB not initialized")
	}

	return DbConfig
}
