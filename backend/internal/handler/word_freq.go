package handler

import (
	"dict-hub/internal/service"
	"dict-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

// WordFreqHandler 词频处理器
type WordFreqHandler struct {
	wordFreqService *service.WordFreqService
}

// NewWordFreqHandler 创建词频处理器
func NewWordFreqHandler(wordFreqService *service.WordFreqService) *WordFreqHandler {
	return &WordFreqHandler{wordFreqService: wordFreqService}
}

// Import 导入词频文件
// POST /api/v1/wordfreq/import
func (h *WordFreqHandler) Import(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}

	// 打开上传的文件
	f, err := file.Open()
	if err != nil {
		response.InternalError(c, "failed to open file: "+err.Error())
		return
	}
	defer f.Close()

	// 导入词频
	result, err := h.wordFreqService.Import(f)
	if err != nil {
		response.InternalError(c, "failed to import: "+err.Error())
		return
	}

	response.Success(c, result)
}
