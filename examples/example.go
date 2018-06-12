// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

package main

import (
	"log"
	"time"

	"anbillon.com/x/sqlbrick/examples/models"
	"anbillon.com/x/sqlbrick/typex"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
		Content:    typex.NewNullString("Golang program"),
		CreateTime: typex.NewNullTime(time.Now()),
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

	var count int
	if err := sqlBrick.Book.CountBooks(&count, 1324); err != nil {
		log.Printf("wrror: %v", err)
	} else {
		log.Printf("select count: %v", count)
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

	if _, err := sqlBrick.Book.UpdateSomeThing(&models.Book{
		Id:         30,
		Uid:        1324,
		Name:       "Sqlbrick",
		CreateTime: typex.NullTime{Time: time.Now(), Valid: true},
		Price:      200,
	}); err != nil {
		log.Printf("wrong: %v", err)
	}

	var txCount int
	if tx, err := sqlBrick.Begin(); err != nil {
		log.Printf("wrong: %v", err)
	} else {
		if _, err = tx.Book.TxInsert(&models.Book{
			Uid:        1234,
			Name:       "Tx",
			Content:    typex.NewNullString("Golang program"),
			CreateTime: typex.NewNullTime(time.Now()),
			Price:      30,
		}); err != nil {
			log.Printf("wrong: %v", err)
			return
		}
		if err := tx.Book.TxSelect(&txCount, 1234); err != nil {
			log.Printf("wrong: %v", err)
			return
		}
		var errTx = models.User{}
		if _, err := tx.Book.TxDeleteById(errTx); err != nil {
			log.Printf("wrong and rollback: %v", err)
			return
		}
		if err := tx.Commit(); err != nil {
			log.Printf("wrong: %v", err)
		} else {
			log.Printf("result: %v", txCount)
		}
	}
}
