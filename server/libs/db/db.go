package db

import (
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/surrealdb/surrealdb.go"
)

var db *surrealdb.DB

type mainDB struct {
	DbPtr *surrealdb.DB

	Host string
	Port int

	User string
	Pass string

	NS       string
	Database string
}

var dbConfig = &mainDB{
	Host: "localhost",
	Port: 8000,

	User: "root",
	Pass: "root",

	// NS:     "prod",
	NS:       "dev",
	Database: "ezMsg",
}

func InitDB() (*mainDB, string) {
	if dbConfig.DbPtr != nil {
		log.Fatal().Msg("DB already initialized")
	}

	var err error
	dbConfig.DbPtr, err = surrealdb.New("wss://" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port))
	if err != nil {
		log.Fatal().Msgf("Failed to connect to DB:\n%v", err)
	}
	log.Debug().Msg("Connected to DB successfully")

	token, err := dbConfig.DbPtr.SignIn(&surrealdb.Auth{
		Username: dbConfig.User,
		Password: dbConfig.Pass,

		Namespace: dbConfig.NS,
		Database:  dbConfig.Database,
	})
	if err != nil {
		log.Fatal().Msgf("Failed to sign in to DB:\n%v", err)
	}
	log.Debug().Msg("Signed in to DB successfully")

	return dbConfig, token
}

func GetDB() *mainDB {
	if dbConfig.DbPtr == nil {
		log.Fatal().Msg("DB not initialized")
	}

	return dbConfig
}
