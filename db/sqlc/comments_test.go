package db

import (
	"Yadier01/neon/util"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewComment(t *testing.T) {
	p, u := createPostUtil(t)

	args := createCommentParams{
		PostID:  p.ID,
		UserID:  u.ID,
		Content: util.RandomString(50),
	}
	comment, err := testQueries.createComment(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, comment)

	require.Equal(t, u.ID, comment.UserID)
	require.Equal(t, p.ID, comment.PostID)
}
