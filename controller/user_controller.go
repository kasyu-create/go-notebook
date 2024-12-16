package controller

import (
	"fmt"
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type IUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	CsrfToken(c echo.Context) error
}

type userController struct {
	uu usecase.IUserUsecase
}

func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

func (uc *userController) SignUp(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, userRes)
}

func (uc *userController) LogIn(c echo.Context) error {
	// リクエストボディからユーザー情報を取得
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid request"})
	}

	// ユーザーをデータベースから取得
	user, err := uc.uu.GetUserByEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid email or password"})
	}

	// パスワードを検証
	if !uc.uu.ValidatePassword(user.Password, req.Password) {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid email or password"})
	}

	// JWT ペイロードを作成
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // トークン有効期限: 24時間
	}

	// トークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET")))
	fmt.Printf("Set-Cookie: %v\n", signedToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "could not generate token"})
	}

	// トークンをクッキーにセット
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    signedToken,
		HttpOnly: true,
		Path:     "/",
		// Domain:   os.Getenv("API_DOMAIN"),
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // ローカル環境ではfalse、本番環境ではtrue
	})

	return c.JSON(http.StatusOK, echo.Map{"message": "logged in successfully"})
}

func (uc *userController) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}

func (uc *userController) CsrfToken(c echo.Context) error {
	token := c.Get("csrf").(string)
	return c.JSON(http.StatusOK, echo.Map{
		"csrf_token": token,
	})
}
