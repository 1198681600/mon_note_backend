package controller

import (
	"awesomeProject/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IEmotionController interface {
	AnalyzeDiary(ctx *gin.Context)
	AnalyzeWeekly(ctx *gin.Context)
}


type emotionController struct {
	claudeService service.IClaudeService
}

func NewEmotionController(claudeService service.IClaudeService) IEmotionController {
	return &emotionController{
		claudeService: claudeService,
	}
}

type AnalyzeDiaryRequest struct {
	DiaryContent string                 `json:"diary_content" binding:"required"`
	DiaryDate    string                 `json:"diary_date" binding:"required"`
	UserContext  map[string]interface{} `json:"user_context,omitempty"`
}

type AnalyzeWeeklyRequest struct {
	WeekStart string                   `json:"week_start" binding:"required"`
	DiaryData []map[string]interface{} `json:"diary_data" binding:"required"`
}

func (c *emotionController) AnalyzeDiary(ctx *gin.Context) {
	var req AnalyzeDiaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	result, err := c.claudeService.AnalyzeDiaryEmotion(req.DiaryContent, req.DiaryDate, req.UserContext)
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

func (c *emotionController) AnalyzeWeekly(ctx *gin.Context) {
	var req AnalyzeWeeklyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	result, err := c.claudeService.AnalyzeWeeklyEmotion(req.WeekStart, req.DiaryData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ApiResponse{
			Code:    500,
			Message: "情绪分析失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Code:    200,
		Message: "一周情绪分析完成",
		Data:    result,
	})
}
