package db

import (
	"Yadier01/neon/util"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) User {

	args := CreateUserParams{
		Username: util.RandomUser(),
		Password: util.RandomPassword(),
		Email:    util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.Password, user.Password)
	require.Equal(t, args.Email, user.Email)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	return user
}
func TestCreateUser(t *testing.T) {
	createRandomAccount(t)
}

func TestGetUser(t *testing.T) {
	u := createRandomAccount(t)
	usr, err := testQueries.GetUser(context.Background(), u.ID)
	require.NoError(t, err)
	require.NotEmpty(t, usr)

	require.Equal(t, u.ID, usr.ID)
	require.Equal(t, u.Password, usr.Password)
	require.Equal(t, u.Email, usr.Email)
	require.Equal(t, u.FollowerCount, usr.FollowerCount)
	require.Equal(t, u.FollowingCount, usr.FollowingCount)

}

func TestUpdateAccount(t *testing.T) {
	u := createRandomAccount(t)

	args := UpdateUserParams{
		Username: util.RandomUser(),
		Password: util.RandomPassword(),
		Email:    util.RandomEmail(),
		ID:       u.ID,
	}
	user, err := testQueries.UpdateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotEqual(t, u.Username, user.Username)
	require.NotEqual(t, u.Password, user.Password)
}

func TestDeleteAccount(t *testing.T) {

	u := createRandomAccount(t)

	err := testQueries.DeleteUser(context.Background(), u.ID)
	require.NoError(t, err)

	usr, err := testQueries.GetUser(context.Background(), u.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, usr)

}
