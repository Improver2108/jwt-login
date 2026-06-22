-- name: CreateUser :one
INSERT INTO
  users (
    username,
    email,
    password_hsh,
    salt,
    phone,
    first_name,
    last_name,
    bio,
    avatar_url
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING
  id,
  username,
  avatar_url;

-- name: GetUserByEmail :one
SELECT
  id,
  username,
  avatar_url,
  password_hsh,
  salt
FROM
  users
WHERE
  email = $1;

-- name: CheckUserExist :one
SELECT
  EXISTS (
    SELECT
      1
    FROM
      users
    WHERE
      email = $1
      or username = $2
      or phone = $3
  );
