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
// GET /api/v1/audio/:word
func (h *AudioHandler) GetAudio(c *gin.Context) {
	word := c.Param("word")
	if word == "" {
		response.BadRequest(c, "word parameter is required")
		return
	}

	reader, contentType, err := h.audioService.GetAudio(word)
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
