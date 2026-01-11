package handler

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"dict-hub/internal/service"
	"dict-hub/internal/service/mdx"
	"dict-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

type MdxHandler struct {
	manager        mdx.DictManager
	historyService *service.HistoryService
}

func NewMdxHandler(manager mdx.DictManager) *MdxHandler {
	return &MdxHandler{manager: manager}
}

// NewMdxHandlerWithHistory 创建带历史服务的 MdxHandler
func NewMdxHandlerWithHistory(manager mdx.DictManager, historyService *service.HistoryService) *MdxHandler {
	return &MdxHandler{
		manager:        manager,
		historyService: historyService,
	}
}

// List 列出已加载的字典
// GET /api/v1/dicts
func (h *MdxHandler) List(c *gin.Context) {
	dicts := h.manager.ListLoaded()
	response.Success(c, dicts)
}

// LoadRequest 加载字典请求
type LoadRequest struct {
	Path string `json:"path" binding:"required"`
}

// Load 加载字典
// POST /api/v1/dicts/load
func (h *MdxHandler) Load(c *gin.Context) {
	var req LoadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	id, err := h.manager.LoadDict(req.Path)
	if err != nil {
		response.InternalError(c, "failed to load dictionary: "+err.Error())
		return
	}

	response.Created(c, gin.H{
		"id":      id,
		"message": "dictionary loaded successfully",
	})
}

// LoadAllRequest 批量加载请求
type LoadAllRequest struct {
	Dir string `json:"dir" binding:"required"`
}

// LoadAll 加载目录下所有字典
// POST /api/v1/dicts/load-all
func (h *MdxHandler) LoadAll(c *gin.Context) {
	var req LoadAllRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	err := h.manager.LoadAll(req.Dir)
	if err != nil {
		response.InternalError(c, "failed to load dictionaries: "+err.Error())
		return
	}

	dicts := h.manager.ListLoaded()
	response.Success(c, gin.H{
		"message": "dictionaries loaded",
		"count":   len(dicts),
		"dicts":   dicts,
	})
}

// Lookup 查询单词
// GET /api/v1/dicts/:id/lookup?word=xxx
func (h *MdxHandler) Lookup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid dictionary id")
		return
	}

	word := c.Query("word")
	if word == "" {
		response.BadRequest(c, "word parameter is required")
		return
	}

	result, err := h.manager.Lookup(uint(id), word)
	if err != nil {
		if err == mdx.ErrDictNotFound {
			response.NotFound(c, "dictionary not found")
			return
		}
		if err == mdx.ErrWordNotFound {
			response.NotFound(c, "word not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"word":       word,
		"definition": string(result),
	})
}

// Search 跨字典搜索
// GET /api/v1/dicts/search?word=xxx&ids=1,2,3
func (h *MdxHandler) Search(c *gin.Context) {
	start := time.Now()

	word := c.Query("word")
	if word == "" {
		response.BadRequest(c, "word parameter is required")
		return
	}

	var dictIDs []uint
	idsStr := c.Query("ids")
	if idsStr != "" {
		for _, idStr := range strings.Split(idsStr, ",") {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32)
			if err != nil {
				continue
			}
			dictIDs = append(dictIDs, uint(id))
		}
	}

	results := h.manager.Search(word, dictIDs...)

	// 记录搜索历史
	if h.historyService != nil {
		go h.historyService.Record(
			word,
			c.GetHeader("X-Session-ID"),
			c.ClientIP(),
			len(results) > 0,
			time.Since(start).Milliseconds(),
		)
	}

	response.Success(c, gin.H{
		"word":    word,
		"count":   len(results),
		"results": results,
	})
}

// GetResource 获取 MDD 资源
// GET /api/v1/dicts/:id/resource/*path
// GET /api/v1/resources/:dictId/*path
func (h *MdxHandler) GetResource(c *gin.Context) {
	// 支持两种路由参数名：id 和 dictId
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Param("dictId")
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid dictionary id")
		return
	}

	resourcePath := c.Param("path")
	if resourcePath == "" {
		response.BadRequest(c, "resource path is required")
		return
	}

	// 生成 ETag
	etag := fmt.Sprintf(`"%d-%s"`, id, resourcePath)

	// 检查条件请求（If-None-Match）
	if match := c.GetHeader("If-None-Match"); match == etag {
		c.Status(http.StatusNotModified)
		return
	}

	reader, err := h.manager.GetResource(uint(id), resourcePath)
	if err != nil {
		if err == mdx.ErrDictNotFound {
			response.NotFound(c, "dictionary not found")
			return
		}
		if err == mdx.ErrNoMDD {
			response.NotFound(c, "no MDD resource file")
			return
		}
		if err == mdx.ErrResourceNotFound {
			response.NotFound(c, "resource not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	// 根据文件扩展名设置 Content-Type
	ext := filepath.Ext(resourcePath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 添加 HTTP 缓存头
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400")
	c.Header("ETag", etag)
	c.Status(http.StatusOK)
	io.Copy(c.Writer, reader)
}

// Unload 卸载字典
// DELETE /api/v1/dicts/:id
func (h *MdxHandler) Unload(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid dictionary id")
		return
	}

	err = h.manager.Unload(uint(id))
	if err != nil {
		if err == mdx.ErrDictNotFound {
			response.NotFound(c, "dictionary not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "dictionary unloaded",
	})
}
