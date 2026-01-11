package handler

import (
	"net/http"
	"strconv"
	"time"

	"dict-hub/internal/service"
	"dict-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

// HistoryHandler 历史记录处理器
type HistoryHandler struct {
	historyService *service.HistoryService
}

// NewHistoryHandler 创建历史记录处理器
func NewHistoryHandler(historyService *service.HistoryService) *HistoryHandler {
	return &HistoryHandler{historyService: historyService}
}

// List 获取历史记录列表
// GET /api/v1/history?page=1&page_size=20
func (h *HistoryHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	histories, total, err := h.historyService.List(page, pageSize)
	if err != nil {
		response.InternalError(c, "failed to get history: "+err.Error())
		return
	}

	response.PagedSuccess(c, histories, total, page, pageSize)
}

// Clear 清空历史记录
// DELETE /api/v1/history
func (h *HistoryHandler) Clear(c *gin.Context) {
	if err := h.historyService.Clear(); err != nil {
		response.InternalError(c, "failed to clear history: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "history cleared"})
}

// Export 导出历史记录为 CSV
// GET /api/v1/history/export
func (h *HistoryHandler) Export(c *gin.Context) {
	data, err := h.historyService.ExportCSV()
	if err != nil {
		response.InternalError(c, "failed to export history: "+err.Error())
		return
	}

	fileName := "search_history_" + time.Now().Format("20060102_150405") + ".csv"
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, "text/csv; charset=utf-8", data)
}
