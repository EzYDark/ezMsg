package db

import (
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/surrealdb/surrealdb.go"
)

var db *surrealdb.DB

var dbConfig = struct {
	Host string
	Port int

	User string
	Pass string

	NS       string
	Database string
}{
	Host: "localhost",
	Port: 8000,

	User: "root",
	Pass: "root",

	// NS:     "prod",
	NS:       "dev",
	Database: "ezMsg",
}

func Init() *surrealdb.DB {
	if db != nil {
		log.Fatal().Msg("DB already initialized")
	}

	var err error
	db, err = surrealdb.New("wss://" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to DB")
	}
	log.Debug().Msg("Connected to DB successfully")

	_, err = db.SignIn(&surrealdb.Auth{
		Username: dbConfig.User,
		Password: dbConfig.Pass,

		Namespace: dbConfig.NS,
		Database:  dbConfig.Database,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to sign in to DB")
	}
	log.Debug().Msg("Signed in to DB successfully")

	return db
}

func GetDB() *surrealdb.DB {
	if db == nil {
		log.Fatal().Msg("DB not initialized")
	}

	return db
}
