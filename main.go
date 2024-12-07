package main

import (
	"go-rest-api/controller"
	"go-rest-api/db"
	"go-rest-api/repository"
	"go-rest-api/router"
	"go-rest-api/usecase"
	"go-rest-api/validator"
)

func main() {
	// データベース接続
	db := db.NewDB()

	// バリデーション
	userValidator := validator.NewUserValidator()
	taskValidator := validator.NewTaskValidator()

	// リポジトリ
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)
	genreRepository := repository.NewGenreRepository(db) // ジャンル用リポジトリ

	// ユースケース
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	taskUsecase := usecase.NewTaskUsecase(taskRepository, taskValidator)
	genreUsecase := usecase.NewGenreUsecase(genreRepository) // ジャンル用ユースケース

	// コントローラー
	userController := controller.NewUserController(userUsecase)
	taskController := controller.NewTaskController(taskUsecase)
	genreController := controller.NewGenreController(genreUsecase) // ジャンル用コントローラー

	// ルーター
	e := router.NewRouter(userController, taskController, genreController) // 追加

	// サーバー開始
	e.Logger.Fatal(e.Start(":8080"))
}
