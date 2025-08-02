package controller

import (
	"awesomeProject/model"
	"awesomeProject/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	SendVerificationCode(ctx *gin.Context)
	Register(ctx *gin.Context)
	VerifyEmail(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetProfile(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)
}

type authController struct {
	authService service.IAuthService
}

func NewAuthController(authService service.IAuthService) IAuthController {
	return &authController{
		authService: authService,
	}
}

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (c *authController) SendVerificationCode(ctx *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	if err := c.authService.SendVerificationCode(req.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Code:    200,
		Message: "验证码已发送",
	})
}

func (c *authController) Register(ctx *gin.Context) {
	var req service.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	if err := c.authService.Register(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Code:    200,
		Message: "注册成功，请验证邮箱",
	})
}

func (c *authController) VerifyEmail(ctx *gin.Context) {
	var req service.VerifyEmailRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	if err := c.authService.VerifyEmail(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Code:    200,
		Message: "邮箱验证成功",
	})
}

func (c *authController) Login(ctx *gin.Context) {
	var req service.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	resp, err := c.authService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Code:    200,
		Message: "登录成功",
		Data:    resp,
	})
}

func (c *authController) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, ApiResponse{
			Code:    401,
			Message: "未提供token",
		})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := c.authService.Logout(token); err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Code:    200,
		Message: "退出登录成功",
	})
}

func (c *authController) GetProfile(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, ApiResponse{
			Code:    401,
			Message: "用户未登录",
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Code:    200,
		Message: "获取用户信息成功",
		Data:    user,
	})
}

func (c *authController) UpdateProfile(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, ApiResponse{
			Code:    401,
			Message: "用户未登录",
		})
		return
	}

	var req service.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	userModel := user.(*model.User)
	updatedUser, err := c.authService.UpdateProfile(userModel.ID, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Code:    200,
		Message: "更新用户信息成功",
		Data:    updatedUser,
	})
}
