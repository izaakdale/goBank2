-- name: CreateUser :one
INSERT INTO users (
  username, 
  hashed_password, 
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
)
RETURNING username, email, created_at;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;