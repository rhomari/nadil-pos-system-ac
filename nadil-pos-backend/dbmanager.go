package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Server struct {
	DATABASE *sql.DB
}

var server = &Server{}

func (server *Server) connectToDB() {
	var err error
	if server.DATABASE != nil {

		return
	}
	dbHost := os.Getenv("NADIL_DB_HOST")
	dbPort := os.Getenv("NADIL_DB_PORT")
	dbUser := os.Getenv("NADIL_DB_USER")
	dbPassword := os.Getenv("NADIL_DB_PASSWORD")
	dbName := os.Getenv("NADIL_DB_NAME")
	dbSSLMode := os.Getenv("NADIL_DB_SSLMODE")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	server.DATABASE, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = server.DATABASE.Ping()
	if err != nil {
		log.Fatal(err)
	}

}
