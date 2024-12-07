package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository" // リポジトリをインポート
)

type IGenreUsecase interface {
	GetAllGenres() ([]model.Genre, error)
}

type genreUsecase struct {
	repo repository.IGenreRepository // 修正: repository.IGenreRepository を使用
}

func NewGenreUsecase(repo repository.IGenreRepository) IGenreUsecase { // 修正: repository.IGenreRepository を使用
	return &genreUsecase{repo}
}

func (gu *genreUsecase) GetAllGenres() ([]model.Genre, error) {
	return gu.repo.GetAllGenres()
}
