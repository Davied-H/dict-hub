package handler

import (
	"errors"
	"strconv"

	"dict-hub/internal/model"
	"dict-hub/internal/service/vocabulary"
	"dict-hub/pkg/response"

	"github.com/gin-gonic/gin"
)

// ReviewHandler 复习处理器
type ReviewHandler struct {
	reviewService *vocabulary.ReviewService
}

// NewReviewHandler 创建复习处理器
func NewReviewHandler(reviewService *vocabulary.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// GetDue 获取待复习的单词列表
// GET /api/v1/review/due?limit=20
func (h *ReviewHandler) GetDue(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	vocabs, err := h.reviewService.GetDueVocabularies(limit)
	if err != nil {
		response.InternalError(c, "failed to get due vocabularies: "+err.Error())
		return
	}

	response.Success(c, vocabs)
}

// SubmitReviewRequest 提交复习请求
type SubmitReviewRequest struct {
	Result       string `json:"result" binding:"required,oneof=forgot hard good easy"`
	ResponseTime int64  `json:"response_time"`
}

// Submit 提交复习结果
// POST /api/v1/review/:id
func (h *ReviewHandler) Submit(c *gin.Context) {
	vocabularyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid vocabulary id")
		return
	}

	var req SubmitReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	result := model.ReviewResult(req.Result)

	resp, err := h.reviewService.SubmitReview(uint(vocabularyID), result, req.ResponseTime)
	if err != nil {
		if errors.Is(err, vocabulary.ErrVocabularyNotFound) {
			response.NotFound(c, "vocabulary not found")
			return
		}
		response.InternalError(c, "failed to submit review: "+err.Error())
		return
	}

	response.Success(c, resp)
}

// GetStats 获取复习统计
// GET /api/v1/review/stats
func (h *ReviewHandler) GetStats(c *gin.Context) {
	stats, err := h.reviewService.GetStats()
	if err != nil {
		response.InternalError(c, "failed to get stats: "+err.Error())
		return
	}

	response.Success(c, stats)
}

// GetHistory 获取复习历史记录
// GET /api/v1/review/history?page=1&page_size=20
func (h *ReviewHandler) GetHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	history, total, err := h.reviewService.GetHistory(page, pageSize)
	if err != nil {
		response.InternalError(c, "failed to get history: "+err.Error())
		return
	}

	response.PagedSuccess(c, history, total, page, pageSize)
}
