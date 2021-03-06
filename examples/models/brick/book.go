// Code generated by github.com/anbillon/sqlbrick. DO NOT EDIT IT.

// This file is generated from: book.sqb

package brick

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/anbillon/sqlbrick/examples/models/entity"
	"github.com/jmoiron/sqlx"
)

// Type definition for BookBrick. This brick will contains all database
// operation from given sqb file. Each sqb file will have only one brick.
type BookBrick struct {
	db *sqlx.DB
}

// Type definition for Book transaction. This aims at sql transaction.
type BookBrickTx struct {
	tx *sqlx.Tx
}

// newBookBrick will create a Book brick. This is used
// invoke the query function generated from sqb file.
func newBookBrick(db *sqlx.DB) *BookBrick {
	return &BookBrick{db: db}
}

// newBookTx will create a new transaction for Book.
func (b *BookBrick) newBookTx(tx *sqlx.Tx) *BookBrickTx {
	return &BookBrickTx{tx: tx}
}

// checkTx will check if tx is available.
func (b *BookBrickTx) checkTx() error {
	if b.tx == nil {
		return errors.New("the Begin func must be invoked first")
	}
	return nil
}

// CreateBook create table if not exsited
func (b *BookBrick) CreateBook() error {
	stmt, err := b.db.Prepare(`CREATE TABLE IF NOT EXISTS book (
  "id"  INTEGER NOT NULL PRIMARY KEY,
  uid INTEGER NOT NULL,
  name TEXT NOT NULL,
  content VARCHAR(255),
  create_time TIMESTAMP,
  price INTEGER NOT NULL
)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

// AddOne generated by sqlbrick, insert data into database.
func (b *BookBrick) AddOne(args *entity.Book) (sql.Result, error) {
	stmt, err := b.db.PrepareNamed(
		`INSERT INTO book (uid, name, content, create_time, price)
  VALUES (:uid, :name, :content, :create_time, :price)`)
	if err != nil {
		return nil, err
	}

	result, err := stmt.Exec(args)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateSomeThing an example to show update.
func (b *BookBrick) UpdateSomeThing(args *entity.Book) (int64, error) {
	conditionQuery := `UPDATE book SET `
	if args.Price > 0 {
		conditionQuery += `price = :price,`
	}
	if args.Content.String != "" {
		conditionQuery += `content = :content,`
	}
	conditionQuery += `name = :name,`
	if args.CreateTime.Time.Unix() != 0 {
		conditionQuery += `create_time = :create_time `
	}
	if strings.HasSuffix(conditionQuery, ",") {
		conditionQuery = strings.TrimSuffix(conditionQuery, ",")
	}
	conditionQuery += ` WHERE id = :id`

	stmt, err := b.db.PrepareNamed(conditionQuery)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(args)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// ComplexUpdate an example to show complex update. Second line comment.
func (b *BookBrick) ComplexUpdate(args *entity.Book) (int64, error) {
	stmt, err := b.db.PrepareNamed(
		`UPDATE book SET price=(SELECT price FROM book, user WHERE book.uid=user.id)
  WHERE book.price <= :price AND name = :name`)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(args)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// SelectAll generated by sqlbrick, select data from database.
func (b *BookBrick) SelectAll(dest *[]entity.CustomBook) error {
	return b.db.Select(dest, `SELECT id, uid, name, content FROM book`)
}

// CountBooks generated by sqlbrick, select data from database.
func (b *BookBrick) CountBooks(dest *int, uid interface{}) error {
	stmt, err := b.db.PrepareNamed(
		`SELECT COUNT(*) FROM book WHERE uid = :uid`)
	if err != nil {
		return err
	}

	// create map arguments for sqlx
	args := map[string]interface{}{
		"uid": uid,
	}

	return stmt.QueryRow(args).Scan(dest)
}

// SelectById an example to show SelectById.
func (b *BookBrick) SelectById(dest *entity.Book, id int, uid int) error {
	stmt, err := b.db.PrepareNamed(
		`SELECT * FROM book WHERE id = :id and uid = :uid`)
	if err != nil {
		return err
	}

	// create map arguments for sqlx
	args := map[string]interface{}{
		"id":  id,
		"uid": uid,
	}

	row := stmt.QueryRowx(args)
	if row.Err() != nil {
		return row.Err()
	}

	return row.StructScan(dest)
}

// SelectByUid generated by sqlbrick, select data from database.
func (b *BookBrick) SelectByUid(dest interface{}, uid interface{}) error {
	stmt, err := b.db.PrepareNamed(
		`SELECT * FROM book WHERE uid = :uid ORDER BY name ASC`)
	if err != nil {
		return err
	}

	// create map arguments for sqlx
	args := map[string]interface{}{
		"uid": uid,
	}

	rows, err := stmt.Queryx(args)
	if err != nil {
		return err
	}

	return sqlx.StructScan(rows, dest)
}

// DeleteById an example to show DeleteById.
func (b *BookBrick) DeleteById(id int) (int64, error) {
	stmt, err := b.db.PrepareNamed(`DELETE FROM book WHERE id = :id`)
	if err != nil {
		return 0, err
	}

	// create map arguments for sqlx
	args := map[string]interface{}{
		"id": id,
	}

	result, err := stmt.Exec(args)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// DeleteByIdAndUid generated by sqlbrick, delete data from database.
// Affected rows will return if there's no error.
func (b *BookBrick) DeleteByIdAndUid(id int, uid int) (int64, error) {
	stmt, err := b.db.PrepareNamed(
		`DELETE FROM book WHERE id = :id and uid = :uid`)
	if err != nil {
		return 0, err
	}

	// create map arguments for sqlx
	args := map[string]interface{}{
		"id":  id,
		"uid": uid,
	}

	result, err := stmt.Exec(args)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// TxInsert generated by sqlbrick, insert data into database.
func (b *BookBrickTx) TxInsert(args *entity.Book) (sql.Result, error) {
	if err := b.checkTx(); err != nil {
		return nil, err
	}

	stmt, err := b.tx.PrepareNamed(
		`INSERT INTO book (uid, name, content, create_time, price)
  VALUES (:uid, :name, :content, :create_time, :price)`)
	if err != nil {
		if rbe := b.tx.Rollback(); rbe != nil {
			return nil, rbe
		}

		return nil, err
	}

	result, err := stmt.Exec(args)
	if err != nil {
		if rbe := b.tx.Rollback(); rbe != nil {
			return nil, rbe
		}

		return nil, err
	}

	return result, nil
}

// TxSelect generated by sqlbrick, select data from database.
func (b *BookBrickTx) TxSelect(dest *int, uid interface{}) error {
	if err := b.checkTx(); err != nil {
		return err
	}

	stmt, err := b.tx.PrepareNamed(
		`SELECT COUNT(*) FROM book WHERE uid = :uid`)
	if err != nil {
		if rbe := b.tx.Rollback(); rbe != nil {
			return rbe
		}

		return err
	}

	// create map arguments for sqlx
	args := map[string]interface{}{
		"uid": uid,
	}

	return stmt.QueryRow(args).Scan(dest)
}

// TxDeleteById generated by sqlbrick, delete data from database.
// Affected rows will return if there's no error.
func (b *BookBrickTx) TxDeleteById(id interface{}) (int64, error) {
	if err := b.checkTx(); err != nil {
		return 0, err
	}

	stmt, err := b.tx.PrepareNamed(`DELETE FROM book WHERE id = :id`)
	if err != nil {
		if rbe := b.tx.Rollback(); rbe != nil {
			return 0, rbe
		}

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

// TxUpdate generated by sqlbrick, update data in database.
func (b *BookBrickTx) TxUpdate(args *entity.Book) (int64, error) {
	conditionQuery := `UPDATE book SET `
	if args.Price > 0 {
		conditionQuery += `price = :price,`
	}
	if args.Content.String != "" {
		conditionQuery += `content = :content,`
	}
	conditionQuery += ` name = :name WHERE id = :id`

	if err := b.checkTx(); err != nil {
		return 0, err
	}

	stmt, err := b.tx.PrepareNamed(conditionQuery)
	if err != nil {
		if rbe := b.tx.Rollback(); rbe != nil {
			return 0, rbe
		}

		return 0, err
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
