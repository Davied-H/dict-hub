package model

import (
	"time"

	"gorm.io/gorm"
)

// 下载任务状态常量
const (
	DownloadStatusPending     = "pending"
	DownloadStatusDownloading = "downloading"
	DownloadStatusCompleted   = "completed"
	DownloadStatusFailed      = "failed"
)

// DownloadTask 字典下载任务
type DownloadTask struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	URL          string         `gorm:"size:2048;not null" json:"url"`                      // 下载URL
	FileName     string         `gorm:"size:255" json:"file_name"`                          // 文件名
	Status       string         `gorm:"size:20;default:'pending';index" json:"status"`      // pending/downloading/completed/failed
	Progress     int            `gorm:"default:0" json:"progress"`                          // 进度 0-100
	TotalSize    int64          `gorm:"default:0" json:"total_size"`                        // 总大小（字节）
	DownloadSize int64          `gorm:"default:0" json:"download_size"`                     // 已下载大小（字节）
	ErrorMsg     string         `gorm:"type:text" json:"error_msg,omitempty"`               // 错误信息
	DictSourceID *uint          `json:"dict_source_id,omitempty"`                           // 下载完成后关联的字典ID
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DownloadTask) TableName() string {
	return "download_tasks"
}
