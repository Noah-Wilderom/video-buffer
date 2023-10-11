package main

import (
	"database/sql"
	"github.com/Noah-Wilderom/video-buffer/auth-service/models"
	"log"
	"os"
	"time"
)

const (
	gRPCPort = 5001
)

var (
	counts int64
)

type (
	Config struct {
		Models models.Models
	}
)

func main() {
	log.Println("Starting authentication service...")

	conn := createConnection()
	if conn == nil {
		log.Panic("Can't connect to postgres")
	}

	app := Config{
		Models: models.New(conn),
	}

}

func openConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createConnection() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		conn, err := openConnection(dsn)
		if err != nil {
			log.Println("Postgres not ready yet")
			counts++
		} else {
			log.Println("Connected to Postgres")
			return conn
		}

		if counts > 10 {
			log.Println(err)

			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
	}
}
