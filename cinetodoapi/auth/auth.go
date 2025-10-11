package auth

import (
    "cinetodoapi/database"
    "cinetodoapi/model"
    "net/http"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/labstack/echo/v4"
)

type UserResponse struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
}

const IdentityKey = "user"

// GetUserFromContext returns the authenticated user from context
func GetUserFromContext(c echo.Context) *UserResponse {
    v := c.Get(IdentityKey)
    if v == nil {
        return nil
    }
    user, ok := v.(*UserResponse)
    if !ok {
        return nil
    }
    return user
}

// Login authenticates user and issues JWT token
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body model.Login true "Login"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func Login(c echo.Context) error {
    var loginVals model.Login
    if err := c.Bind(&loginVals); err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing credentials"})
    }
    var user model.User
    database.Instance.Where("username = ?", loginVals.Username).First(&user)
    if user.ID == 0 || user.CheckPassword(loginVals.Password) != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
    }

    claims := jwt.MapClaims{
        "ID":  user.ID,
        "exp": time.Now().Add(time.Hour).Unix(),
        "iat": time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "token error"})
    }
    return c.JSON(http.StatusOK, map[string]string{"token": signed})
}

// JWTMiddleware validates JWT and loads user into context
func JWTMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            authHeader := c.Request().Header.Get("Authorization")
            if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing bearer token"})
            }
            tokenStr := authHeader[7:]
            token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
                if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid signing method")
                }
                return []byte(os.Getenv("JWT_SECRET")), nil
            })
            if err != nil || !token.Valid {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
            }
            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid claims"})
            }
            // Load user from DB
            var user model.User
            database.Instance.First(&user, claims["ID"])
            if user.ID == 0 {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
            }
            c.Set(IdentityKey, &UserResponse{ID: user.ID, Username: user.Username})
            return next(c)
        }
    }
}

