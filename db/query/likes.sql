-- name: AddLike :exec
INSERT INTO likes (user_id, post_id)
VALUES ($1, $2);

-- name: UnAddLike :exec
DELETE FROM likes 
WHERE user_id = $1 AND post_id = $2;

-- name: UpdateLikesCount :exec
UPDATE posts
SET likes = (
    SELECT COUNT(*)
    FROM likes
    WHERE post_id = $1
)
WHERE id = $1;
-- name: GetLike :one
SELECT *
  FROM likes l 
  WHERE user_id = $1 AND post_id = $2;
