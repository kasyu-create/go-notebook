package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IAuthController interface {
	Me(c echo.Context) error
}

type AuthController struct{}

func NewAuthController() IAuthController {
	return &AuthController{}
}

func (ac *AuthController) Me(c echo.Context) error {
	user := c.Get("user") // JWT ミドルウェアからユーザー情報を取得
	if user == nil {
		fmt.Println("JWT token is missing or invalid")
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthenticated"})
	}
	fmt.Printf("Authenticated user: %+v\n", user)
	return c.JSON(http.StatusOK, echo.Map{"user": user})
}
