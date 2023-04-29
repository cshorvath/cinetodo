package auth

import (
	"cinetodoapi/database"
	"cinetodoapi/model"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

const IdentityKey = "user"

func InitAuthMiddleware() *jwt.GinJWTMiddleware {

	ret, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "cinetodo",
		Key:         []byte(os.Getenv("JWT_SECRET")),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					"ID": v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			var user model.User
			database.Instance.First(&user, claims["ID"])
			return &UserResponse{ID: user.ID, Username: user.Username}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals model.Login
			if err := c.ShouldBindJSON(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			var user model.User
			database.Instance.Where("username = ?", loginVals.Username).First(&user)
			if user.CheckPassword(loginVals.Password) == nil {
				return &user, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		panic("JWT Error:" + err.Error())
	}

	return ret

}
