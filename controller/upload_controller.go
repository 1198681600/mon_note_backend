package controller

import (
	"awesomeProject/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IUploadController interface {
	GenerateUploadURL(ctx *gin.Context)
}

type uploadController struct {
	uploadService service.IUploadService
}

func NewUploadController(uploadService service.IUploadService) IUploadController {
	return &uploadController{
		uploadService: uploadService,
	}
}

type GenerateUploadURLRequest struct {
	FileType string `json:"file_type" binding:"required"`
}

func (c *uploadController) GenerateUploadURL(ctx *gin.Context) {
	var req GenerateUploadURLRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	// 验证文件类型
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	if !allowedTypes[req.FileType] {
		ctx.JSON(http.StatusBadRequest, ApiResponse{
			Code:    400,
			Message: "不支持的文件类型，仅支持 JPEG、PNG、WebP 格式",
		})
		return
	}

	result, err := c.uploadService.GeneratePresignedURL(req.FileType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Code:    200,
		Message: "生成上传链接成功",
		Data:    result,
	})
}