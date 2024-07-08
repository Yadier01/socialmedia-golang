package db

import (
	"Yadier01/neon/util"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createPostUtil(t *testing.T) (Post, User) {

	u := createRandomAccount(t)

	args := CreatePostParams{
		UserID: u.ID,
		Body:   util.RandomString(30),
	}

	post, err := testQueries.CreatePost(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, post)
	return post, u
}
func TestCreatePost(t *testing.T) {
	createPostUtil(t)
}

func TestListPost(t *testing.T) {

	for i := 0; i < 10; i++ {
		createPostUtil(t)
	}
	args := ListPostsParams{
		Limit:  5,
		Offset: 5,
	}

	posts, err := testQueries.ListPosts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, posts, 5)

	for _, post := range posts {
		require.NotEmpty(t, post)
	}

}

func TestDeletePost(t *testing.T) {

	p, _ := createPostUtil(t)

	err := testQueries.DeletePost(context.Background(), p.ID)
	require.NoError(t, err)

	post, err := testQueries.GetPost(context.Background(), p.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, post)

}
func TestEditPost(t *testing.T) {
	p, u := createPostUtil(t)

	args := EditPostParams{
		Body:   util.RandomString(41),
		ID:     p.ID,
		UserID: u.ID,
	}

	post, err := testQueries.EditPost(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, post)

	require.NotEqual(t, p.Body, args.Body)
	require.Equal(t, p.UserID, u.ID)

}
