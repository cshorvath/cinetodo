package main

import (
    "cinetodoapi/auth"
    "cinetodoapi/controller"
    "cinetodoapi/database"
    "cinetodoapi/tmdb"
    "log"
    "os"

    "github.com/joho/godotenv"
    "github.com/labstack/echo/v4"
    echoMiddleware "github.com/labstack/echo/v4/middleware"
    echoSwagger "github.com/swaggo/echo-swagger"
    _ "cinetodoapi/docs"
)

// @title Cinetodo API
// @version 1.0
// @description API for managing user movie lists.
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Provide your JWT token as: Bearer <token>
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

    // Public routes
    e.POST("/login", auth.Login)
    e.POST("/user", controller.NewUser)

    // Swagger UI
    e.GET("/swagger/*", echoSwagger.WrapHandler)

    // Protected routes
    g := e.Group("")
    g.Use(auth.JWTMiddleware())
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
