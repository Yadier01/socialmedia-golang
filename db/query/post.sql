-- name: GetPost :many
WITH RECURSIVE PostWithComments AS (
  SELECT p.*, u.username
  FROM posts p
  JOIN users u ON p.user_id = u.id
  WHERE p.id = $1

  UNION ALL

  SELECT p2.*, u2.username
  FROM posts p2
  JOIN PostWithComments pwc ON p2.parent_post_id = pwc.id
  JOIN users u2 ON p2.user_id = u2.id
)
SELECT *
FROM PostWithComments
ORDER BY created_at;

-- name: ListPosts :many
SELECT p.* , u.username 
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE p.parent_post_id IS NULL  -- Only select top-level posts
ORDER BY p.id
LIMIT $1
OFFSET $2;

-- name: EditPost :one
UPDATE posts
SET body = $1
WHERE id = $2 AND user_id = $3
RETURNING *;

-- name: CreatePost :one
INSERT INTO posts (
  user_id, body, parent_post_id
) VALUES (
  $1,$2, $3
) RETURNING *;


-- name: UpdatePost :one
UPDATE posts
SET comments = comments + $1  -- Use $1 for increment/decrement directly
WHERE id = $2
RETURNING *;


-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;
