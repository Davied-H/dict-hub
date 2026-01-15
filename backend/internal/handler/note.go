package handler

import (
	"errors"
	"strconv"

	"dict-hub/internal/model"
	"dict-hub/internal/service/vocabulary"
	"dict-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

// NoteHandler 笔记处理器
type NoteHandler struct {
	noteService *vocabulary.NoteService
}

// NewNoteHandler 创建笔记处理器
func NewNoteHandler(noteService *vocabulary.NoteService) *NoteHandler {
	return &NoteHandler{noteService: noteService}
}

// CreateNoteRequest 创建笔记请求
type CreateNoteRequest struct {
	Content string `json:"content" binding:"required"`
}

// Create 创建笔记
// POST /api/v1/vocabulary/:id/notes
func (h *NoteHandler) Create(c *gin.Context) {
	vocabularyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid vocabulary id")
		return
	}

	var req CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	note := &model.Note{
		VocabularyID: uint(vocabularyID),
		Content:      req.Content,
	}

	if err := h.noteService.Create(note); err != nil {
		if errors.Is(err, vocabulary.ErrVocabularyNotFound) {
			response.NotFound(c, "vocabulary not found")
			return
		}
		response.InternalError(c, "failed to create note: "+err.Error())
		return
	}

	response.Created(c, note)
}

// UpdateNoteRequest 更新笔记请求
type UpdateNoteRequest struct {
	Content string `json:"content" binding:"required"`
}

// Update 更新笔记
// PUT /api/v1/notes/:id
func (h *NoteHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid note id")
		return
	}

	var req UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if err := h.noteService.Update(uint(id), req.Content); err != nil {
		if errors.Is(err, vocabulary.ErrNoteNotFound) {
			response.NotFound(c, "note not found")
			return
		}
		response.InternalError(c, "failed to update note: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "note updated"})
}

// Delete 删除笔记
// DELETE /api/v1/notes/:id
func (h *NoteHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid note id")
		return
	}

	if err := h.noteService.Delete(uint(id)); err != nil {
		if errors.Is(err, vocabulary.ErrNoteNotFound) {
			response.NotFound(c, "note not found")
			return
		}
		response.InternalError(c, "failed to delete note: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "note deleted"})
}

// ListByVocabulary 获取某单词的笔记列表
// GET /api/v1/vocabulary/:id/notes
func (h *NoteHandler) ListByVocabulary(c *gin.Context) {
	vocabularyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid vocabulary id")
		return
	}

	notes, err := h.noteService.ListByVocabulary(uint(vocabularyID))
	if err != nil {
		response.InternalError(c, "failed to list notes: "+err.Error())
		return
	}

	response.Success(c, notes)
}
