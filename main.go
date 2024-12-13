package main

import (
	"go-rest-api/controller"
	"go-rest-api/db"
	"go-rest-api/repository"
	"go-rest-api/router"
	"go-rest-api/usecase"
	"go-rest-api/validator"
	"log"
	"os"
)

func main() {
	// ログファイルを開く
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// ログ出力先を設定
	log.SetOutput(logFile)

	// ログフォーマットを設定（例: プリフィックス）
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 初期化ログ
	log.Println("Starting the application...")

	// データベース接続
	db := db.NewDB()
	log.Println("Database connected successfully")

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
	log.Println("Starting the server on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
