package main

import (
	"awesomeProject/config"
	"awesomeProject/controller"
	"awesomeProject/middleware"
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
	userStorage := storage.NewUserStorage(db)
	authService := service.NewAuthService(userStorage)
	authController := controller.NewAuthController(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

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
			protected.POST("/logout", authController.Logout)
		}
	}

	log.Println("MoodTrace 后端服务启动在 :8080 端口")
	log.Fatal(router.Run(":8080"))
}