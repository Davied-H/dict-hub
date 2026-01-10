package model

import (
	"time"

	"gorm.io/gorm"
)

type SearchHistory struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Word         string         `gorm:"size:255;not null;index" json:"word"`
	SessionID    string         `gorm:"size:100;index" json:"session_id"`
	IPAddress    string         `gorm:"size:50" json:"ip_address"`
	Found        bool           `gorm:"default:false" json:"found"`
	ResponseTime int64          `json:"response_time"`
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SearchHistory) TableName() string {
	return "search_histories"
}
