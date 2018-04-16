// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

package main

import (
	"log"

	"anbillon.com/sqlbrick/examples/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	sqlBrick := models.NewSqlBrick(db)

	var books interface{}
	if err := sqlBrick.Book.SelectAll(books); err != nil {
		log.Printf("error: %v", err)
	}
}
