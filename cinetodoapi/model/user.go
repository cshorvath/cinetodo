package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserTokenPayload struct {
	ID uint
}

type User struct {
	gorm.Model
	Username     string
	PasswordHash string
	Movies       []Movie `gorm:"many2many:user_movies;"`
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.PasswordHash = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
