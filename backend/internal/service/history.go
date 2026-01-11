package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"dict-hub/internal/model"

	"gorm.io/gorm"
)

// HistoryService 搜索历史服务
type HistoryService struct {
	db *gorm.DB
}

// NewHistoryService 创建历史服务
func NewHistoryService(db *gorm.DB) *HistoryService {
	return &HistoryService{db: db}
}

// List 分页查询历史记录
func (s *HistoryService) List(page, pageSize int) ([]model.SearchHistory, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	if err := s.db.Model(&model.SearchHistory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var histories []model.SearchHistory
	offset := (page - 1) * pageSize
	if err := s.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// Record 记录搜索历史
func (s *HistoryService) Record(word, sessionID, ip string, found bool, responseTime int64) error {
	history := &model.SearchHistory{
		Word:         word,
		SessionID:    sessionID,
		IPAddress:    ip,
		Found:        found,
		ResponseTime: responseTime,
		CreatedAt:    time.Now(),
	}
	return s.db.Create(history).Error
}

// Clear 清空历史记录（软删除）
func (s *HistoryService) Clear() error {
	return s.db.Where("1 = 1").Delete(&model.SearchHistory{}).Error
}

// ExportCSV 导出CSV格式的历史记录
func (s *HistoryService) ExportCSV() ([]byte, error) {
	var histories []model.SearchHistory
	if err := s.db.Order("created_at DESC").Find(&histories).Error; err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// 写入表头
	if err := writer.Write([]string{"word", "found", "response_time_ms", "created_at"}); err != nil {
		return nil, err
	}

	// 写入数据行
	for _, h := range histories {
		record := []string{
			h.Word,
			fmt.Sprintf("%t", h.Found),
			fmt.Sprintf("%d", h.ResponseTime),
			h.CreatedAt.Format(time.RFC3339),
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
