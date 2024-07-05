-- name: FollowUser :one
INSERT INTO followers(
    follower_id, following_id
) VALUES (
    $1, $2
) RETURNING *;

