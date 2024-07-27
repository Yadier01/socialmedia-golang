// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: post.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (
  user_id, body, parent_post_id
) VALUES (
  $1,$2, $3
) RETURNING id, user_id, body, likes, comments, parent_post_id, created_at
`

type CreatePostParams struct {
	UserID       int64         `json:"user_id"`
	Body         string        `json:"body"`
	ParentPostID sql.NullInt64 `json:"parent_post_id"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost, arg.UserID, arg.Body, arg.ParentPostID)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Body,
		&i.Likes,
		&i.Comments,
		&i.ParentPostID,
		&i.CreatedAt,
	)
	return i, err
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1
`

func (q *Queries) DeletePost(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePost, id)
	return err
}

const editPost = `-- name: EditPost :one
UPDATE posts
SET body = $1
WHERE id = $2 AND user_id = $3
RETURNING id, user_id, body, likes, comments, parent_post_id, created_at
`

type EditPostParams struct {
	Body   string `json:"body"`
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
}

func (q *Queries) EditPost(ctx context.Context, arg EditPostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, editPost, arg.Body, arg.ID, arg.UserID)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Body,
		&i.Likes,
		&i.Comments,
		&i.ParentPostID,
		&i.CreatedAt,
	)
	return i, err
}

const getPost = `-- name: GetPost :many
WITH RECURSIVE PostWithComments AS (
  SELECT p.id, p.user_id, p.body, p.likes, p.comments, p.parent_post_id, p.created_at, u.username
  FROM posts p
  JOIN users u ON p.user_id = u.id
  WHERE p.id = $1

  UNION ALL

  SELECT p2.id, p2.user_id, p2.body, p2.likes, p2.comments, p2.parent_post_id, p2.created_at, u2.username
  FROM posts p2
  JOIN PostWithComments pwc ON p2.parent_post_id = pwc.id
  JOIN users u2 ON p2.user_id = u2.id
)
SELECT id, user_id, body, likes, comments, parent_post_id, created_at, username
FROM PostWithComments
ORDER BY created_at
`

type GetPostRow struct {
	ID           int64         `json:"id"`
	UserID       int64         `json:"user_id"`
	Body         string        `json:"body"`
	Likes        int64         `json:"likes"`
	Comments     int64         `json:"comments"`
	ParentPostID sql.NullInt64 `json:"parent_post_id"`
	CreatedAt    time.Time     `json:"created_at"`
	Username     string        `json:"username"`
}

func (q *Queries) GetPost(ctx context.Context, id int64) ([]GetPostRow, error) {
	rows, err := q.db.QueryContext(ctx, getPost, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPostRow{}
	for rows.Next() {
		var i GetPostRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Body,
			&i.Likes,
			&i.Comments,
			&i.ParentPostID,
			&i.CreatedAt,
			&i.Username,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPosts = `-- name: ListPosts :many
SELECT p.id, p.user_id, p.body, p.likes, p.comments, p.parent_post_id, p.created_at , u.username 
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE p.parent_post_id IS NULL  -- Only select top-level posts
ORDER BY p.id
LIMIT $1
OFFSET $2
`

type ListPostsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListPostsRow struct {
	ID           int64         `json:"id"`
	UserID       int64         `json:"user_id"`
	Body         string        `json:"body"`
	Likes        int64         `json:"likes"`
	Comments     int64         `json:"comments"`
	ParentPostID sql.NullInt64 `json:"parent_post_id"`
	CreatedAt    time.Time     `json:"created_at"`
	Username     string        `json:"username"`
}

func (q *Queries) ListPosts(ctx context.Context, arg ListPostsParams) ([]ListPostsRow, error) {
	rows, err := q.db.QueryContext(ctx, listPosts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListPostsRow{}
	for rows.Next() {
		var i ListPostsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Body,
			&i.Likes,
			&i.Comments,
			&i.ParentPostID,
			&i.CreatedAt,
			&i.Username,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePost = `-- name: UpdatePost :one
UPDATE posts
SET comments = comments + $1  -- Use $1 for increment/decrement directly
WHERE id = $2
RETURNING id, user_id, body, likes, comments, parent_post_id, created_at
`

type UpdatePostParams struct {
	Comments int64 `json:"comments"`
	ID       int64 `json:"id"`
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, updatePost, arg.Comments, arg.ID)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Body,
		&i.Likes,
		&i.Comments,
		&i.ParentPostID,
		&i.CreatedAt,
	)
	return i, err
}
