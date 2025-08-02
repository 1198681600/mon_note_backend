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
	var userStorage storage.IUserStorage = storage.NewUserStorage(db)
	var diaryStorage storage.IDiaryStorage = storage.NewDiaryStorage(db)
	var authService service.IAuthService = service.NewAuthService(userStorage)
	var claudeService service.IClaudeService = service.NewClaudeService()
	var uploadService service.IUploadService = service.NewUploadService()
	var diaryService service.IDiaryService = service.NewDiaryService(diaryStorage)
	var authController controller.IAuthController = controller.NewAuthController(authService)
	var emotionController controller.IEmotionController = controller.NewEmotionController(claudeService)
	var uploadController controller.IUploadController = controller.NewUploadController(uploadService)
	var diaryController = controller.NewDiaryController(diaryService)
	var authMiddleware middleware.IAuthMiddleware = middleware.NewAuthMiddleware(authService)

	router := gin.Default()

	router.Use(authMiddleware.CORS())

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/send-code", authController.SendVerificationCode)
			auth.POST("/register", authController.Register)
			auth.POST("/verify-email", authController.VerifyEmail)
			auth.POST("/login", authController.Login)
		}

		protected := api.Group("/")
		protected.Use(authMiddleware.RequireAuth())
		{
			protected.GET("/profile", authController.GetProfile)
			protected.POST("/profile", authController.UpdateProfile)
			protected.POST("/upload/generate-url", uploadController.GenerateUploadURL)
			protected.POST("/emotion/analyze-daily", emotionController.AnalyzeDailyPattern)

			diary := protected.Group("/diaries")
			{
				diary.POST("", diaryController.CreateDiary)
				diary.GET("", diaryController.GetDiaries)
				diary.GET("/:id", diaryController.GetDiary)
				diary.PUT("/:id", diaryController.UpdateDiary)
				diary.DELETE("/:id", diaryController.DeleteDiary)
			}
		}
	}

	log.Println("MoodTrace 后端服务启动在 :8080 端口")
	log.Fatal(router.Run(":8080"))
}
