package middlewares

import (
	. "booking-api/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

type Claims struct {
	Id string
	jwt.RegisteredClaims
}

func VerifyJwt(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "No token provided")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &UserJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		claims, ok := token.Claims.(*UserJwtClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}
		c.Set("user", claims)
		return next(c)
	}
}
