-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username;

-- name: CreateUser :one
INSERT INTO users (
  username, password, email 
) VALUES (
  $1, $2, $3 
) RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: LogIn :one
SELECT * FROM users 
WHERE username  = $1 AND password = $2 LIMIT 1;

-- name: IncrementFollowerCount :exec
UPDATE users
SET follower_count = follower_count + 1
WHERE id = $1;

-- name: IncrementFollowingCount :exec
UPDATE users
SET following_count = following_count + 1
WHERE id = $1;
