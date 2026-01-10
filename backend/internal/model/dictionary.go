package model

import (
	"time"

	"gorm.io/gorm"
)

type Dictionary struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Word         string         `gorm:"size:255;not null;uniqueIndex" json:"word"`
	Phonetic     string         `gorm:"size:255" json:"phonetic"`
	Definition   string         `gorm:"type:text;not null" json:"definition"`
	Example      string         `gorm:"type:text" json:"example"`
	PartOfSpeech string         `gorm:"size:50" json:"part_of_speech"`
	Synonyms     string         `gorm:"type:text" json:"synonyms"`
	Antonyms     string         `gorm:"type:text" json:"antonyms"`
	Difficulty   int            `gorm:"default:1" json:"difficulty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Dictionary) TableName() string {
	return "dictionaries"
}
