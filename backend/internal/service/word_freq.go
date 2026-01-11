package service

import (
	"bufio"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"time"

	"dict-hub/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// WordFreqService 词频服务
type WordFreqService struct {
	db *gorm.DB
}

// NewWordFreqService 创建词频服务
func NewWordFreqService(db *gorm.DB) *WordFreqService {
	return &WordFreqService{db: db}
}

// ImportResult 导入结果
type ImportResult struct {
	TotalLines    int      `json:"total_lines"`
	ImportedCount int      `json:"imported_count"`
	SkippedCount  int      `json:"skipped_count"`
	Errors        []string `json:"errors,omitempty"`
}

// Import 导入词频文件（CSV格式: word,count）
func (s *WordFreqService) Import(reader io.Reader) (*ImportResult, error) {
	result := &ImportResult{
		Errors: make([]string, 0),
	}

	csvReader := csv.NewReader(bufio.NewReader(reader))
	csvReader.FieldsPerRecord = -1 // 允许可变字段数
	csvReader.TrimLeadingSpace = true

	lineNum := 0
	batch := make([]*model.WordFrequency, 0, 100)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		lineNum++
		result.TotalLines++

		if err != nil {
			result.Errors = append(result.Errors, "line "+strconv.Itoa(lineNum)+": "+err.Error())
			result.SkippedCount++
			continue
		}

		// 跳过表头
		if lineNum == 1 && len(record) >= 1 && strings.ToLower(strings.TrimSpace(record[0])) == "word" {
			result.TotalLines--
			continue
		}

		// 至少需要 word 字段
		if len(record) < 1 {
			result.SkippedCount++
			continue
		}

		word := strings.TrimSpace(record[0])
		if word == "" {
			result.SkippedCount++
			continue
		}

		var count int64 = 1
		if len(record) >= 2 {
			c, err := strconv.ParseInt(strings.TrimSpace(record[1]), 10, 64)
			if err == nil && c > 0 {
				count = c
			}
		}

		batch = append(batch, &model.WordFrequency{
			Word:         word,
			SearchCount:  count,
			LastSearched: time.Now(),
		})

		// 批量插入
		if len(batch) >= 100 {
			if err := s.batchUpsert(batch); err != nil {
				result.Errors = append(result.Errors, "batch insert error: "+err.Error())
			} else {
				result.ImportedCount += len(batch)
			}
			batch = batch[:0]
		}
	}

	// 处理剩余数据
	if len(batch) > 0 {
		if err := s.batchUpsert(batch); err != nil {
			result.Errors = append(result.Errors, "batch insert error: "+err.Error())
		} else {
			result.ImportedCount += len(batch)
		}
	}

	return result, nil
}

// batchUpsert 批量插入或更新
func (s *WordFreqService) batchUpsert(items []*model.WordFrequency) error {
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "word"}},
		DoUpdates: clause.AssignmentColumns([]string{"search_count", "last_searched", "updated_at"}),
	}).Create(&items).Error
}

// IncrementSearch 增加单词搜索次数
func (s *WordFreqService) IncrementSearch(word string) error {
	return s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "word"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"search_count":  gorm.Expr("search_count + 1"),
			"last_searched": time.Now(),
			"updated_at":    time.Now(),
		}),
	}).Create(&model.WordFrequency{
		Word:         word,
		SearchCount:  1,
		LastSearched: time.Now(),
	}).Error
}
