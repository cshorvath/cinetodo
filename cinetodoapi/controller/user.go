package controller

import (
    "cinetodoapi/auth"
    "cinetodoapi/database"
    "cinetodoapi/model"
    "errors"

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
        return c.JSON(400, map[string]string{"error": parseErr.Error()})
    }
    var user model.User
    user.Username = login.Username
    user.HashPassword(login.Password)
    err := database.Instance.Create(&user).Error
    if err != nil {
        var mysqlErr *mysql.MySQLError
        if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
            return c.JSON(409, map[string]string{"error": "User already exists"})
        }
        return c.JSON(400, map[string]string{"error": "Invalid request"})
    }
    return c.JSON(204, nil)
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
    return c.JSON(200, user)
}
