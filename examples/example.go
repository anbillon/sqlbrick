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

//go:generate sqlbrick -w ./sql/ -o ./models/
func main() {
	db, err := sqlx.Connect("postgres", "postgres://dev:developer@localhost:5432/dev?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	sqlBrick := models.NewSqlBrick(db)

	sqlBrick.Book.CreateBook()

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

	var book models.Book
	if err := sqlBrick.Book.SelectById(&book, 25); err != nil {
		log.Printf("wrror: %v", err)
	} else {
		log.Printf("select by id: %v", book)
	}

	var booksx []models.Book
	if err := sqlBrick.Book.SelectByUid(&booksx, 1324); err != nil {
		log.Printf("wrong: %v", err)
	} else {
		log.Printf("select by id: %v", booksx)
	}

	if _, err := sqlBrick.Book.DeleteById(26); err != nil {
		log.Printf("wrong: %v", err)
	}

	if _, err := sqlBrick.Book.UpdatePrice(&models.Book{
		Id:         30,
		Uid:        1324,
		Name:       "Sqlbrick",
		CreateTime: typex.NullTime{Time: time.Now(), Valid: true},
		Price:      200,
	}); err != nil {
		log.Printf("wrong: %v", err)
	}
}
