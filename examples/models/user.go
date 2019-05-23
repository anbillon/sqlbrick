// Code generated by github.com/anbillon/sqlbrick. DO NOT EDIT IT.

// This file is generated from: user.sql

package models

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// Type definition for User which defined in sql file.
// This can be used as a model in database operation.
type User struct {
	Id   int32  `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

// Type definition for UserBrick. This brick will contains all database
// operation from given sql file. Each sql file will have only one brick.
type UserBrick struct {
	db *sqlx.DB
}

// Type definition for User transaction. This aims at sql transaction.
type UserBrickTx struct {
	tx *sqlx.Tx
}

// newUserBrick will create a User brick. This is used
// invoke the query function generated from sql file.
func newUserBrick(db *sqlx.DB) *UserBrick {
	return &UserBrick{db: db}
}

// newUserTx will create a new transaction for User.
func (b *UserBrick) newUserTx(tx *sqlx.Tx) *UserBrickTx {
	return &UserBrickTx{tx: tx}
}

// checkTx will check if tx is available.
func (b *UserBrickTx) checkTx() error {
	if b.tx == nil {
		return errors.New("the Begin func must be invoked first")
	}
	return nil
}

// CreateUser generated by sqlbrick, used to operate database table.
func (b *UserBrick) CreateUser() error {
	stmt, err := b.db.Prepare(`CREATE TABLE IF NOT EXISTS user (
  id serial NOT NULL PRIMARY KEY,
  name text NOT NULL,
  age int NOT NULL
)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

// SelectById name: SelectById
func (b *UserBrick) SelectById(dest interface{}, id interface{}) error {
	stmt, err := b.db.PrepareNamed(
		`SELECT * FROM user WHERE id = :id`)
	if err != nil {
		return err
	}

	// create map arguments for sqlx
	args := map[string]interface{}{
		"id": id,
	}

	rows, err := stmt.Queryx(args)
	if err != nil {
		return err
	}

	return sqlx.StructScan(rows, dest)
}

// UnionSelect generated by sqlbrick, select data from database.
func (b *UserBrick) UnionSelect(dest interface{}) error {
	return b.db.Select(dest, `SELECT name FROM user
UNION ALL
SELECT age FROM user`)
}

// TxInsert generated by sqlbrick, insert data into database.
func (b *UserBrickTx) TxInsert(args *User) (sql.Result, error) {
	if err := b.checkTx(); err != nil {
		return nil, err
	}

	stmt, err := b.tx.PrepareNamed(
		`INSERT INTO user(name, age) VALUES (:name, :age)`)
	if err != nil {
		return nil, err
	}

	if result, err := stmt.Exec(args); err != nil {
		if rbe := b.tx.Rollback(); rbe != nil {
			return nil, rbe
		}
		return nil, err
	} else {
		return result, nil
	}
}

// TxDelete generated by sqlbrick, delete data from database.
// Affected rows will return if there's no error.
func (b *UserBrickTx) TxDelete(id interface{}) (int64, error) {
	if err := b.checkTx(); err != nil {
		return 0, err
	}

	stmt, err := b.tx.PrepareNamed(`DELETE FROM user WHERE id := :id`)
	if err != nil {
		return 0, err
	}

	// create map arguments for sqlx
	args := map[string]interface{}{
		"id": id,
	}

	result, err := stmt.Exec(args)
	if err != nil {
		if rbe := b.tx.Rollback(); rbe != nil {
			return 0, rbe
		}
		return 0, err
	}

	return result.RowsAffected()
}

// CountAll generated by sqlbrick, select data from database.
func (b *UserBrick) CountAll(dest interface{}) error {
	stmt, err := b.db.Prepare(`SELECT COUNT (*) FROM user`)
	if err != nil {
		return err
	}
	return stmt.QueryRow().Scan(dest)
}
