package api

import (
	"Yadier01/neon/cmd/token"
	db "Yadier01/neon/db/sqlc"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// type PostResponse struct {
// 	PostID    int64     `json:"post_id"`
// 	UserID    int64     `json:"user_id"`
// 	Body      string    `json:"body"`
// 	Likes     int64     `json:"likes"`
// 	Comments  int64     `json:"comments"`
// 	CreatedAt time.Time `json:"created_at"`
// }

func (server *Server) GetPosts(c *gin.Context) {
	args := db.ListPostsParams{
		Limit:  5,
		Offset: 0,
	}
	posts, err := server.store.ListPosts(c, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "SOrry"})
		return
	}

	c.JSON(http.StatusOK, posts)

}

// func (server *Server) GetPostById(c *gin.Context) {
// 	idString := c.Param("id")
// 	id, err := strconv.ParseInt(idString, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error"})
// 		return
// 	}

// 	post, err := server.store.GetPost(context.Background(), sql.NullInt64{
// 		Int64: id,
// 		Valid: true,
// 	})
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
// 		} else {

// 			log.Println("Error getting post:", err) // Log the error for debugging
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get post"})
// 		}
// 		return
// 	}

// 	retrPost := PostResponse{
// 		PostID:   post.ID,
// 		UserID:   post.UserID,
// 		Body:     post.Body,
// 		Comments: post.Comments,
// 	}
// 	c.JSON(http.StatusOK, retrPost)
// }

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
