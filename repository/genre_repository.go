package repository

import (
	"go-rest-api/model"

	"gorm.io/gorm"
)

type IGenreRepository interface {
	GetAllGenres() ([]model.Genre, error)
}

type genreRepository struct {
	db *gorm.DB
}

func NewGenreRepository(db *gorm.DB) IGenreRepository {
	return &genreRepository{db}
}

func (gr *genreRepository) GetAllGenres() ([]model.Genre, error) {
	var genres []model.Genre
	if err := gr.db.Find(&genres).Error; err != nil {
		return nil, err
	}
	return genres, nil
}
