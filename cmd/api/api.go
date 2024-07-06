package api

import (
	"Yadier01/neon/cmd/token"
	db "Yadier01/neon/db/sqlc"
	"Yadier01/neon/util"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Queries    *db.Queries
	store      db.Store
	router     *gin.Engine
	TokenMaker token.Maker
	Config     util.Config
}


func NewServer(config util.Config, store db.Store) *Server {
	jwt, err := token.NewJWTMaker(config.TokenSymmetricKey)

	    if err != nil {
		log.Fatal(err)
	}
	server := &Server{
		store:      store,
		TokenMaker: jwt,
		Config:     config,
	}
	server.setupRouter()
	return server
}

func (server *Server) setupRouter() {
	r := gin.Default()

	r.POST("/login", server.UserLogIn)
	r.GET("/:id", server.GetUserById)
	r.GET("/", server.GetUsers)
	r.POST("/", server.AddUser)
	authRoutes := r.Group("/").Use(authMiddleware(server.TokenMaker))

	r.GET("/post/:id", server.GetPostById)
	r.GET("/post/", server.GetPosts)
	authRoutes.POST("/post/", server.CreatePost)
	authRoutes.POST("/follow/", server.FollowUser)
	server.router = r
}
func (server *Server) Start() error {
	return server.router.Run()
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
