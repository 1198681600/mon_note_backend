package controller

import (
	"awesomeProject/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IEmotionController interface {
	AnalyzeDailyPattern(ctx *gin.Context)
}

type emotionController struct {
	claudeService service.IClaudeService
}

func NewEmotionController(claudeService service.IClaudeService) IEmotionController {
	return &emotionController{
		claudeService: claudeService,
	}
}

func (c *emotionController) AnalyzeDailyPattern(ctx *gin.Context) {
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