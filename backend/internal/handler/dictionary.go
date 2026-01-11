package handler

import (
	"strconv"

	"dict-hub/internal/cache"
	"dict-hub/internal/service"
	"dict-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

// DictionaryHandler 字典管理处理器
type DictionaryHandler struct {
	dictSourceSvc *service.DictSourceService
	downloadSvc   *service.DownloadService
	cache         *cache.Cache
}

// NewDictionaryHandler 创建字典管理处理器
func NewDictionaryHandler(dictSourceSvc *service.DictSourceService, downloadSvc *service.DownloadService, cache *cache.Cache) *DictionaryHandler {
	return &DictionaryHandler{
		dictSourceSvc: dictSourceSvc,
		downloadSvc:   downloadSvc,
		cache:         cache,
	}
}

// List 获取字典列表
// GET /api/v1/dictionaries
func (h *DictionaryHandler) List(c *gin.Context) {
	sources, err := h.dictSourceSvc.List()
	if err != nil {
		response.InternalError(c, "failed to list dictionaries: "+err.Error())
		return
	}

	response.Success(c, sources)
}

// AddRequest 添加字典请求
type AddRequest struct {
	Path string `json:"path" binding:"required"`
}

// Add 添加字典
// POST /api/v1/dictionaries
func (h *DictionaryHandler) Add(c *gin.Context) {
	var req AddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	source, err := h.dictSourceSvc.Add(req.Path)
	if err != nil {
		switch err {
		case service.ErrDictFileNotFound:
			response.NotFound(c, "dictionary file not found")
		case service.ErrDictAlreadyExists:
			response.BadRequest(c, "dictionary already exists")
		default:
			response.InternalError(c, "failed to add dictionary: "+err.Error())
		}
		return
	}

	// 清除搜索缓存，新增字典可能影响搜索结果
	h.cache.Clear()

	response.Created(c, source)
}

// Toggle 切换字典启用状态
// PUT /api/v1/dictionaries/:id/toggle
func (h *DictionaryHandler) Toggle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid dictionary id")
		return
	}

	source, err := h.dictSourceSvc.Toggle(uint(id))
	if err != nil {
		if err == service.ErrDictSourceNotFound {
			response.NotFound(c, "dictionary not found")
			return
		}
		response.InternalError(c, "failed to toggle dictionary: "+err.Error())
		return
	}

	// 清除搜索缓存，字典启用/禁用状态变更会影响搜索结果
	h.cache.Clear()

	response.Success(c, source)
}

// Reorder 重新排序字典
// PUT /api/v1/dictionaries/reorder
func (h *DictionaryHandler) Reorder(c *gin.Context) {
	var orders []service.ReorderItem
	if err := c.ShouldBindJSON(&orders); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.dictSourceSvc.Reorder(orders); err != nil {
		response.InternalError(c, "failed to reorder dictionaries: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "reorder successful"})
}

// DownloadRequest 下载字典请求
type DownloadRequest struct {
	URL string `json:"url" binding:"required"`
}

// Download 启动异步下载
// POST /api/v1/dictionaries/download
func (h *DictionaryHandler) Download(c *gin.Context) {
	var req DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	task, err := h.downloadSvc.StartDownload(req.URL)
	if err != nil {
		switch err {
		case service.ErrInvalidURL:
			response.BadRequest(c, "invalid URL")
		case service.ErrInvalidFileFormat:
			response.BadRequest(c, "invalid file format, only .mdx files are supported")
		default:
			response.InternalError(c, "failed to start download: "+err.Error())
		}
		return
	}

	response.Created(c, task)
}

// GetDownloadStatus 获取下载任务状态
// GET /api/v1/dictionaries/download/:taskId
func (h *DictionaryHandler) GetDownloadStatus(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	task, err := h.downloadSvc.GetTaskStatus(uint(taskID))
	if err != nil {
		if err == service.ErrTaskNotFound {
			response.NotFound(c, "download task not found")
			return
		}
		response.InternalError(c, "failed to get task status: "+err.Error())
		return
	}

	response.Success(c, task)
}

// Delete 删除字典
// DELETE /api/v1/dictionaries/:id
func (h *DictionaryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid dictionary id")
		return
	}

	if err := h.dictSourceSvc.Delete(uint(id)); err != nil {
		if err == service.ErrDictSourceNotFound {
			response.NotFound(c, "dictionary not found")
			return
		}
		response.InternalError(c, "failed to delete dictionary: "+err.Error())
		return
	}

	// 清除搜索缓存，删除字典会影响搜索结果
	h.cache.Clear()

	response.Success(c, gin.H{"message": "dictionary deleted"})
}
