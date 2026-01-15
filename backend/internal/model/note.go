package model

import (
	"time"

	"gorm.io/gorm"
)

// Note 单词笔记
type Note struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	VocabularyID uint           `gorm:"not null;index" json:"vocabulary_id"`
	Content      string         `gorm:"type:text;not null" json:"content"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Note) TableName() string {
	return "notes"
}
