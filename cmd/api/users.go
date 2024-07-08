package api

import (
	"Yadier01/neon/cmd/token"
	db "Yadier01/neon/db/sqlc"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func (server *Server) GetUserById(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error"})
		return
	}

	usr, err := server.store.GetUser(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "No user"})
		return
	}
	c.JSON(http.StatusOK, usr)
}

func (server *Server) GetUsers(c *gin.Context) {
	post, err := server.store.ListUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, post)
}

// this creates a db transaction
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

func (server *Server) CreateUser(c *gin.Context) {
	var req *db.User
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 16)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	args := db.CreateUserParams{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	usr, err := server.store.CreateUser(context.Background(), args)
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
	c.JSON(http.StatusCreated, gin.H{"user": usr, "jwt": jwt})

}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (server *Server) UserLogIn(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user, err := server.store.LogIn(c, req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong credentials"})
		return
	}

	jwt, err := server.TokenMaker.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"usr": req.Username, "token": jwt})
}

func (server *Server) DeleteUser(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	err := server.store.DeleteUser(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
