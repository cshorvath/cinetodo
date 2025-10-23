package controller

import (
	"cinetodoapi/auth"
	"cinetodoapi/database"
	"cinetodoapi/model"
	"cinetodoapi/tmdb"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserMovieResponse struct {
	model.Movie
	Seen bool `json:"seen"`
}

type FlashMessage struct {
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

func parseIdParam(c echo.Context, param string) (ID int, err error) {
	return strconv.Atoi(c.Param(param))
}

func getUser(c echo.Context) *auth.UserResponse {
	return auth.GetUserFromContext(c)
}

func isHTMXRequest(c echo.Context) bool {
	return strings.EqualFold(c.Request().Header.Get("HX-Request"), "true")
}

func fetchUserMovies(userID uint) ([]UserMovieResponse, error) {
	var user model.User
	dbErr := database.Instance.Preload("UserMovies.Movie").First(&user, userID).Error
	if dbErr != nil {
		return nil, dbErr
	}
	ret := make([]UserMovieResponse, 0, len(user.UserMovies))
	for _, um := range user.UserMovies {
		ret = append(ret, UserMovieResponse{Movie: um.Movie, Seen: um.Seen})
	}
	return ret, nil
}

func renderWatchlist(c echo.Context, status int, movies []UserMovieResponse, isOwner bool, flash *FlashMessage) error {
	data := map[string]interface{}{
		"Movies":  movies,
		"IsOwner": isOwner,
	}
	if flash != nil {
		data["Flash"] = flash
	}
	return c.Render(status, "partials/watchlist.html", data)
}

func renderWatchlistForUser(c echo.Context, userID uint, isOwner bool, flash *FlashMessage) error {
	movies, err := fetchUserMovies(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			movies = []UserMovieResponse{}
		} else {
			log.Printf("failed to load movies for user %d: %v", userID, err)
			if flash == nil {
				flash = &FlashMessage{
					Kind:    "error",
					Message: "We couldn't load the movie list. Please try again.",
				}
			}
			movies = []UserMovieResponse{}
		}
	}
	return renderWatchlist(c, http.StatusOK, movies, isOwner, flash)
}

func ListUserMovies(c echo.Context) error {
	userID, err := parseIdParam(c, "userID")
	if err != nil {
		flash := &FlashMessage{Kind: "error", Message: "Invalid user id."}
		if isHTMXRequest(c) {
			return renderWatchlist(c, http.StatusBadRequest, []UserMovieResponse{}, false, flash)
		}
		return c.Render(http.StatusBadRequest, "pages/public_watchlist.html", map[string]interface{}{
			"Movies":  []UserMovieResponse{},
			"Flash":   flash,
			"Title":   "Shared Watchlist",
			"IsOwner": false,
		})
	}

	var owner model.User
	ownerErr := database.Instance.First(&owner, userID).Error
	if ownerErr != nil {
		flash := &FlashMessage{Kind: "warning", Message: "That user could not be found."}
		if isHTMXRequest(c) {
			return renderWatchlist(c, http.StatusOK, []UserMovieResponse{}, false, flash)
		}
		return c.Render(http.StatusOK, "pages/public_watchlist.html", map[string]interface{}{
			"Movies":  []UserMovieResponse{},
			"Flash":   flash,
			"Title":   "Shared Watchlist",
			"IsOwner": false,
		})
	}

	movies, fetchErr := fetchUserMovies(uint(userID))
	var flash *FlashMessage
	if fetchErr != nil {
		if errors.Is(fetchErr, gorm.ErrRecordNotFound) {
			movies = []UserMovieResponse{}
		} else {
			log.Printf("failed to load public user movies: %v", fetchErr)
			flash = &FlashMessage{Kind: "error", Message: "Unable to load movies right now."}
			movies = []UserMovieResponse{}
		}
	}

	if isHTMXRequest(c) {
		return renderWatchlist(c, http.StatusOK, movies, false, flash)
	}

	title := fmt.Sprintf("%s's Watchlist", owner.Username)
	data := map[string]interface{}{
		"Movies":  movies,
		"Owner":   owner.Username,
		"Title":   title,
		"IsOwner": false,
	}
	if flash != nil {
		data["Flash"] = flash
	}
	return c.Render(http.StatusOK, "pages/public_watchlist.html", data)
}

func ListCurrentUserMovies(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		return renderWatchlist(c, http.StatusUnauthorized, []UserMovieResponse{}, true, &FlashMessage{
			Kind:    "warning",
			Message: "Please sign in again.",
		})
	}
	return renderWatchlistForUser(c, user.ID, true, nil)
}

func SearchMovies(c echo.Context) error {
	query := strings.TrimSpace(c.QueryParam("query"))
	if query == "" {
		return c.Render(http.StatusOK, "partials/search_results.html", map[string]interface{}{
			"Results": []*model.Movie{},
			"Query":   query,
		})
	}

	res, err := tmdb.Instance.SearchMovies(query)
	if err != nil {
		log.Printf("search error: %v", err)
		return c.Render(http.StatusInternalServerError, "partials/search_results.html", map[string]interface{}{
			"Results": []*model.Movie{},
			"Query":   query,
			"Error":   "We couldn't search TMDB right now. Try again later.",
		})
	}
	return c.Render(http.StatusOK, "partials/search_results.html", map[string]interface{}{
		"Results": res,
		"Query":   query,
	})
}

func DeleteMovieFromUser(c echo.Context) error {
	movieID, err := parseIdParam(c, "movieID")
	user := getUser(c)
	if err != nil {
		return renderWatchlistForUser(c, user.ID, true, &FlashMessage{
			Kind:    "error",
			Message: "Invalid movie id.",
		})
	}
	db := database.Instance.Where("user_id = ? AND movie_id = ?", user.ID, movieID).Delete(&model.UserMovie{})
	if db.Error != nil {
		log.Printf("failed to delete movie %d for user %d: %v", movieID, user.ID, db.Error)
		return renderWatchlistForUser(c, user.ID, true, &FlashMessage{
			Kind:    "error",
			Message: "We couldn't remove that movie. Please try again.",
		})
	}
	if db.RowsAffected == 0 {
		return renderWatchlistForUser(c, user.ID, true, &FlashMessage{
			Kind:    "warning",
			Message: "Movie was not found in your list.",
		})
	}
	return renderWatchlistForUser(c, user.ID, true, &FlashMessage{
		Kind:    "success",
		Message: "Removed the movie from your watchlist.",
	})
}

func UpdateUserMovie(c echo.Context) error {
	movieID, err := parseIdParam(c, "movieID")
	user := getUser(c)
	fmt.Println(user.ID)
	if err != nil {
		return renderWatchlistForUser(c, user.ID, true, &FlashMessage{
			Kind:    "error",
			Message: "Invalid movie id.",
		})
	}
	seen, err := strconv.ParseBool(c.FormValue("seen"))
	if err != nil {
		return renderWatchlistForUser(c, user.ID, true, &FlashMessage{
			Kind:    "error",
			Message: "Could not read your update.",
		})
	}

	db := database.Instance.Model(&model.UserMovie{}).
		Where("user_id = ? AND movie_id = ?", user.ID, movieID).
		Update("seen", seen)
	if db.Error != nil {
		log.Printf("failed to update movie %d for user %d: %v", movieID, user.ID, db.Error)
		return renderWatchlistForUser(c, user.ID, true, &FlashMessage{
			Kind:    "error",
			Message: "We couldn't update that movie. Please try again.",
		})
	}
	if db.RowsAffected == 0 {
		return renderWatchlistForUser(c, user.ID, true, &FlashMessage{
			Kind:    "warning",
			Message: "Movie was not found in your list.",
		})
	}
	status := &FlashMessage{Kind: "success", Message: "Marked as unseen."}
	if seen {
		status.Message = "Marked as seen."
	}
	return renderWatchlistForUser(c, user.ID, true, status)
}

func AddMovieToUser(c echo.Context) error {
	movieID, err := parseIdParam(c, "movieID")
	user := getUser(c)
	if err != nil {
		return renderWatchlistForUser(c, user.ID, true, &FlashMessage{
			Kind:    "error",
			Message: "Invalid movie id.",
		})
	}
	movie := model.Movie{}
	database.Instance.Find(&movie, movieID)
	if movie.ID == 0 {
		tmdbMovie, tmdbErr := tmdb.Instance.GetMovie(movieID)
		if tmdbErr != nil {
			log.Printf("tmdb fetch failed for %d: %v", movieID, tmdbErr)
			return renderWatchlistForUser(c, user.ID, true, &FlashMessage{
				Kind:    "error",
				Message: "We couldn't fetch that movie from TMDB.",
			})
		}
		movie.ID = tmdbMovie.ID
		movie.Title = tmdbMovie.Title
		movie.OriginalTitle = tmdbMovie.OriginalTitle
		movie.Director = tmdbMovie.Director
		movie.Year = tmdbMovie.Year
		movie.PosterPath = tmdbMovie.PosterPath
		if createErr := database.Instance.Create(&movie).Error; createErr != nil {
			log.Printf("failed to create movie %d: %v", movieID, createErr)
			return renderWatchlistForUser(c, user.ID, true, &FlashMessage{
				Kind:    "error",
				Message: "We couldn't save that movie. Try again.",
			})
		}
	}
	if err := database.Instance.Create(&model.UserMovie{MovieID: movie.ID, UserID: user.ID, Seen: false}).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return renderWatchlistForUser(c, user.ID, true, &FlashMessage{
				Kind:    "warning",
				Message: "That movie is already on your list.",
			})
		}
		log.Printf("failed to attach movie %d to user %d: %v", movie.ID, user.ID, err)
		return renderWatchlistForUser(c, user.ID, true, &FlashMessage{
			Kind:    "error",
			Message: "We couldn't add that movie right now.",
		})
	}
	success := &FlashMessage{
		Kind:    "success",
		Message: fmt.Sprintf("Added \"%s\" to your watchlist.", movie.Title),
	}
	c.Response().Header().Set("HX-Trigger", `{"search-reset": ""}`)
	return renderWatchlistForUser(c, user.ID, true, success)
}
