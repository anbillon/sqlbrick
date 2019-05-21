sqlbrick
========

[![Build Status](https://travis-ci.org/Tourbillon/sqlbrick.svg?branch=develop)](https://travis-ci.org/Tourbillon/sqlbrick) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/Tourbillon/sqlbrick/master/LICENSE)

SQLBrick generates golang function from your SQL statements. It's not another orm library, but a tool to generate golang function and models from given SQL files. The generated source code is totally based on  [sqlx][1]. 

# Why this
As metioned above, this is not an orm library. If you are looking for some orm library in go, this is not suitable for you. If you like to write SQL statements, but you don't want to write SQL function again and again, then this tool will help you to reduce workload.

# Install
```shell
go get -u anbillon.com/sqlbrick/cmd/sqlbrick
```
Add the following to your dependency if you are working with `go modules`
```mod
require (
	github.com/anbillon/sqlbrick/typex develop
)
```

# Usage
To use sqlbrick, put your SQL statements in `.sql` file. Typically the first statement creates a table. Each SQL file can include only one `CREATE TABLE`. The statement will be a little different from standard SQL statement, it uses `{}` as some simple syntax and `${}` as  placeholder. Other syntax in standard SQL can be used as usual such as comment  `--`. Here's an example:

```sql
{define name CreateBook}
CREATE TABLE IF NOT EXISTS book (
  "id"  serial NOT NULL PRIMARY KEY,
  uid int NOT NULL,
  name text NOT NULL,
  content varchar(255),
  create_time TIMESTAMP,
  price int NOT NULL
);
{end define}

{define name InsertOne}
INSERT INTO book (uid, name, content, create_time, price)
  VALUES (${uid}, ${name}, ${content}, ${create_time}, ${price});
{end define}

-- An example to show update price.
{define name UpdatePrice}
UPDATE book SET
{if price > 0} price = ${price}, {end if}
{if content != ""} content = ${content}, {end if}
name = ${name} WHERE id = ${id};
{end define}

-- An example to show complex update.
-- Second line comment.
{define name ComplexUpdate}
UPDATE book SET price=(SELECT price FROM book, user WHERE book.uid=user.id)
  WHERE book.price <= ${price} AND name = ${name};
{end define}

{define name SelectAll}
SELECT * FROM book;
{end define}

-- An example to show SelectById.
{define name SelectById, mapper single}
SELECT * FROM book WHERE id = ${id} ORDER BY name ASC;
{end define}

{define name SelectByUid, mapper array}
SELECT * FROM book WHERE uid = ${uid} ORDER BY name ASC;
{end define}

-- An example to show DeleteById.
{define name DeleteById}
DELETE FROM book WHERE id = ${id};
{end define}

{define name TxInsert, tx true}
INSERT INTO book (uid, name, content, create_time, price)
  VALUES (${uid}, ${name}, ${content}, ${create_time}, ${price});
{end define}

{define name TxDeleteById, tx true}
DELETE FROM book WHERE id = ${id};
{end define}
```
Go to your work directory where contains `.sql` files, then run:
```shel
$ sqlbrick
```
Or you can specify the directory of sql files. For more usage, check
```shell
$ sqlbrick -help
```
You  can also use SQLBrick with `go generate`, just add the following somewhere in your source code and then run `go generate`:
```text
//go:generate sqlbrick -w your/sql/dir -o output/dir
```

From this sqlbrick will generate `sqlbrick.go`, `book.go`. You'll see all SQL function in `book.go` and `Book` model. At finally, you can use them in your code.
```go
import (
	"log"
	
	"github.com/jmoiron/sqlx"
	_"github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", "postgres://user:pass@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	sqlBrick := models.NewSqlBrick(db)

	var books []Book
	err := sqlBrick.Book.SelectById(&books, "someid")
	if err != nil {
		log.Printf("error: %v", err)
	}
	
	// your code here
}
```
The value in placeholder should keep the same with SQL field, or the same with your custom struct which tagged with `db`. For more detail, you can check the document of [sqlx][1]. 

# Syntax

### Definition
SQLBrick uses `{define ...}...{end define}` to define a SQL function:
```sql
{define name SelectById, mapper single}
....
{end define}
```
Definition has there parameters, `name`,  `mapper` and `tx`, they must be split by `,`. 
* The `name` is necessary to define the name of current SQL function.
* The `mapper` is optional, default is `array` which means the result will map to an array. If you want to map to only one result, then use `single`. If you want to map to basic type, then use `basicType`(ps: this should be based on your query result).
* The `tx` is optional, default is `false`. If set `true`, the SQL function will work in a transaction. Consuming the generated code like this:
```go
...
if tx, err := sqlBrick.Begin(); err != nil {
	// your code here
} else {
	tx.Book.SomeTxFunc()
	tx.Book.AnotherTxFunc()
	tx.Commit()
}
...
```
For detail usage,  check the `examples` in source code.

### Condition
SQLBrick uses `{if ...}...{end if}` as condition to make dynamic queries.

License
======
```text
MIT License

Copyright (C) 2018-present Anbillon Team

This source code is licensed under the MIT license found in the
LICENSE file in the root directory of this source tree.
```

[1]: https://github.com/jmoiron/sqlx
