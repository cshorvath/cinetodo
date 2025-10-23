package main

import (
	"cinetodoapi/auth"
	"cinetodoapi/controller"
	"cinetodoapi/database"
	"cinetodoapi/tmdb"
	"cinetodoapi/view"
	"cinetodoapi/views"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Println("Running on port " + port)

	database.Connect(os.Getenv("DB_CONNECTION_STRING"))
	tmdb.InitFromEnv()
	e := echo.New()
	e.HideBanner = true
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Logger())

	renderer, err := view.NewRenderer(views.Files, nil)
	if err != nil {
		log.Fatal(err)
	}
	e.Renderer = renderer

	// Public routes
	e.GET("/login", auth.ShowLogin)
	e.POST("/login", auth.Login)
	e.POST("/logout", auth.Logout)
	e.POST("/user", controller.NewUser)

	// Protected routes
	g := e.Group("")
	g.Use(auth.JWTMiddleware())
	g.GET("/", controller.Dashboard)
	g.GET("/user", controller.GetCurrentUser)
	g.GET("/user/movie", controller.ListCurrentUserMovies)
	g.POST("/user/movie/:movieID", controller.AddMovieToUser)
	g.PATCH("/user/movie/:movieID", controller.UpdateUserMovie)
	g.DELETE("/user/movie/:movieID", controller.DeleteMovieFromUser)
	g.GET("/movie", controller.SearchMovies)

	// Public user movies by userID
	e.GET("/user/:userID/movie", controller.ListUserMovies)

	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
