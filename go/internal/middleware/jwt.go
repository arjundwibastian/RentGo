package middleware

import (
	"milestone-02/internal/dto"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type JWTMiddleware struct {
	jwtSecret []byte
}

// Constructor
func NewJWTMiddleware(jwtSecret []byte) *JWTMiddleware {
	return &JWTMiddleware{jwtSecret: jwtSecret}
}

func (m *JWTMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: Token not found",
				ResponseData:    nil,
			})
		}

		tokenString := ""
		splitString := strings.Split(authHeader, " ")
		if len(splitString) > 1 && splitString[0] == "Bearer" {
			tokenString = splitString[1]
		} else {
			return c.JSON(http.StatusUnauthorized, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: Token not found",
				ResponseData:    nil,
			})
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return m.jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: Token is invalid",
				ResponseData:    nil,
			})
		}

		claims := token.Claims.(jwt.MapClaims)
		userIDFloat,ok := claims["user_id"].(float64)
		if !ok {
			return c.JSON(http.StatusUnauthorized, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: Token is invalid",
				ResponseData:    nil,
			})
		}
		userID := int(userIDFloat)

		role, ok := claims["role"].(string)
		if !ok {
				return c.JSON(http.StatusUnauthorized, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: role is invalid",
				ResponseData:    nil,
			})
		}
		c.Set("user_id", userID)
		c.Set("role", role)

		return next(c)
	}
}


func (m *JWTMiddleware) AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("role").(string)
		if !ok || role != "admin" {
			return c.JSON(http.StatusForbidden, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: Admin access required",
			})
		}
		return next(c)
	}
}

func (m *JWTMiddleware) UserOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("role").(string)
		if !ok || role != "user" {
			return c.JSON(http.StatusForbidden, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: User access required",
			})
		}
		return next(c)
	}
}