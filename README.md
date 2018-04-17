sqlbrick
======
SQLBrick generates golang function from your SQL statements. It's not another orm library, but a tool to generate golang function and models from given SQL files. The generated source code is totally based on  [sqlx][1]. 

Why this
======
As metioned above, this is not an orm library. If you are looking for some orm library in go, this is not suitable for you. If you like to write SQL statements, but you don't want to write SQL function again and again, then this tool will help you to reduce workload.

# Install
```shell
go get -u anbillon.com/sqlbrick/cmd/sqlbrick
```
Add the following to your dependency if you are working with `dep`
```text
anbillon.com/sqlbrick/typex
```

# Usage
To use sqlbrick, put your SQL statements in `.sql` file. Typically the first statement creates a table. The statement will be a little different from standrad SQL statement, it uses `${}` as  placeholder. Here's an example:
```sql
CREATE TABLE IF NOT EXISTS book (
  id  int NOT NULL PRIMARY KEY,
  uid int NOT NULL,
  name text NOT NULL,
  content varchar(255),
  create_time TIMESTAMP,
  price int NOT NULL
);

-- name: InsertOne
INSERT INTO book (name, content, price) VALUES (${name}, ${content}, ${price});

-- name: UpdatePrice
UPDATE book SET price = ${price}, content = ${content} WHERE id = ${id};

-- name: SelectById
SELECT * FROM book WHERE id = ${id};

-- name: DeleteById
DELETE FROM book WHERE id = ${id};
```
Go to your work directory where contains `.sql` files, then run:
```shel
sqlbrick
```
Or you can specify the directory of sql files. For more usage, check
```shell
sqlbrick -help
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
The value in wildcards should keep the same with SQL field, or the same with your custom struct which tagged with `db`. For more detail, you can check the document of [sqlx][1]. 

License
======
```text
MIT License

Copyright (C) 2018-present Anbillon Team

This source code is licensed under the MIT license found in the
LICENSE file in the root directory of this source tree.
```

[1]: https://github.com/jmoiron/sqlx
