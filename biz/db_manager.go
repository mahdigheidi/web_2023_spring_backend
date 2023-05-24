package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func connectToDB() *sql.DB {
	psqlInfo := "postgres://postgres:postgres@postgres_db:5432/web_spring?sslmode=disable"
	print(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to database")
	return db
}
