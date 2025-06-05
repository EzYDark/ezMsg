package main

import "github.com/ezydark/ezMsg/server/libs/db"

func main() {
	db, token := db.InitDB()
	// Invalidate the token when the program exits
	defer func(token string) {
		if err := db.DbPtr.Invalidate(); err != nil {
			panic(err)
		}
	}(token)
}
