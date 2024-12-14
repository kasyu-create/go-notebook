package usecase

import (
	"fmt"
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"
)

type ITaskUsecase interface {
	GetAllTasks(userId uint) ([]model.TaskResponse, error)
	GetTaskById(userId uint, taskId uint) (model.TaskResponse, error)
	CreateTask(task model.Task) (model.TaskResponse, error)
	UpdateTask(task model.Task, userId uint, taskId uint) (model.TaskResponse, error)
	DeleteTask(userId uint, taskId uint) error
}

type taskUsecase struct {
	tr repository.ITaskRepository
	tv validator.ITaskValidator
}

func NewTaskUsecase(tr repository.ITaskRepository, tv validator.ITaskValidator) ITaskUsecase {
	return &taskUsecase{tr, tv}
}

func (tu *taskUsecase) GetAllTasks(userId uint) ([]model.TaskResponse, error) {
	tasks := []model.Task{}
	if err := tu.tr.GetAllTasks(&tasks, userId); err != nil {
		return nil, err
	}
	resTasks := []model.TaskResponse{}
	for _, v := range tasks {
		// GenreID をポインタ型に変換
		var genreID *uint
		if v.GenreID != nil { // GenreID が nil でない場合のみ代入
			genreID = v.GenreID
		}

		// タスクレスポンスを作成し、リストに追加
		resTasks = append(resTasks, model.TaskResponse{
			ID:        v.ID,
			Title:     v.Title,
			GenreID:   genreID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}
	return resTasks, nil
}

func (tu *taskUsecase) GetTaskById(userId uint, taskId uint) (model.TaskResponse, error) {
	task := model.Task{}
	if err := tu.tr.GetTaskById(&task, userId, taskId); err != nil {
		return model.TaskResponse{}, err
	}
	resTask := model.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	return resTask, nil
}

func (tu *taskUsecase) CreateTask(task model.Task) (model.TaskResponse, error) {
	if err := tu.tv.TaskValidate(task); err != nil {
		return model.TaskResponse{}, err
	}
	if err := tu.tr.CreateTask(&task); err != nil {
		return model.TaskResponse{}, err
	}
	resTask := model.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	return resTask, nil
}

func (tu *taskUsecase) UpdateTask(task model.Task, userId uint, taskId uint) (model.TaskResponse, error) {
	// 既存データを取得
	existingTask := model.Task{}
	if err := tu.tr.GetTaskById(&existingTask, userId, taskId); err != nil {
		return model.TaskResponse{}, err
	}

	// 必要なフィールドを更新
	if task.Title != "" {
		existingTask.Title = task.Title
	}
	if task.GenreID != nil { // task.GenreID が nil でない場合に更新
		fmt.Printf("Updating GenreID: %v\n", *task.GenreID) // デバッグログ
		existingTask.GenreID = task.GenreID
	}

	// 更新をリポジトリに反映
	if err := tu.tr.UpdateTask(&existingTask, userId, taskId); err != nil {
		return model.TaskResponse{}, err
	}

	// 更新後のレスポンスを作成
	return model.TaskResponse{
		ID:        existingTask.ID,
		Title:     existingTask.Title,
		GenreID:   existingTask.GenreID,
		Order:     existingTask.Order,
		CreatedAt: existingTask.CreatedAt,
		UpdatedAt: existingTask.UpdatedAt,
	}, nil
}

func (tu *taskUsecase) DeleteTask(userId uint, taskId uint) error {
	if err := tu.tr.DeleteTask(userId, taskId); err != nil {
		return err
	}
	return nil
}
