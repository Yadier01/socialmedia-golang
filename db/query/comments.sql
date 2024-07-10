-- name: createComment :one
INSERT INTO comments(
  post_id, user_id, content, likes,  parent_comment_id
) VALUES (
   $1, $2, $3, $4, $5
)RETURNING *;



