package vocabulary

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"time"

	"dict-hub/internal/model"

	"gorm.io/gorm"
)

var (
	ErrVocabularyNotFound = errors.New("vocabulary not found")
	ErrWordAlreadyExists  = errors.New("word already exists in vocabulary")
)

// VocabularyService 生词本服务
type VocabularyService struct {
	db *gorm.DB
}

// NewVocabularyService 创建生词本服务
func NewVocabularyService(db *gorm.DB) *VocabularyService {
	return &VocabularyService{db: db}
}

// Add 添加单词到生词本
func (s *VocabularyService) Add(vocab *model.Vocabulary) error {
	// 检查是否已存在
	var existing model.Vocabulary
	if err := s.db.Where("word = ?", vocab.Word).First(&existing).Error; err == nil {
		return ErrWordAlreadyExists
	}

	// 设置初始复习时间
	now := time.Now()
	vocab.NextReviewAt = &now
	vocab.CreatedAt = now
	vocab.UpdatedAt = now

	return s.db.Create(vocab).Error
}

// Remove 从生词本移除单词
func (s *VocabularyService) Remove(id uint) error {
	result := s.db.Delete(&model.Vocabulary{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrVocabularyNotFound
	}
	// 同时删除相关笔记
	s.db.Where("vocabulary_id = ?", id).Delete(&model.Note{})
	return nil
}

// GetByID 根据ID获取生词
func (s *VocabularyService) GetByID(id uint) (*model.Vocabulary, error) {
	var vocab model.Vocabulary
	if err := s.db.First(&vocab, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVocabularyNotFound
		}
		return nil, err
	}
	return &vocab, nil
}

// GetByWord 根据单词获取生词
func (s *VocabularyService) GetByWord(word string) (*model.Vocabulary, error) {
	var vocab model.Vocabulary
	if err := s.db.Where("word = ?", word).First(&vocab).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVocabularyNotFound
		}
		return nil, err
	}
	return &vocab, nil
}

// CheckWord 检查单词是否已收藏
func (s *VocabularyService) CheckWord(word string) (bool, uint) {
	var vocab model.Vocabulary
	if err := s.db.Where("word = ?", word).First(&vocab).Error; err != nil {
		return false, 0
	}
	return true, vocab.ID
}

// ListParams 列表查询参数
type ListParams struct {
	Page      int
	PageSize  int
	SortBy    string // created_at, next_review_at, level, word
	SortOrder string // asc, desc
	Level     *int
	Tag       string
	Keyword   string
}

// List 分页查询生词本
func (s *VocabularyService) List(params ListParams) ([]model.Vocabulary, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}

	query := s.db.Model(&model.Vocabulary{})

	// 筛选条件
	if params.Level != nil {
		query = query.Where("level = ?", *params.Level)
	}
	if params.Tag != "" {
		query = query.Where("tags LIKE ?", "%"+params.Tag+"%")
	}
	if params.Keyword != "" {
		query = query.Where("word LIKE ?", "%"+params.Keyword+"%")
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	orderBy := "created_at"
	if params.SortBy != "" {
		switch params.SortBy {
		case "created_at", "next_review_at", "level", "word":
			orderBy = params.SortBy
		}
	}
	order := "DESC"
	if params.SortOrder == "asc" {
		order = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", orderBy, order))

	// 分页
	var vocabs []model.Vocabulary
	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Find(&vocabs).Error; err != nil {
		return nil, 0, err
	}

	return vocabs, total, nil
}

// UpdateTags 更新标签
func (s *VocabularyService) UpdateTags(id uint, tags string) error {
	result := s.db.Model(&model.Vocabulary{}).Where("id = ?", id).Update("tags", tags)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrVocabularyNotFound
	}
	return nil
}

// ExportCSV 导出CSV格式的生词本
func (s *VocabularyService) ExportCSV() ([]byte, error) {
	var vocabs []model.Vocabulary
	if err := s.db.Order("created_at DESC").Find(&vocabs).Error; err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// 写入表头
	if err := writer.Write([]string{"word", "phonetic", "level", "review_count", "tags", "created_at"}); err != nil {
		return nil, err
	}

	// 写入数据行
	for _, v := range vocabs {
		record := []string{
			v.Word,
			v.Phonetic,
			fmt.Sprintf("%d", v.Level),
			fmt.Sprintf("%d", v.ReviewCount),
			v.Tags,
			v.CreatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GetStats 获取生词本统计
func (s *VocabularyService) GetStats() (*VocabularyStats, error) {
	stats := &VocabularyStats{}

	// 总数
	s.db.Model(&model.Vocabulary{}).Count(&stats.Total)

	// 按等级统计
	s.db.Model(&model.Vocabulary{}).Where("level = 0").Count(&stats.New)
	s.db.Model(&model.Vocabulary{}).Where("level BETWEEN 1 AND 5").Count(&stats.Learning)
	s.db.Model(&model.Vocabulary{}).Where("level >= 6").Count(&stats.Mastered)

	return stats, nil
}

// VocabularyStats 生词本统计
type VocabularyStats struct {
	Total    int64 `json:"total"`
	New      int64 `json:"new"`
	Learning int64 `json:"learning"`
	Mastered int64 `json:"mastered"`
}
