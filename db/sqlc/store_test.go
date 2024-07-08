package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFollowTx(t *testing.T) {
	store := NewStore(testDB)

	user1 := createRandomAccount(t)
	user2 := createRandomAccount(t)

	err := store.Follow(context.Background(), FollowTxParams{
		UserID:       52,
		TargetUserID: user1.ID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, user2.FollowingCount)
	require.NotEmpty(t, user1.FollowerCount)

}
func TestUnFollow(t *testing.T) {
	store := NewStore(testDB)

	// user1 := createRandomAccount(t)
	// user2 := createRandomAccount(t)

	err := store.UnFollow(context.Background(), UnFollowTxParams{
		UserID:       51,
		TargetUserID: 52,
	})

	require.NoError(t, err)
	// require.NotEmpty(t, user2.FollowingCount)
	// require.NotEmpty(t, user1.FollowerCount)

}
