package main

import (
	"database/sql"
	"os"

	"github.com/Piyushhbhutoria/grpc-api/logger"
	"github.com/Piyushhbhutoria/grpc-api/server"
	"github.com/Piyushhbhutoria/grpc-api/store"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	m, err := migrate.New("file://migration", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	logger.Init()
	store.InitSQL()
	server.Init()
}
