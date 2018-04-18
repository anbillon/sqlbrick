{define name CreateUser}
CREATE TABLE IF NOT EXISTS user (
  id int NOT NULL PRIMARY KEY,
  name text NOT NULL,
  age int NOT NULL
);
{end define}

-- name: SelectById
{define name SelectById}
SELECT * FROM user WHERE id = ${id}
{end define}