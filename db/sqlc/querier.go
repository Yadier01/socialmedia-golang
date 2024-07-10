// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
)

type Querier interface {
	CreatePost(ctx context.Context, arg CreatePostParams) (Post, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DecreaseFollowerCount(ctx context.Context, id int64) error
	DecreaseFollowingCount(ctx context.Context, id int64) error
	DeletePost(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	EditPost(ctx context.Context, arg EditPostParams) (Post, error)
	FollowUser(ctx context.Context, arg FollowUserParams) (Follower, error)
	GetPost(ctx context.Context, id int64) (Post, error)
	GetUser(ctx context.Context, id int64) (User, error)
	IncrementFollowerCount(ctx context.Context, id int64) error
	IncrementFollowingCount(ctx context.Context, id int64) error
	ListPosts(ctx context.Context, arg ListPostsParams) ([]Post, error)
	ListUsers(ctx context.Context) ([]User, error)
	LogIn(ctx context.Context, username string) (User, error)
	UnFollowUser(ctx context.Context, arg UnFollowUserParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	createComment(ctx context.Context, arg createCommentParams) (Comment, error)
}

var _ Querier = (*Queries)(nil)
