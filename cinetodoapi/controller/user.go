package controller

import (
	"cinetodoapi/database"
	"cinetodoapi/model"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func NewUser(c *gin.Context) {
	var login model.Login
	if parseErr := c.ShouldBindJSON(&login); parseErr != nil {
		c.JSON(400, parseErr.Error())
		return
	}
	var user model.User
	user.Username = login.Username
	user.HashPassword(login.Password)
	err := database.Instance.Create(&user).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			c.JSON(409, gin.H{"error": "User already exists"})
			return
		}
		c.JSON(400, gin.H{"error": "Invalid request"})
	}
	c.JSON(204, nil)
}

func GetCurrentUser(c *gin.Context) {
	user, _ := c.Get("ID")
	c.JSON(200, &user)
}
