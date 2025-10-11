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

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserMovieResponse struct {
	model.Movie
	Seen bool `json:"seen"`
}

type UpdateRequest struct {
	Seen bool `json:"seen"`
}

func parseIdParam(c echo.Context, param string) (ID int, err error) {
	return strconv.Atoi(c.Param(param))
}

func handleUpdateResult(c echo.Context, db *gorm.DB) {
	if db.Error != nil {
		log.Fatal(db.Error.Error())
		_ = c.JSON(500, map[string]string{"error": "Internal server error"})
		return
	}

	if db.RowsAffected == 0 {
		_ = c.JSON(404, map[string]string{"error": "Movie not found."})
		return
	}

	_ = c.JSON(204, nil)
}

func getUser(c echo.Context) *auth.UserResponse {
	return auth.GetUserFromContext(c)
}

func listUserMovies(c echo.Context, userID uint) {
	var user model.User
	dbErr := database.Instance.Preload("UserMovies.Movie").First(&user, userID).Error
	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		_ = c.JSON(404, map[string]string{"error": "User not found."})
		return
	}
	ret := make([]UserMovieResponse, 0, len(user.UserMovies))
	for _, um := range user.UserMovies {
		ret = append(ret, UserMovieResponse{Movie: um.Movie, Seen: um.Seen})
	}
	_ = c.JSON(200, ret)
}

// ListUserMovies lists movies for a given user
// @Summary List movies for user
// @Tags movie
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {array} controller.UserMovieResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /user/{userID}/movie [get]
func ListUserMovies(c echo.Context) error {
	userID, err := parseIdParam(c, "userID")
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid userID"})
	}
	listUserMovies(c, uint(userID))
	return nil
}

// ListCurrentUserMovies lists movies for current user
// @Summary List current user movies
// @Tags movie
// @Produce json
// @Security Bearer
// @Success 200 {array} controller.UserMovieResponse
// @Router /user/movie [get]
func ListCurrentUserMovies(c echo.Context) error {
	user := getUser(c)
	listUserMovies(c, user.ID)
	return nil
}

// SearchMovies searches TMDB by query
// @Summary Search movies
// @Tags movie
// @Security Bearer
// @Produce json
// @Param query query string true "Search query"
// @Success 200 {array} model.Movie
// @Router /movie [get]
func SearchMovies(c echo.Context) error {
	res, err := tmdb.Instance.SearchMovies(c.QueryParam("query"))
	if err == nil {
		return c.JSON(200, res)
	}
	return c.JSON(500, err.Error())
}

// DeleteMovieFromUser removes a movie from current user
// @Summary Delete movie from current user
// @Tags movie
// @Security Bearer
// @Param movieID path int true "Movie ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Router /user/movie/{movieID} [delete]
func DeleteMovieFromUser(c echo.Context) error {
	movieID, err := parseIdParam(c, "movieID")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}

	db := database.Instance.Where("user_id = ? AND movie_id = ?", getUser(c).ID, movieID).Delete(&model.UserMovie{})
	handleUpdateResult(c, db)
	return nil
}

// UpdateUserMovie updates seen flag for current user movie
// @Summary Update movie for current user
// @Tags movie
// @Security Bearer
// @Param movieID path int true "Movie ID"
// @Param body body controller.UpdateRequest true "Update payload"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Router /user/movie/{movieID} [patch]
func UpdateUserMovie(c echo.Context) error {
	movieID, err := parseIdParam(c, "movieID")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}

	var request UpdateRequest
	jsonErr := c.Bind(&request)

	if jsonErr != nil {
		return c.JSON(400, jsonErr.Error())
	}

	db := database.Instance.Model(&model.UserMovie{}).Where("user_id = ? AND movie_id = ?", getUser(c).ID, movieID).Update("seen", request.Seen)
	handleUpdateResult(c, db)
	return nil
}

// AddMovieToUser adds a movie to current user
// @Summary Add movie to current user
// @Tags movie
// @Security Bearer
// @Param movieID path int true "Movie ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Router /user/movie/{movieID} [post]
func AddMovieToUser(c echo.Context) error {
	movieID, err := parseIdParam(c, "movieID")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}

	movie := model.Movie{}
	database.Instance.Find(&movie, movieID)
	if movie.ID == 0 {
		tmdbMovie, err := tmdb.Instance.GetMovie(movieID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
		movie.ID = tmdbMovie.ID
		movie.Title = tmdbMovie.Title
		movie.OriginalTitle = tmdbMovie.OriginalTitle
		movie.Director = tmdbMovie.Director
		movie.Year = tmdbMovie.Year
		database.Instance.Create(&movie)
	}
	database.Instance.Create(model.UserMovie{MovieID: movie.ID, UserID: getUser(c).ID, Seen: false})
	return c.JSON(204, nil)
}
