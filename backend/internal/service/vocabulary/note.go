package vocabulary

import (
	"errors"
	"time"

	"dict-hub/internal/model"

	"gorm.io/gorm"
)

var (
	ErrNoteNotFound = errors.New("note not found")
)

// NoteService 笔记服务
type NoteService struct {
	db *gorm.DB
}

// NewNoteService 创建笔记服务
func NewNoteService(db *gorm.DB) *NoteService {
	return &NoteService{db: db}
}

// Create 创建笔记
func (s *NoteService) Create(note *model.Note) error {
	// 验证 vocabulary 存在
	var vocab model.Vocabulary
	if err := s.db.First(&vocab, note.VocabularyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVocabularyNotFound
		}
		return err
	}

	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()
	return s.db.Create(note).Error
}

// GetByID 根据ID获取笔记
func (s *NoteService) GetByID(id uint) (*model.Note, error) {
	var note model.Note
	if err := s.db.First(&note, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}
	return &note, nil
}

// Update 更新笔记
func (s *NoteService) Update(id uint, content string) error {
	result := s.db.Model(&model.Note{}).Where("id = ?", id).Updates(map[string]interface{}{
		"content":    content,
		"updated_at": time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoteNotFound
	}
	return nil
}

// Delete 删除笔记
func (s *NoteService) Delete(id uint) error {
	result := s.db.Delete(&model.Note{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoteNotFound
	}
	return nil
}

// ListByVocabulary 获取某单词的所有笔记
func (s *NoteService) ListByVocabulary(vocabularyID uint) ([]model.Note, error) {
	var notes []model.Note
	if err := s.db.Where("vocabulary_id = ?", vocabularyID).Order("created_at DESC").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

// CountByVocabulary 获取某单词的笔记数量
func (s *NoteService) CountByVocabulary(vocabularyID uint) (int64, error) {
	var count int64
	if err := s.db.Model(&model.Note{}).Where("vocabulary_id = ?", vocabularyID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
