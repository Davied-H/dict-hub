package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"dict-hub/internal/model"
	"dict-hub/internal/service/vocabulary"
	"dict-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

// VocabularyHandler 生词本处理器
type VocabularyHandler struct {
	vocabService *vocabulary.VocabularyService
	noteService  *vocabulary.NoteService
}

// NewVocabularyHandler 创建生词本处理器
func NewVocabularyHandler(vocabService *vocabulary.VocabularyService, noteService *vocabulary.NoteService) *VocabularyHandler {
	return &VocabularyHandler{
		vocabService: vocabService,
		noteService:  noteService,
	}
}

// VocabularyAddRequest 添加生词请求
type VocabularyAddRequest struct {
	Word       string `json:"word" binding:"required"`
	DictID     uint   `json:"dict_id"`
	DictTitle  string `json:"dict_title"`
	Definition string `json:"definition"`
	Phonetic   string `json:"phonetic"`
	Tags       string `json:"tags"`
}

// Add 添加单词到生词本
// POST /api/v1/vocabulary
func (h *VocabularyHandler) Add(c *gin.Context) {
	var req VocabularyAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	vocab := &model.Vocabulary{
		Word:       req.Word,
		DictID:     req.DictID,
		DictTitle:  req.DictTitle,
		Definition: req.Definition,
		Phonetic:   req.Phonetic,
		Tags:       req.Tags,
	}

	if err := h.vocabService.Add(vocab); err != nil {
		if errors.Is(err, vocabulary.ErrWordAlreadyExists) {
			response.BadRequest(c, "word already exists in vocabulary")
			return
		}
		response.InternalError(c, "failed to add vocabulary: "+err.Error())
		return
	}

	response.Created(c, vocab)
}

// Remove 从生词本移除单词
// DELETE /api/v1/vocabulary/:id
func (h *VocabularyHandler) Remove(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.vocabService.Remove(uint(id)); err != nil {
		if errors.Is(err, vocabulary.ErrVocabularyNotFound) {
			response.NotFound(c, "vocabulary not found")
			return
		}
		response.InternalError(c, "failed to remove vocabulary: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "vocabulary removed"})
}

// Get 获取单个生词详情
// GET /api/v1/vocabulary/:id
func (h *VocabularyHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	vocab, err := h.vocabService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, vocabulary.ErrVocabularyNotFound) {
			response.NotFound(c, "vocabulary not found")
			return
		}
		response.InternalError(c, "failed to get vocabulary: "+err.Error())
		return
	}

	// 获取笔记数量
	noteCount, _ := h.noteService.CountByVocabulary(uint(id))

	response.Success(c, gin.H{
		"vocabulary": vocab,
		"note_count": noteCount,
	})
}

// List 获取生词本列表
// GET /api/v1/vocabulary?page=1&page_size=20&sort_by=created_at&sort_order=desc&level=0&tag=xxx&keyword=xxx
func (h *VocabularyHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	params := vocabulary.ListParams{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		Tag:       c.Query("tag"),
		Keyword:   c.Query("keyword"),
	}

	if levelStr := c.Query("level"); levelStr != "" {
		level, err := strconv.Atoi(levelStr)
		if err == nil {
			params.Level = &level
		}
	}

	vocabs, total, err := h.vocabService.List(params)
	if err != nil {
		response.InternalError(c, "failed to list vocabulary: "+err.Error())
		return
	}

	response.PagedSuccess(c, vocabs, total, page, pageSize)
}

// CheckWord 检查单词是否已收藏
// GET /api/v1/vocabulary/check/:word
func (h *VocabularyHandler) CheckWord(c *gin.Context) {
	word := c.Param("word")
	if word == "" {
		response.BadRequest(c, "word is required")
		return
	}

	exists, id := h.vocabService.CheckWord(word)
	response.Success(c, gin.H{
		"exists": exists,
		"id":     id,
	})
}

// UpdateTagsRequest 更新标签请求
type UpdateTagsRequest struct {
	Tags string `json:"tags"`
}

// UpdateTags 更新标签
// PUT /api/v1/vocabulary/:id/tags
func (h *VocabularyHandler) UpdateTags(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req UpdateTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.vocabService.UpdateTags(uint(id), req.Tags); err != nil {
		if errors.Is(err, vocabulary.ErrVocabularyNotFound) {
			response.NotFound(c, "vocabulary not found")
			return
		}
		response.InternalError(c, "failed to update tags: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "tags updated"})
}

// Export 导出生词本为 CSV
// GET /api/v1/vocabulary/export
func (h *VocabularyHandler) Export(c *gin.Context) {
	data, err := h.vocabService.ExportCSV()
	if err != nil {
		response.InternalError(c, "failed to export vocabulary: "+err.Error())
		return
	}

	fileName := "vocabulary_" + time.Now().Format("20060102_150405") + ".csv"
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, "text/csv; charset=utf-8", data)
}

// Stats 获取生词本统计
// GET /api/v1/vocabulary/stats
func (h *VocabularyHandler) Stats(c *gin.Context) {
	stats, err := h.vocabService.GetStats()
	if err != nil {
		response.InternalError(c, "failed to get stats: "+err.Error())
		return
	}

	response.Success(c, stats)
}
