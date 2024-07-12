-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1
LIMIT 1;

-- name: ListPosts :many
select * from posts
order by id
limit $1
offset $2;

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
