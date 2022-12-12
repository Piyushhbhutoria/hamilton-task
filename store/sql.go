package store

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Piyushhbhutoria/grpc-api/logger"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitSQL() {
	dbURL := os.Getenv("DATABASE_URL")

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	db = conn

	err = db.Ping()
	if err != nil {
		logger.LogMessage("error", "error pinging to db: %v", err)
		logger.LogMessage("debug", "reconnecting")
		InitSQL()
	}
	logger.LogMessage("info", "postgres db connected")
}

func GetSQL() *sql.DB {
	return db
}

func CloseSQL() {
	db.Close()
}
