{define name CreateUser}
CREATE TABLE IF NOT EXISTS user (
  id INTEGER NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  age INTEGER NOT NULL
);
{end define}

-- name: SelectById
{define name SelectById}
SELECT * FROM user WHERE id = ${id};
{end define}

{define name UnionSelect}
SELECT name FROM user
UNION ALL
SELECT age FROM user;
{end define}

{define name TxInsert, tx true}
INSERT INTO user(name, age) VALUES (${name}, ${age});
{end define}

{define name TxDelete, tx true}
DELETE FROM user WHERE id := ${id};
{end define}

{define name CountAll, mapper basicType}
SELECT COUNT (*) FROM user;
{end define}