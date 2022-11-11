package auth

import (
	"cinetodoapi/database"
	"cinetodoapi/model"
	"os"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func InitAuthMiddleware() *jwt.GinJWTMiddleware {
	const identityKey = "ID"

	ret, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "cinetodo",
		Key:         []byte(os.Getenv("JWT_SECRET")),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.UserTokenPayload); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model.UserTokenPayload{
				ID: claims[identityKey].(uint),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals model.Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			var user model.User

			database.Instance.Where("username = ?", username, user)
			if user.CheckPassword(password) == nil {
				return model.UserTokenPayload{ID: user.ID}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			requestID, err := strconv.Atoi(c.Param("userID"))
			if v, ok := data.(*model.UserTokenPayload); ok && err == nil && int(v.ID) == requestID {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		panic("JWT Error:" + err.Error())
	}

	return ret

}
