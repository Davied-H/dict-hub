package model

import (
	"time"

	"gorm.io/gorm"
)

type WordFrequency struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Word         string         `gorm:"size:255;not null;uniqueIndex" json:"word"`
	SearchCount  int64          `gorm:"default:0" json:"search_count"`
	ViewCount    int64          `gorm:"default:0" json:"view_count"`
	LastSearched time.Time      `json:"last_searched"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (WordFrequency) TableName() string {
	return "word_frequencies"
}
