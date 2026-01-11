package service

import (
	"errors"
	"os"
	"sync"

	"dict-hub/internal/model"
	"dict-hub/internal/service/mdx"

	"gorm.io/gorm"
)

var (
	ErrDictSourceNotFound = errors.New("dictionary source not found")
	ErrDictFileNotFound   = errors.New("dictionary file not found")
	ErrDictAlreadyExists  = errors.New("dictionary already exists")
)

// ReorderItem 排序项
type ReorderItem struct {
	ID        uint `json:"id"`
	SortOrder int  `json:"sort_order"`
}

// DictSourceService 字典来源服务
type DictSourceService struct {
	db         *gorm.DB
	mdxManager mdx.DictManager
	dictDir    string
	runtimeIDs map[uint]uint // DB ID -> Runtime ID 映射
	mu         sync.RWMutex
}

// NewDictSourceService 创建字典来源服务
func NewDictSourceService(db *gorm.DB, mdxManager mdx.DictManager, dictDir string) *DictSourceService {
	return &DictSourceService{
		db:         db,
		mdxManager: mdxManager,
		dictDir:    dictDir,
		runtimeIDs: make(map[uint]uint),
	}
}

// DictSourceResponse 字典响应（包含加载状态）
type DictSourceResponse struct {
	model.DictSource
	Loaded bool `json:"loaded"`
}

// List 列出所有字典，按 SortOrder 排序
func (s *DictSourceService) List() ([]DictSourceResponse, error) {
	var sources []model.DictSource
	if err := s.db.Order("sort_order ASC, id ASC").Find(&sources).Error; err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	responses := make([]DictSourceResponse, len(sources))
	for i, src := range sources {
		_, loaded := s.runtimeIDs[src.ID]
		responses[i] = DictSourceResponse{
			DictSource: src,
			Loaded:     loaded,
		}
	}

	return responses, nil
}

// Add 添加字典
func (s *DictSourceService) Add(path string) (*DictSourceResponse, error) {
	// 检查文件是否存在
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, ErrDictFileNotFound
	}

	// 检查是否已存在
	var existing model.DictSource
	if err := s.db.Where("path = ?", path).First(&existing).Error; err == nil {
		return nil, ErrDictAlreadyExists
	}

	// 加载到 MDX 管理器获取元信息
	runtimeID, err := s.mdxManager.LoadDict(path)
	if err != nil {
		return nil, err
	}

	// 从加载的字典获取元信息
	var dictInfo mdx.DictInfo
	for _, info := range s.mdxManager.ListLoaded() {
		if info.ID == runtimeID {
			dictInfo = info
			break
		}
	}

	// 获取最大排序值
	var maxOrder int
	s.db.Model(&model.DictSource{}).Select("COALESCE(MAX(sort_order), -1)").Scan(&maxOrder)

	// 创建数据库记录
	source := &model.DictSource{
		Name:        dictInfo.Name,
		Title:       dictInfo.Title,
		Description: dictInfo.Description,
		Path:        path,
		Enabled:     true,
		SortOrder:   maxOrder + 1,
		WordCount:   dictInfo.WordCount,
		HasMDD:      dictInfo.HasMDD,
		FileSize:    fileInfo.Size(),
	}

	if err := s.db.Create(source).Error; err != nil {
		// 回滚：卸载已加载的字典
		s.mdxManager.Unload(runtimeID)
		return nil, err
	}

	// 维护 ID 映射
	s.mu.Lock()
	s.runtimeIDs[source.ID] = runtimeID
	s.mu.Unlock()

	return &DictSourceResponse{
		DictSource: *source,
		Loaded:     true,
	}, nil
}

// Toggle 切换字典启用状态
func (s *DictSourceService) Toggle(id uint) (*DictSourceResponse, error) {
	var source model.DictSource
	if err := s.db.First(&source, id).Error; err != nil {
		return nil, ErrDictSourceNotFound
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if source.Enabled {
		// 禁用：卸载字典
		if runtimeID, ok := s.runtimeIDs[id]; ok {
			if err := s.mdxManager.Unload(runtimeID); err != nil {
				return nil, err
			}
			delete(s.runtimeIDs, id)
		}
		source.Enabled = false
	} else {
		// 启用：加载字典
		runtimeID, err := s.mdxManager.LoadDict(source.Path)
		if err != nil {
			return nil, err
		}
		s.runtimeIDs[id] = runtimeID
		source.Enabled = true
	}

	if err := s.db.Save(&source).Error; err != nil {
		return nil, err
	}

	return &DictSourceResponse{
		DictSource: source,
		Loaded:     source.Enabled,
	}, nil
}

// Reorder 重新排序字典
func (s *DictSourceService) Reorder(orders []ReorderItem) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range orders {
			if err := tx.Model(&model.DictSource{}).Where("id = ?", item.ID).Update("sort_order", item.SortOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Delete 删除字典
func (s *DictSourceService) Delete(id uint) error {
	var source model.DictSource
	if err := s.db.First(&source, id).Error; err != nil {
		return ErrDictSourceNotFound
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 从 MDX 管理器卸载
	if runtimeID, ok := s.runtimeIDs[id]; ok {
		s.mdxManager.Unload(runtimeID)
		delete(s.runtimeIDs, id)
	}

	// 软删除数据库记录
	return s.db.Delete(&source).Error
}

// GetByID 根据 ID 获取字典
func (s *DictSourceService) GetByID(id uint) (*model.DictSource, error) {
	var source model.DictSource
	if err := s.db.First(&source, id).Error; err != nil {
		return nil, ErrDictSourceNotFound
	}
	return &source, nil
}

// GetEnabled 获取已启用的字典列表
func (s *DictSourceService) GetEnabled() ([]model.DictSource, error) {
	var sources []model.DictSource
	if err := s.db.Where("enabled = ?", true).Order("sort_order ASC").Find(&sources).Error; err != nil {
		return nil, err
	}
	return sources, nil
}

// SyncOnStartup 启动时同步：加载数据库中已启用的字典
func (s *DictSourceService) SyncOnStartup() error {
	sources, err := s.GetEnabled()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, src := range sources {
		// 检查文件是否存在
		if _, err := os.Stat(src.Path); err != nil {
			// 文件不存在，禁用该字典
			s.db.Model(&src).Update("enabled", false)
			continue
		}

		// 加载字典
		runtimeID, err := s.mdxManager.LoadDict(src.Path)
		if err != nil {
			// 加载失败，禁用该字典
			s.db.Model(&src).Update("enabled", false)
			continue
		}

		// 维护 ID 映射
		s.runtimeIDs[src.ID] = runtimeID
	}

	return nil
}

// GetRuntimeID 获取字典的运行时 ID（供其他服务使用）
func (s *DictSourceService) GetRuntimeID(dbID uint) (uint, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	runtimeID, ok := s.runtimeIDs[dbID]
	return runtimeID, ok
}
