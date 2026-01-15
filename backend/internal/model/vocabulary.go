package model

import (
	"time"

	"gorm.io/gorm"
)

// Vocabulary 生词本 - 收藏的单词
type Vocabulary struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Word       string         `gorm:"size:255;not null;uniqueIndex" json:"word"`
	DictID     uint           `gorm:"index" json:"dict_id"`
	DictTitle  string         `gorm:"size:255" json:"dict_title"`
	Definition string         `gorm:"type:text" json:"definition"`
	Phonetic   string         `gorm:"size:255" json:"phonetic"`

	// 复习相关
	Level        int        `gorm:"default:0" json:"level"`
	NextReviewAt *time.Time `gorm:"index" json:"next_review_at"`
	LastReviewAt *time.Time `json:"last_review_at"`
	ReviewCount  int        `gorm:"default:0" json:"review_count"`

	// 标签
	Tags string `gorm:"size:500" json:"tags"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Vocabulary) TableName() string {
	return "vocabularies"
}
