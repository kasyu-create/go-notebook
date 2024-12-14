package repository

import (
	"fmt"
	"go-rest-api/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ITaskRepository interface {
	GetAllTasks(tasks *[]model.Task, userId uint) error
	GetTaskById(task *model.Task, userId uint, taskId uint) error
	CreateTask(task *model.Task) error
	UpdateTask(task *model.Task, userId uint, taskId uint) error
	DeleteTask(userId uint, taskId uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &taskRepository{db}
}

func (tr *taskRepository) GetAllTasks(tasks *[]model.Task, userId uint) error {
	if err := tr.db.
		Preload("Genre").
		Where("user_id = ?", userId).
		Order("created_at").
		Find(tasks).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) GetTaskById(task *model.Task, userId uint, taskId uint) error {
	if err := tr.db.Joins("User").Where("user_id=?", userId).First(task, taskId).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) CreateTask(task *model.Task) error {
	if err := tr.db.Create(task).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) UpdateTask(task *model.Task, userId uint, taskId uint) error {
	// 更新処理
	result := tr.db.Model(task).
		Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
			{Name: "title"},
			{Name: "genre_id"}, // ジャンルも返却に含める
			{Name: "updated_at"},
		}}).
		Where("id=? AND user_id=?", taskId, userId).
		Updates(map[string]interface{}{
			"title":    task.Title,
			"genre_id": task.GenreID, // ジャンルを更新
		})

	// エラー処理
	if result.Error != nil {
		fmt.Printf("UpdateTask Error: %v\n", result.Error)
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}

func (tr *taskRepository) DeleteTask(userId uint, taskId uint) error {
	result := tr.db.Where("id=? AND user_id=?", taskId, userId).Delete(&model.Task{})
	if result.Error != nil {
		fmt.Printf("DeleteTask Error: %v\n", result.Error)
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}
