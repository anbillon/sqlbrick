{define name CreateUser}
CREATE TABLE IF NOT EXISTS user (
  id serial NOT NULL PRIMARY KEY,
  name text NOT NULL,
  age int NOT NULL
);
{end define}

-- name: SelectById
{define name SelectById}
SELECT * FROM user WHERE id = ${id};
{end define}

{define name TxInsert, tx true}
INSERT INTO user(name, age) VALUES (${name}, ${age});
{end define}

{define name TxDelete, tx true}
DELETE FROM user WHERE id := ${id};
{end define}