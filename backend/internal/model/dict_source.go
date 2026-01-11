package model

import (
	"time"

	"gorm.io/gorm"
)

// DictSource 字典来源/文件元数据
// 区别于 Dictionary（词条模型），这是字典文件本身的元信息
type DictSource struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`                   // 字典名称（来自MDX元数据）
	Title       string         `gorm:"size:255" json:"title"`                           // 字典标题
	Description string         `gorm:"type:text" json:"description"`                    // 字典描述
	Path        string         `gorm:"size:1024;not null;uniqueIndex" json:"path"`      // MDX文件路径
	Enabled     bool           `gorm:"default:true" json:"enabled"`                     // 是否启用
	SortOrder   int            `gorm:"default:0;index" json:"sort_order"`               // 排序顺序
	WordCount   int64          `gorm:"default:0" json:"word_count"`                     // 词条数量
	HasMDD      bool           `gorm:"default:false" json:"has_mdd"`                    // 是否有MDD资源文件
	SourceURL   string         `gorm:"size:1024" json:"source_url,omitempty"`           // 下载来源URL（可选）
	FileSize    int64          `gorm:"default:0" json:"file_size"`                      // 文件大小（字节）
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DictSource) TableName() string {
	return "dict_sources"
}
