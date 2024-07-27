package api

import (
	"Yadier01/neon/cmd/token"
	db "Yadier01/neon/db/sqlc"
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PostResponse struct {
	Usernme   string         `json:"username"`
	PostID    int64          `json:"post_id"`
	UserID    int64          `json:"user_id"`
	Body      string         `json:"body"`
	Likes     int64          `json:"likes"`
	Comments  int64          `json:"comments"`
	Reply     []PostResponse `json:"reply"`
	CreatedAt time.Time      `json:"created_at"`
}

// recursive function to build the reply tree
func buildReplyTree(parent *PostResponse, rows []db.GetPostRow) {
	for _, row := range rows {

		if row.ParentPostID.Valid && row.ParentPostID.Int64 == parent.PostID {
			reply := &PostResponse{
				PostID:    row.ID,
				UserID:    row.UserID,
				Body:      row.Body,
				Likes:     row.Likes,
				Comments:  row.Comments,
				Usernme:   row.Username,
				Reply:     []PostResponse{},
				CreatedAt: row.CreatedAt,
			}

			parent.Reply = append(parent.Reply, *reply)
			buildReplyTree(reply, rows)
		}
	}
}

func (server *Server) GetPostById(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	rows, err := server.store.GetPost(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
		return
	}

	if len(rows) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// The root post is always the first row because of the ORDER BY in the SQL query
	rootPostRow := rows[0]

	rootPost := &PostResponse{
		PostID:    rootPostRow.ID,
		UserID:    rootPostRow.UserID,
		Body:      rootPostRow.Body,
		Likes:     rootPostRow.Likes,
		Usernme:   rootPostRow.Username,
		Comments:  rootPostRow.Comments,
		Reply:     []PostResponse{},
		CreatedAt: rootPostRow.CreatedAt,
	}

	buildReplyTree(rootPost, rows)

	c.JSON(http.StatusOK, rootPost)
}

func (server *Server) GetPosts(c *gin.Context) {
	args := db.ListPostsParams{
		Limit:  10,
		Offset: 0,
	}
	posts, err := server.store.ListPosts(c, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusOK, posts)

}

func (server *Server) CreatePost(c *gin.Context) {
	var upost *db.Post
	err := c.ShouldBindJSON(&upost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.CreatePostParams{
		UserID:       authPayload.UserID,
		Body:         upost.Body,
		ParentPostID: upost.ParentPostID,
	}

	post, err := server.store.CreatePost(context.Background(), args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	//this will update the comment count
	if post.ParentPostID.Valid {
		_, err := server.store.UpdatePost(c, db.UpdatePostParams{
			ID:       post.ParentPostID.Int64,
			Comments: 1,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusCreated, post)
}

type ReceiveHandleLike struct {
	PostID int64 `json:"post_id"`
}

// TODO: make this tx
func (server *Server) HandleLike(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	var handleLike ReceiveHandleLike
	if err := c.ShouldBindJSON(&handleLike); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.GetLikeParams{
		UserID: authPayload.UserID,
		PostID: sql.NullInt64{Int64: handleLike.PostID, Valid: true},
	}

	like, err := server.store.GetLike(c, args)
	if err != nil {
		if err == sql.ErrNoRows {
			server.handleAddLike(c, authPayload.UserID, handleLike.PostID)
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.handleRemoveLike(c, authPayload.UserID, like.PostID)
}

func (server *Server) handleAddLike(c *gin.Context, userID int64, postID int64) {
	err := server.store.AddLike(c, db.AddLikeParams{
		UserID: userID,
		PostID: sql.NullInt64{Int64: postID, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := server.store.UpdateLikesCount(c, sql.NullInt64{Int64: postID, Valid: true}); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Liked post"})
}

func (server *Server) handleRemoveLike(c *gin.Context, userID int64, postID sql.NullInt64) {
	err := server.store.UnAddLike(c, db.UnAddLikeParams{
		UserID: userID,
		PostID: postID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to unlike post"})
		return
	}

	if err := server.store.UpdateLikesCount(c, sql.NullInt64{Int64: postID.Int64, Valid: true}); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unliked post"})
}
