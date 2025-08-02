package controller

import (
	"awesomeProject/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmotionController struct {
	claudeService *service.ClaudeService
}

func NewEmotionController(claudeService *service.ClaudeService) *EmotionController {
	return &EmotionController{
		claudeService: claudeService,
	}
}

func (c *EmotionController) AnalyzeDailyPattern(ctx *gin.Context) {
	var req service.EmotionAnalysisRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	if len(req.UserData) == 0 {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: "用户数据不能为空",
		})
		return
	}

	result, err := c.claudeService.AnalyzeDailyEmotionPattern(req.UserData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ApiResponse{
			Code:    500,
			Message: "情绪分析失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Code:    200,
		Message: "情绪分析完成",
		Data:    result,
	})
}