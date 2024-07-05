-- name: GetPost :one
SELECT * FROM posts 
WHERE id = $1 LIMIT 1;

-- name: ListPosts :many
SELECT * FROM posts 
ORDER BY id;

-- name: CreatePost :one
INSERT INTO posts (
  user_id, title, body 
) VALUES (
  $1, $2, $3 
) RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts 
WHERE id = $1;