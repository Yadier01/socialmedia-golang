package api

import (
	"Yadier01/neon/cmd/token"
	db "Yadier01/neon/db/sqlc"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) GetPostById(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error"})
		return
	}

	usr, err := server.store.GetPost(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "No user"})
		return
	}

	c.JSON(http.StatusOK, usr)
}

func (server *Server) CreatePost(c *gin.Context) {
	var post *db.Post
	err := c.ShouldBindJSON(&post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	args := db.CreatePostParams{
		UserID: authPayload.UserID,
		Body:   post.Body,
	}

	usr, err := server.store.CreatePost(context.Background(), args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, usr)
}

// func (server *Server) GetPosts(c *gin.Context) {
// 	post, err := server.store.ListPosts(context.Background(), )
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, err)
// 		return
// 	}

// 	c.JSON(http.StatusCreated, post)
// }
