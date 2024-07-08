-- name: GetPost :one
SELECT * FROM posts 
WHERE id = $1 LIMIT 1;

-- name: ListPosts :many
SELECT * FROM posts 
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: EditPost :one
UPDATE posts
SET body = $1
WHERE id = $2 AND user_id = $3  
RETURNING *;

-- name: CreatePost :one
INSERT INTO posts (
  user_id, body 
) VALUES (
  $1,$2 
) RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts 
WHERE id = $1;
