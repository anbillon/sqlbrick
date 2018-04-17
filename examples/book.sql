CREATE TABLE IF NOT EXISTS book (
  "id"  serial NOT NULL PRIMARY KEY,
  uid int NOT NULL,
  name text NOT NULL,
  content varchar(255),
  create_time TIMESTAMP,
  price int NOT NULL
);

-- name: InsertOne
INSERT INTO book (uid, name, content, create_time, price)
  VALUES (${uid}, ${name}, ${content}, ${create_time}, ${price});

-- name: UpdatePrice
UPDATE book SET price = ${price}, content = ${content} WHERE id = ${id};

-- name: ComplexUpdate
UPDATE book SET price=(SELECT price FROM book, user WHERE book.uid=user.id)
  WHERE book.price <= ${price} AND name = ${name};

-- name: SelectAll
SELECT * FROM book;

-- name: SelectById
SELECT * FROM book WHERE id = ${id} ORDER BY name ASC;

-- name: DeleteById
DELETE FROM book WHERE id = ${id};