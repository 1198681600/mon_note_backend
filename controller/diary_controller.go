package controller

import (
	"awesomeProject/model"
	"encoding/json"
	"net/http"
	"time"

	"awesomeProject/service"
	"github.com/gin-gonic/gin"
)

type DiaryController struct {
	service service.IDiaryService
}

func NewDiaryController(service service.IDiaryService) *DiaryController {
	return &DiaryController{service: service}
}

type CreateDiaryRequest struct {
	Content string `json:"content" binding:"required"`
}

type DiaryResponse struct {
	*model.Diary
	EmotionData *service.DiaryEmotionResponse `json:"emotion_data,omitempty"`
}

func (c *DiaryController) convertToResponse(diary *model.Diary) *DiaryResponse {
	response := &DiaryResponse{Diary: diary}
	
	if diary.EmotionAnalysis != "" {
		var emotionData service.DiaryEmotionResponse
		if err := json.Unmarshal([]byte(diary.EmotionAnalysis), &emotionData); err == nil {
			response.EmotionData = &emotionData
		}
	}
	
	return response
}

func (c *DiaryController) convertToResponseList(diaries []*model.Diary) []*DiaryResponse {
	responses := make([]*DiaryResponse, len(diaries))
	for i, diary := range diaries {
		responses[i] = c.convertToResponse(diary)
	}
	return responses
}

func (c *DiaryController) CreateDiary(ctx *gin.Context) {
	var req CreateDiaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := ctx.Get("user")
	if !ok {
		return
	}

	userModel := user.(*model.User)

	loc, _ := time.LoadLocation("Asia/Shanghai")
	dateStr := ctx.DefaultQuery("date", time.Now().In(loc).Format("2006-01-02"))
	date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	diary, err := c.service.CreateDiary(userModel.ID, req.Content, date)
	if err != nil {
		if err.Error() == "diary already exists for this date" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, c.convertToResponse(diary))
}

func (c *DiaryController) GetDiary(ctx *gin.Context) {
	var req GetDiaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := ctx.Get("user")
	if !ok {
		return
	}

	userModel := user.(*model.User)
	diary, err := c.service.GetDiary(req.ID, userModel.ID)
	if err != nil {
		if err.Error() == "unauthorized" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Diary not found"})
		}
		return
	}

	ctx.JSON(http.StatusOK, c.convertToResponse(diary))
}

func (c *DiaryController) GetDiaries(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		return
	}

	userModel := user.(*model.User)
	diaries, err := c.service.GetDiaries(userModel.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, c.convertToResponseList(diaries))
}

type GetDiaryRequest struct {
	ID uint `json:"id" binding:"required"`
}

type UpdateDiaryRequest struct {
	ID      uint   `json:"id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type DeleteDiaryRequest struct {
	ID uint `json:"id" binding:"required"`
}

func (c *DiaryController) UpdateDiary(ctx *gin.Context) {
	var req UpdateDiaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := ctx.Get("user")
	if !ok {
		return
	}

	userModel := user.(*model.User)
	diary, err := c.service.UpdateDiary(req.ID, userModel.ID, req.Content)
	if err != nil {
		if err.Error() == "unauthorized" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, c.convertToResponse(diary))
}

func (c *DiaryController) DeleteDiary(ctx *gin.Context) {
	var req DeleteDiaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := ctx.Get("user")
	if !ok {
		return
	}

	userModel := user.(*model.User)
	if err := c.service.DeleteDiary(req.ID, userModel.ID); err != nil {
		if err.Error() == "unauthorized" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Diary deleted successfully"})
}
