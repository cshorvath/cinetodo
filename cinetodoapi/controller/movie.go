package controller

import (
	"cinetodoapi/database"
	"cinetodoapi/model"
	"cinetodoapi/tmdb"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserMovieResponse struct {
	model.Movie
	Seen bool `json:"seen"`
}

type UpdateRequest struct {
	Seen bool `json:"seen"`
}

func parseUserIDAndMovieID(c *gin.Context) (userID int, movieID int, err error) {
	movieID, err = strconv.Atoi(c.Param("movieID"))
	userID, _ = strconv.Atoi(c.Param("userID"))
	return
}

func SearchMovies(c *gin.Context) {
	res, err := tmdb.Instance.SearchMovies(c.Query("query"))
	if err == nil {
		c.JSON(200, res)
	}
	c.JSON(500, err.Error())
}

func DeleteMovieFromUser(c *gin.Context) {
	userID, movieID, err := parseUserIDAndMovieID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	database.Instance.Where("userID = ? AND movieID = ?", userID, movieID).Delete(&model.UserMovie{})
	c.JSON(204, nil)
}

func AddMovieToUser(c *gin.Context) {
	userID, movieID, err := parseUserIDAndMovieID(c)
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
		movie.Year = uint8(tmdbMovie.Year)
		database.Instance.Create(&movie)
	}
	database.Instance.Create(model.UserMovie{MovieID: int(movie.ID), UserID: userID, Seen: false})
	c.JSON(204, nil)
}
