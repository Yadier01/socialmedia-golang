-- name: FollowUser :one
INSERT INTO followers(
    follower_id, following_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: UnFollowUser :exec
DELETE FROM followers
WHERE follower_id = $1 AND following_id = $2;


