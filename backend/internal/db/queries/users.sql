-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  password_hash,
  display_name
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET display_name = $2,
    avatar_url = $3,
    bio = $4,
    is_public = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
