package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	Follow(ctx context.Context, args FollowTxParams) error
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execute a funciton witin a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:  %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type FollowTxParams struct {
	UserID       int64 // The user initiating the follow
	TargetUserID int64 // The user being followed
}

// increments the following and follower for both users
func (store *SQLStore) Follow(ctx context.Context, args FollowTxParams) error {
	return store.execTx(ctx, func(q *Queries) error {
		if err := q.IncrementFollowingCount(ctx, args.UserID); err != nil {
			return err
		}
		if err := q.IncrementFollowerCount(ctx, args.TargetUserID); err != nil {
			return err
		}

		_, err := q.FollowUser(ctx, FollowUserParams{
			FollowerID:  args.UserID,
			FollowingID: args.TargetUserID,
		})

		return err
	})
}
