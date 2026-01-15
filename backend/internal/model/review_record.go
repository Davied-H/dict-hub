package model

import (
	"time"

	"gorm.io/gorm"
)

// ReviewResult 复习结果类型
type ReviewResult string

const (
	ReviewResultForgot ReviewResult = "forgot" // 完全忘记
	ReviewResultHard   ReviewResult = "hard"   // 记得但困难
	ReviewResultGood   ReviewResult = "good"   // 记得
	ReviewResultEasy   ReviewResult = "easy"   // 非常容易
)

// ReviewRecord 复习记录
type ReviewRecord struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	VocabularyID uint           `gorm:"not null;index" json:"vocabulary_id"`
	Result       ReviewResult   `gorm:"size:20;not null" json:"result"`
	LevelBefore  int            `json:"level_before"`
	LevelAfter   int            `json:"level_after"`
	ResponseTime int64          `json:"response_time"`
	ReviewedAt   time.Time      `gorm:"index" json:"reviewed_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ReviewRecord) TableName() string {
	return "review_records"
}
