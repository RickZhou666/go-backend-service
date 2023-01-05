package main

import (
	"database/sql"
	"log"

	"github.com/RickZhou666/go-backend-service/api"
	db "github.com/RickZhou666/go-backend-service/db/sqlc"
	"github.com/RickZhou666/go-backend-service/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot laod config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// 1. create a new store
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
