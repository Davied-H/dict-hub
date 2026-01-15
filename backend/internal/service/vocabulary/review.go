package vocabulary

import (
	"time"

	"dict-hub/internal/model"

	"gorm.io/gorm"
)

// ReviewService 复习服务
type ReviewService struct {
	db       *gorm.DB
	vocabSvc *VocabularyService
}

// NewReviewService 创建复习服务
func NewReviewService(db *gorm.DB, vocabSvc *VocabularyService) *ReviewService {
	return &ReviewService{
		db:       db,
		vocabSvc: vocabSvc,
	}
}

// GetDueVocabularies 获取待复习的单词列表
func (s *ReviewService) GetDueVocabularies(limit int) ([]model.Vocabulary, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var vocabs []model.Vocabulary
	now := time.Now()

	// 查询需要复习的单词（next_review_at <= now 或者 next_review_at 为空）
	if err := s.db.Where("next_review_at IS NULL OR next_review_at <= ?", now).
		Order("next_review_at ASC, level ASC").
		Limit(limit).
		Find(&vocabs).Error; err != nil {
		return nil, err
	}

	return vocabs, nil
}

// SubmitReview 提交复习结果
func (s *ReviewService) SubmitReview(vocabularyID uint, result model.ReviewResult, responseTime int64) (*ReviewResponse, error) {
	// 获取当前单词
	vocab, err := s.vocabSvc.GetByID(vocabularyID)
	if err != nil {
		return nil, err
	}

	// 计算新等级和下次复习时间
	oldLevel := vocab.Level
	newLevel, nextReview := CalculateNextReview(oldLevel, result)

	// 更新单词
	now := time.Now()
	vocab.Level = newLevel
	vocab.NextReviewAt = &nextReview
	vocab.LastReviewAt = &now
	vocab.ReviewCount++

	if err := s.db.Save(vocab).Error; err != nil {
		return nil, err
	}

	// 记录复习历史
	record := &model.ReviewRecord{
		VocabularyID: vocabularyID,
		Result:       result,
		LevelBefore:  oldLevel,
		LevelAfter:   newLevel,
		ResponseTime: responseTime,
		ReviewedAt:   now,
	}
	if err := s.db.Create(record).Error; err != nil {
		return nil, err
	}

	return &ReviewResponse{
		VocabularyID: vocabularyID,
		NewLevel:     newLevel,
		NextReviewAt: nextReview,
		LevelChange:  newLevel - oldLevel,
	}, nil
}

// ReviewResponse 复习响应
type ReviewResponse struct {
	VocabularyID uint      `json:"vocabulary_id"`
	NewLevel     int       `json:"new_level"`
	NextReviewAt time.Time `json:"next_review_at"`
	LevelChange  int       `json:"level_change"`
}

// GetStats 获取复习统计
func (s *ReviewService) GetStats() (*ReviewStats, error) {
	stats := &ReviewStats{}
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 总词汇量
	s.db.Model(&model.Vocabulary{}).Count(&stats.TotalWords)

	// 今日待复习
	s.db.Model(&model.Vocabulary{}).
		Where("next_review_at IS NULL OR next_review_at <= ?", now).
		Count(&stats.DueToday)

	// 已掌握 (Level >= 6)
	s.db.Model(&model.Vocabulary{}).Where("level >= 6").Count(&stats.Mastered)

	// 学习中 (Level 1-5)
	s.db.Model(&model.Vocabulary{}).Where("level BETWEEN 1 AND 5").Count(&stats.Learning)

	// 新词 (Level 0)
	s.db.Model(&model.Vocabulary{}).Where("level = 0").Count(&stats.New)

	// 今日已复习
	s.db.Model(&model.ReviewRecord{}).
		Where("reviewed_at >= ?", todayStart).
		Count(&stats.ReviewedToday)

	// 连续学习天数
	stats.StreakDays = s.calculateStreakDays()

	return stats, nil
}

// calculateStreakDays 计算连续学习天数
func (s *ReviewService) calculateStreakDays() int {
	streak := 0
	now := time.Now()

	for i := 0; i < 365; i++ { // 最多检查一年
		dayStart := time.Date(now.Year(), now.Month(), now.Day()-i, 0, 0, 0, 0, now.Location())
		dayEnd := dayStart.Add(24 * time.Hour)

		var count int64
		s.db.Model(&model.ReviewRecord{}).
			Where("reviewed_at >= ? AND reviewed_at < ?", dayStart, dayEnd).
			Count(&count)

		if count > 0 {
			streak++
		} else if i > 0 { // 允许今天还没学习
			break
		}
	}

	return streak
}

// ReviewStats 复习统计
type ReviewStats struct {
	TotalWords    int64 `json:"total_words"`
	DueToday      int64 `json:"due_today"`
	Mastered      int64 `json:"mastered"`
	Learning      int64 `json:"learning"`
	New           int64 `json:"new"`
	ReviewedToday int64 `json:"reviewed_today"`
	StreakDays    int   `json:"streak_days"`
}

// GetHistory 获取复习历史记录
func (s *ReviewService) GetHistory(page, pageSize int) ([]ReviewHistoryItem, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	if err := s.db.Model(&model.ReviewRecord{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var records []model.ReviewRecord
	offset := (page - 1) * pageSize
	if err := s.db.Order("reviewed_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	// 获取关联的单词信息
	var items []ReviewHistoryItem
	for _, r := range records {
		item := ReviewHistoryItem{
			ID:           r.ID,
			VocabularyID: r.VocabularyID,
			Result:       string(r.Result),
			LevelBefore:  r.LevelBefore,
			LevelAfter:   r.LevelAfter,
			ResponseTime: r.ResponseTime,
			ReviewedAt:   r.ReviewedAt,
		}

		// 获取单词
		var vocab model.Vocabulary
		if err := s.db.First(&vocab, r.VocabularyID).Error; err == nil {
			item.Word = vocab.Word
		}

		items = append(items, item)
	}

	return items, total, nil
}

// ReviewHistoryItem 复习历史项
type ReviewHistoryItem struct {
	ID           uint      `json:"id"`
	VocabularyID uint      `json:"vocabulary_id"`
	Word         string    `json:"word"`
	Result       string    `json:"result"`
	LevelBefore  int       `json:"level_before"`
	LevelAfter   int       `json:"level_after"`
	ResponseTime int64     `json:"response_time"`
	ReviewedAt   time.Time `json:"reviewed_at"`
}
