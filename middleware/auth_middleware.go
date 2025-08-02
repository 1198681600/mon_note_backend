package middleware

import (
	"awesomeProject/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type IAuthMiddleware interface {
	RequireAuth() gin.HandlerFunc
	CORS() gin.HandlerFunc
}

type authMiddleware struct {
	authService service.IAuthService
}

func NewAuthMiddleware(authService service.IAuthService) IAuthMiddleware {
	return &authMiddleware{
		authService: authService,
	}
}

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (m *authMiddleware) RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, ApiResponse{
				Code:    401,
				Message: "未提供认证token",
			})
			ctx.Abort()
			return
		}

		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = authHeader[7:]
		}

		user, err := m.authService.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, ApiResponse{
				Code:    401,
				Message: err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Set("token", token)
		ctx.Next()
	}
}

func (m *authMiddleware) CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}