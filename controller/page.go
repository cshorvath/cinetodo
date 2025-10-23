package controller

import (
	"cinetodoapi/auth"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Dashboard renders the main application page for the authenticated user.
func Dashboard(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	if user == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	movies, err := fetchUserMovies(user.ID)
	loadFailure := err != nil && !errors.Is(err, gorm.ErrRecordNotFound)
	if loadFailure {
		log.Printf("failed to preload dashboard movies for user %d: %v", user.ID, err)
		// fall back to empty list but still render the page with error message flag.
		movies = []UserMovieResponse{}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		movies = []UserMovieResponse{}
	}

	data := map[string]interface{}{
		"User":    user,
		"Movies":  movies,
		"IsOwner": true,
	}
	if loadFailure {
		data["InitialFlash"] = map[string]string{
			"Kind":    "error",
			"Message": "We couldn't load your movies. The list below may be incomplete.",
		}
	}

	return c.Render(http.StatusOK, "pages/dashboard.html", data)
}
