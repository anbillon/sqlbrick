// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

package main

import (
	"log"

	"anbillon.com/sqlbrick/examples/models"
	"anbillon.com/sqlbrick/typex"
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

func main() {
	db, err := sqlx.Connect("postgres", "postgres://dev:developer@localhost:5432/dev?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	sqlBrick := models.NewSqlBrick(db)

	if _, err = sqlBrick.Book.InsertOne(&models.Book{
		Uid:        1324,
		Name:       "Golang",
		Content:    sql.NullString{String: "Golang program", Valid: true},
		CreateTime: typex.NullTime{Time: time.Now(), Valid: true},
		Price:      20,
	}); err != nil {
		log.Printf("insert error: %v", err)
	}

	var books []models.Book
	if err := sqlBrick.Book.SelectAll(&books); err != nil {
		log.Printf("error: %v", err)
	} else {
		log.Printf("select all: %v", books)
	}
}
