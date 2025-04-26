package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Conn *sql.DB

func Init() {
	var err error
	err = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Printf("database url: %s", os.Getenv("DATABASE_URL"))

	Conn, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal("Failed to connect to DB ", err)
	}

	if err = Conn.Ping(); err != nil {
		log.Fatal("DB is unreachable: ", err)
	}
}
