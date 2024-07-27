// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: likes.sql

package db

import (
	"context"
	"database/sql"
)

const addLike = `-- name: AddLike :exec
INSERT INTO likes (user_id, post_id)
VALUES ($1, $2)
`

type AddLikeParams struct {
	UserID int64         `json:"user_id"`
	PostID sql.NullInt64 `json:"post_id"`
}

func (q *Queries) AddLike(ctx context.Context, arg AddLikeParams) error {
	_, err := q.db.ExecContext(ctx, addLike, arg.UserID, arg.PostID)
	return err
}

const getLike = `-- name: GetLike :one
SELECT id, user_id, post_id, created_at
  FROM likes l 
  WHERE user_id = $1 AND post_id = $2
`

type GetLikeParams struct {
	UserID int64         `json:"user_id"`
	PostID sql.NullInt64 `json:"post_id"`
}

func (q *Queries) GetLike(ctx context.Context, arg GetLikeParams) (Like, error) {
	row := q.db.QueryRowContext(ctx, getLike, arg.UserID, arg.PostID)
	var i Like
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PostID,
		&i.CreatedAt,
	)
	return i, err
}

const unAddLike = `-- name: UnAddLike :exec
DELETE FROM likes 
WHERE user_id = $1 AND post_id = $2
`

type UnAddLikeParams struct {
	UserID int64         `json:"user_id"`
	PostID sql.NullInt64 `json:"post_id"`
}

func (q *Queries) UnAddLike(ctx context.Context, arg UnAddLikeParams) error {
	_, err := q.db.ExecContext(ctx, unAddLike, arg.UserID, arg.PostID)
	return err
}

const updateLikesCount = `-- name: UpdateLikesCount :exec
UPDATE posts
SET likes = (
    SELECT COUNT(*)
    FROM likes
    WHERE post_id = $1
)
WHERE id = $1
`

func (q *Queries) UpdateLikesCount(ctx context.Context, postID sql.NullInt64) error {
	_, err := q.db.ExecContext(ctx, updateLikesCount, postID)
	return err
}