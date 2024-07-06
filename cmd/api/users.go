package api

import (
	"Yadier01/neon/cmd/token"
	db "Yadier01/neon/db/sqlc"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) GetUserById(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error"})
		return
	}

	usr, err := server.store.GetUser(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "No user"})
		return
	}
	c.JSON(http.StatusOK, usr)
}

func (server *Server) GetUsers(c *gin.Context) {
	post, err := server.store.ListUsers(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (server *Server) FollowUser(c *gin.Context) {
	var follow *db.Follower
	err := c.ShouldBindJSON(&follow)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.FollowTxParams{
		UserID:       authPayload.UserID,
		TargetUserID: follow.FollowingID,
	}

	err = server.store.Follow(c, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "all good"})

}

func (server *Server) AddUser(c *gin.Context) {
	var u *db.User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	args := db.CreateUserParams{
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
	}
	usr, err := server.store.CreateUser(context.Background(), args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		fmt.Print(err)
		return
	}

	jwt, err := server.TokenMaker.CreateToken(usr.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		fmt.Print(err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": usr, "jwt": jwt})

}

func (server *Server) UserLogIn(c *gin.Context) {
	var u *db.User
	err := c.ShouldBindJSON(&u)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	usr, err := server.store.LogIn(context.Background(), db.LogInParams{Username: u.Username, Password: u.Password})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	jwt, err := server.TokenMaker.CreateToken(usr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		fmt.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"usr": usr, "token": jwt})

}
