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