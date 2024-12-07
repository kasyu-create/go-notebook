package controller

import (
	"go-rest-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IGenreController interface {
	GetAllGenres(c echo.Context) error
}

type genreController struct {
	gu usecase.IGenreUsecase
}

func NewGenreController(gu usecase.IGenreUsecase) IGenreController {
	return &genreController{gu}
}

func (gc *genreController) GetAllGenres(c echo.Context) error {
	genres, err := gc.gu.GetAllGenres()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, genres)
}
