package model

import "time"

type Task struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	GenreID   *uint     `json:"genre_id" gorm:""`                // 外部キー
	Genre     Genre     `json:"genre" gorm:"foreignKey:GenreID"` // リレーション
	Order     *int      `json:"order" gorm:""`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user" gorm:"foreignKey:UserId; constraint:OnDelete:CASCADE"`
	UserId    uint      `json:"user_id" gorm:"not null"`
}

type TaskResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	GenreID   *uint     `json:"genre_id"`
	GenreName string    `json:"genre_name"` // ジャンル名をレスポンスに含める
	Order     *int      `json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
