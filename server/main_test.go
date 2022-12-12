package server

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Piyushhbhutoria/grpc-api/logger"
	"github.com/Piyushhbhutoria/grpc-api/store"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var cc *grpc.ClientConn

func TestMain(m *testing.M) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	migrateDB, err := migrate.New("file://migration", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	err = migrateDB.Down()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	err = migrateDB.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	logger.Init()
	store.InitSQL()
	go Init()
	generateUser()

	cc, err = grpc.Dial("localhost:"+os.Getenv("PORT"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}

	os.Exit(m.Run())
}
