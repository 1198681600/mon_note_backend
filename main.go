package main

import (
	"awesomeProject/config"
	"awesomeProject/controller"
	"awesomeProject/middleware"
	"awesomeProject/model"
	"awesomeProject/service"
	"awesomeProject/storage"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.InitDatabase(); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	db := config.GetDB()

	// 自动迁移
	if err := db.AutoMigrate(&model.User{}, &model.Diary{}); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 依赖注入
	var userStorage = storage.NewUserStorage(db)
	var diaryStorage = storage.NewDiaryStorage(db)
	var authService = service.NewAuthService(userStorage)
	var claudeService = service.NewClaudeService()
	var uploadService = service.NewUploadService()
	var diaryService = service.NewDiaryService(diaryStorage, claudeService)
	var authController = controller.NewAuthController(authService)
	var emotionController = controller.NewEmotionController(claudeService)
	var uploadController = controller.NewUploadController(uploadService)
	var diaryController = controller.NewDiaryController(diaryService)
	var authMiddleware = middleware.NewAuthMiddleware(authService)

	router := gin.Default()

	router.Use(authMiddleware.CORS())

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/send-code", authController.SendVerificationCode)
			auth.POST("/login", authController.Login)
		}

		protected := api.Group("/")
		protected.Use(authMiddleware.RequireAuth())
		{
			protected.GET("/profile", authController.GetProfile)
			protected.POST("/profile", authController.UpdateProfile)
			protected.POST("/upload/generate-url", uploadController.GenerateUploadURL)
			protected.POST("/emotion/analyze-diary", emotionController.AnalyzeDiary)
			protected.POST("/emotion/analyze-weekly", emotionController.AnalyzeWeekly)

			diary := protected.Group("/diary")
			{
				diary.POST("/create", diaryController.CreateDiary)
				diary.POST("/list", diaryController.GetDiaries)
				diary.POST("/get", diaryController.GetDiary)
				diary.POST("/update", diaryController.UpdateDiary)
				diary.POST("/delete", diaryController.DeleteDiary)
			}
		}
	}

	log.Println("MoodTrace 后端服务启动在 :8080 端口")
	log.Fatal(router.Run(":8080"))
}
