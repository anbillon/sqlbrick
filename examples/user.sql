CREATE TABLE IF NOT EXISTS user (
  id int NOT NULL PRIMARY KEY,
  name text NOT NULL,
  age int NOT NULL
);

-- name: SelectById
SELECT * FROM user WHERE id = ${id}