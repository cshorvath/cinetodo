package controller

import (
	"cinetodoapi/auth"
	"cinetodoapi/database"
	"cinetodoapi/model"
	"cinetodoapi/tmdb"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserMovieResponse struct {
	model.Movie
	Seen bool `json:"seen"`
}

type UpdateRequest struct {
	Seen bool `json:"seen"`
}

func parseIdParam(c *gin.Context, param string) (ID int, err error) {
	return strconv.Atoi(c.Param(param))
}

func handleUpdateResult(c *gin.Context, db *gorm.DB) {
	if db.Error != nil {
		log.Fatal(db.Error.Error())
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	if db.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Movie not found."})
		return
	}

	c.JSON(204, nil)
}

func getUser(c *gin.Context) *auth.UserResponse {
	user, ok := c.Get(auth.IdentityKey)
	if !ok {
		return nil
	}
	return user.(*auth.UserResponse)
}

func matchUser(c *gin.Context) bool {
	user := getUser(c)
	pathUserID, err := parseIdParam(c, "userID")
	if err == nil && user != nil && user.ID == uint(pathUserID) {
		return true
	}
	c.JSON(403, gin.H{"error": "Unauthorized."})
	return false
}

func listUserMovies(c *gin.Context, userID uint) {
	var user model.User
	dbErr := database.Instance.Preload("UserMovies.Movie").First(&user, userID).Error
	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{"error": "User not found."})
		return
	}
	ret := make([]UserMovieResponse, 0, len(user.UserMovies))
	for _, um := range user.UserMovies {
		ret = append(ret, UserMovieResponse{Movie: um.Movie, Seen: um.Seen})
	}
	c.JSON(200, ret)
}

func ListUserMovies(c *gin.Context) {
	userID, err := parseIdParam(c, "userID")
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid userID"})
		return
	}
	listUserMovies(c, uint(userID))
}

func ListCurrentUserMovies(c *gin.Context) {
	user := getUser(c)
	listUserMovies(c, user.ID)
}

func SearchMovies(c *gin.Context) {
	res, err := tmdb.Instance.SearchMovies(c.Query("query"))
	if err == nil {
		c.JSON(200, res)
		return
	}
	c.JSON(500, err.Error())
}

func DeleteMovieFromUser(c *gin.Context) {
	if !matchUser(c) {
		return
	}
	movieID, err := parseIdParam(c, "movieID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	db := database.Instance.Where("user_id = ? AND movie_id = ?", getUser(c).ID, movieID).Delete(&model.UserMovie{})
	handleUpdateResult(c, db)
}

func UpdateUserMovie(c *gin.Context) {
	if !matchUser(c) {
		return
	}
	movieID, err := parseIdParam(c, "movieID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	var request UpdateRequest
	jsonErr := c.ShouldBindJSON(&request)

	if jsonErr != nil {
		c.JSON(400, jsonErr.Error())
		return
	}

	db := database.Instance.Model(&model.UserMovie{}).Where("user_id = ? AND movie_id = ?", getUser(c).ID, movieID).Update("seen", request.Seen)
	handleUpdateResult(c, db)
}

func AddMovieToUser(c *gin.Context) {
	if !matchUser(c) {
		return
	}
	movieID, err := parseIdParam(c, "movieID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	movie := model.Movie{}
	database.Instance.Find(&movie, movieID)
	if movie.ID == 0 {
		tmdbMovie, err := tmdb.Instance.GetMovie(movieID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		movie.ID = tmdbMovie.ID
		movie.Title = tmdbMovie.Title
		movie.OriginalTitle = tmdbMovie.OriginalTitle
		movie.Director = tmdbMovie.Director
		movie.Year = tmdbMovie.Year
		database.Instance.Create(&movie)
	}
	database.Instance.Create(model.UserMovie{MovieID: movie.ID, UserID: getUser(c).ID, Seen: false})
	c.JSON(204, nil)
}
