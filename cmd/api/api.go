package api

import (
	"Yadier01/neon/cmd/token"
	db "Yadier01/neon/db/sqlc"
	"Yadier01/neon/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow your frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Add any necessary headers here
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/:id", server.GetUserById)

	r.GET("/posts", server.GetPosts)
	r.POST("/auth/register", server.CreateUser)
	r.POST("/auth/login/", server.UserLogIn)
	authRoutes := r.Group("/").Use(authMiddleware(server.TokenMaker))

	r.GET("/post/:id", server.GetPostById)

	authRoutes.POST("/add-like", server.HandleLike)
	authRoutes.POST("/post/", server.CreatePost)
	authRoutes.POST("/follow/", server.FollowUser)
	authRoutes.DELETE("/user/", server.DeleteUser)
	server.router = r
}

func (server *Server) Start() error {
	return server.router.Run()
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
