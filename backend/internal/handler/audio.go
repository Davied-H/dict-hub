package handler

import (
	"io"
	"net/http"

	"dict-hub/internal/service/audio"
	"dict-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

// AudioHandler 音频处理器
type AudioHandler struct {
	audioService *audio.AudioService
}

// NewAudioHandler 创建音频处理器
func NewAudioHandler(audioService *audio.AudioService) *AudioHandler {
	return &AudioHandler{
		audioService: audioService,
	}
}

// GetAudio 获取单词音频
// GET /api/v1/audio/:word?type=gb|us
func (h *AudioHandler) GetAudio(c *gin.Context) {
	word := c.Param("word")
	audioType := c.Query("type") // 新增可选参数

	if word == "" {
		response.BadRequest(c, "word parameter is required")
		return
	}

	var reader io.Reader
	var contentType string
	var err error

	if audioType != "" {
		// 有类型参数，使用新方法
		reader, contentType, err = h.audioService.GetAudioByType(word, audioType)
	} else {
		// 无类型参数，保持原有逻辑（向后兼容）
		reader, contentType, err = h.audioService.GetAudio(word)
	}

	if err != nil {
		if err == audio.ErrAudioNotFound {
			response.NotFound(c, "audio not found for word: "+word)
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	// 设置缓存头
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400") // 缓存 24 小时
	c.Status(http.StatusOK)
	io.Copy(c.Writer, reader)
}

// CheckAudioAvailability 检查音频可用性
// GET /api/v1/audio/:word/availability
func (h *AudioHandler) CheckAudioAvailability(c *gin.Context) {
	word := c.Param("word")
	if word == "" {
		response.BadRequest(c, "word parameter is required")
		return
	}

	hasGB, hasUS := h.audioService.CheckAudioAvailability(word)
	response.Success(c, gin.H{
		"gb": hasGB,
		"us": hasUS,
	})
}
