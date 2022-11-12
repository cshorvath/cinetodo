package main

import (
	"cinetodoapi/auth"
	"cinetodoapi/controller"
	"cinetodoapi/database"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	r := gin.Default()

	if port == "" {
		port = "8000"
	}

	database.Connect(os.Getenv("DB_CONNECTION_STRING"))

	authMiddleware := auth.InitAuthMiddleware()
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Not found"})
	})

	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/refreshToken", authMiddleware.RefreshHandler)
	r.POST("/user", controller.NewUser)

	auth := r.Group("/")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/user", controller.GetCurrentUser)
		auth.GET("/user/movie", controller.ListCurrentUserMovies)
		auth.POST("/user/movie/:movieID", controller.AddMovieToUser)
		auth.PATCH("/user/movie/:movieID", controller.UpdateUserMovie)
		auth.DELETE("/user/movie/:movieID", controller.DeleteMovieFromUser)
		auth.GET("/movie", controller.SearchMovies)
	}

	r.GET("/user/:userID/movie", controller.ListUserMovies)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
