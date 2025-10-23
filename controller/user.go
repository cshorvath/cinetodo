package controller

import (
	"cinetodoapi/auth"
	"cinetodoapi/database"
	"cinetodoapi/model"
	"errors"
	"net/http"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

// NewUser creates a new user
// @Summary Create user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.Login true "User credentials"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /user [post]
func NewUser(c echo.Context) error {
	var login model.Login
	if parseErr := c.Bind(&login); parseErr != nil {
		return auth.RenderLogin(c, http.StatusBadRequest, map[string]interface{}{
			"RegisterError": "Please complete the registration form.",
		})
	}

	login.Username = strings.TrimSpace(login.Username)
	if login.Username == "" || login.Password == "" {
		return auth.RenderLogin(c, http.StatusBadRequest, map[string]interface{}{
			"RegisterError": "Both username and password are required.",
		})
	}

	var user model.User
	user.Username = login.Username
	if err := user.HashPassword(login.Password); err != nil {
		return auth.RenderLogin(c, http.StatusInternalServerError, map[string]interface{}{
			"RegisterError": "Unable to create your account right now. Please try again later.",
		})
	}
	err := database.Instance.Create(&user).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return auth.RenderLogin(c, http.StatusConflict, map[string]interface{}{
				"RegisterError": "That username is already taken. Try another one.",
				"Username":      login.Username,
			})
		}
		return auth.RenderLogin(c, http.StatusBadRequest, map[string]interface{}{
			"RegisterError": "We couldn't process your registration. Double-check the details and try again.",
		})
	}

	return auth.RenderLogin(c, http.StatusCreated, map[string]interface{}{
		"RegisterSuccess": "Account created! Sign in to start building your list.",
		"Username":        login.Username,
	})
}

// GetCurrentUser returns the current authenticated user
// @Summary Get current user
// @Tags user
// @Produce json
// @Security Bearer
// @Success 200 {object} auth.UserResponse
// @Router /user [get]
func GetCurrentUser(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	if user == nil {
		return auth.RenderLogin(c, http.StatusUnauthorized, map[string]interface{}{
			"Error": "Please sign in to continue.",
		})
	}
	return c.Render(http.StatusOK, "partials/user_badge.html", map[string]interface{}{
		"User": user,
	})
}
