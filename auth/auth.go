package auth

import (
	"cinetodoapi/database"
	"cinetodoapi/model"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

const (
	IdentityKey     = "user"
	tokenCookieName = "auth_token"
)

// GetUserFromContext returns the authenticated user from context.
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

// RenderLogin renders the login page with optional data.
func RenderLogin(c echo.Context, status int, data map[string]any) error {
	if data == nil {
		data = map[string]any{}
	}
	if _, ok := data["Title"]; !ok {
		data["Title"] = "Sign in"
	}
	return c.Render(status, "auth/login.html", data)
}

// ShowLogin displays the login screen.
func ShowLogin(c echo.Context) error {
	return RenderLogin(c, http.StatusOK, nil)
}

// Login authenticates the user and issues a JWT that is stored in a cookie.
func Login(c echo.Context) error {
	var loginVals model.Login
	if err := c.Bind(&loginVals); err != nil {
		return RenderLogin(c, http.StatusBadRequest, map[string]interface{}{
			"Error": "Please provide your username and password.",
		})
	}

	loginVals.Username = strings.TrimSpace(loginVals.Username)
	var user model.User
	database.Instance.Where("username = ?", loginVals.Username).First(&user)
	if user.ID == 0 || user.CheckPassword(loginVals.Password) != nil {
		return RenderLogin(c, http.StatusUnauthorized, map[string]interface{}{
			"Error":    "Invalid username or password.",
			"Username": loginVals.Username,
		})
	}

	signed, err := issueToken(user.ID)
	if err != nil {
		return RenderLogin(c, http.StatusInternalServerError, map[string]interface{}{
			"Error": "We could not sign you in right now. Please try again.",
		})
	}

	setAuthCookie(c, signed, time.Hour)

	if isHTMXRequest(c) {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusSeeOther)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

// Logout clears the auth cookie and redirects to the login page.
func Logout(c echo.Context) error {
	clearAuthCookie(c)
	if isHTMXRequest(c) {
		c.Response().Header().Set("HX-Redirect", "/login")
		return c.NoContent(http.StatusSeeOther)
	}
	return c.Redirect(http.StatusSeeOther, "/login")
}

// JWTMiddleware validates JWT and loads user into context.
func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenStr := extractToken(c)
			if tokenStr == "" {
				return redirectToLogin(c)
			}

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid signing method")
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil || !token.Valid {
				clearAuthCookie(c)
				return redirectToLogin(c)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				clearAuthCookie(c)
				return redirectToLogin(c)
			}

			userID, err := claimToUint(claims["ID"])
			if err != nil {
				clearAuthCookie(c)
				return redirectToLogin(c)
			}

			var user model.User
			if dbErr := database.Instance.First(&user, userID).Error; dbErr != nil || user.ID == 0 {
				clearAuthCookie(c)
				return redirectToLogin(c)
			}

			c.Set(IdentityKey, &UserResponse{ID: user.ID, Username: user.Username})
			return next(c)
		}
	}
}

func issueToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"ID":  userID,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func extractToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if len(authHeader) > 7 && strings.EqualFold(authHeader[:7], "Bearer ") {
		return strings.TrimSpace(authHeader[7:])
	}
	if cookie, err := c.Cookie(tokenCookieName); err == nil {
		return cookie.Value
	}
	return ""
}

func claimToUint(value interface{}) (uint, error) {
	switch v := value.(type) {
	case float64:
		return uint(v), nil
	case float32:
		return uint(v), nil
	case int64:
		return uint(v), nil
	case int32:
		return uint(v), nil
	case int:
		return uint(v), nil
	case uint64:
		return uint(v), nil
	case uint32:
		return uint(v), nil
	case uint:
		return v, nil
	default:
		return 0, errors.New("invalid claim")
	}
}

func redirectToLogin(c echo.Context) error {
	if isHTMXRequest(c) {
		c.Response().Header().Set("HX-Redirect", "/login")
		return c.NoContent(http.StatusUnauthorized)
	}
	return c.Redirect(http.StatusSeeOther, "/login")
}

func setAuthCookie(c echo.Context, token string, ttl time.Duration) {
	secure := strings.EqualFold(os.Getenv("COOKIE_SECURE"), "true")
	cookie := &http.Cookie{
		Name:     tokenCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(ttl),
	}
	c.SetCookie(cookie)
}

func clearAuthCookie(c echo.Context) {
	cookie := &http.Cookie{
		Name:     tokenCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	}
	c.SetCookie(cookie)
}

func isHTMXRequest(c echo.Context) bool {
	return strings.EqualFold(c.Request().Header.Get("HX-Request"), "true")
}
