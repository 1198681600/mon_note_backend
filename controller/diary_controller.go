package controller

import (
	"net/http"
	"strconv"
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

func (c *DiaryController) CreateDiary(ctx *gin.Context) {
	var req CreateDiaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint("userID")

	loc, _ := time.LoadLocation("Asia/Shanghai")
	dateStr := ctx.DefaultQuery("date", time.Now().In(loc).Format("2006-01-02"))
	date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	diary, err := c.service.CreateDiary(userID, req.Content, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, diary)
}

func (c *DiaryController) GetDiary(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diary ID"})
		return
	}

	userID := ctx.GetUint("userID")
	diary, err := c.service.GetDiary(uint(id), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Diary not found"})
		return
	}

	ctx.JSON(http.StatusOK, diary)
}

func (c *DiaryController) GetDiaries(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	diaries, err := c.service.GetDiaries(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, diaries)
}

type UpdateDiaryRequest struct {
	Content string `json:"content" binding:"required"`
}

func (c *DiaryController) UpdateDiary(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diary ID"})
		return
	}

	var req UpdateDiaryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint("userID")
	diary, err := c.service.UpdateDiary(uint(id), userID, req.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, diary)
}

func (c *DiaryController) DeleteDiary(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diary ID"})
		return
	}

	userID := ctx.GetUint("userID")
	if err := c.service.DeleteDiary(uint(id), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Diary deleted successfully"})
}
