package main

import (
	"github.com/rs/zerolog/log"

	"github.com/ezydark/ezMsg/server/libs/comm"
	"github.com/ezydark/ezMsg/server/libs/db"
)

func main() {
	db, token := db.InitDB()
	// Invalidate the token when the program exits
	defer func(token string) {
		if err := db.DbPtr.Invalidate(); err != nil {
			log.Fatal().Msgf("Failed to invalidate DB connection token:\n%v", err)
		}
	}(token)

	if err := comm.InitServer(); err != nil {
		log.Fatal().Msgf("Failed to initialize Server:\n%v", err)
	}
}
