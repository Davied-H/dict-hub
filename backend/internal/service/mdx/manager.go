package mdx

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	libmdx "github.com/lib-x/mdx"
)

var (
	ErrDictNotFound     = errors.New("dictionary not found")
	ErrResourceNotFound = errors.New("resource not found")
	ErrNoMDD            = errors.New("no MDD resource file associated")
	ErrWordNotFound     = errors.New("word not found")
)

// dictEntry 内部字典条目
type dictEntry struct {
	id   uint
	mdx  *libmdx.Mdict
	mdd  *libmdx.Mdict // 可选的 MDD 资源文件
	path string
}

// manager DictManager 实现
type manager struct {
	mu     sync.RWMutex
	dicts  map[uint]*dictEntry
	nextID uint
}

// NewManager 创建新的字典管理器
func NewManager() DictManager {
	return &manager{
		dicts:  make(map[uint]*dictEntry),
		nextID: 1,
	}
}

// LoadDict 加载单个 MDX 字典文件
func (m *manager) LoadDict(path string) (uint, error) {
	// 验证文件存在
	if _, err := os.Stat(path); err != nil {
		return 0, err
	}

	// 创建 MDX 实例
	mdx, err := libmdx.New(path)
	if err != nil {
		return 0, err
	}

	// 构建索引
	if err := mdx.BuildIndex(); err != nil {
		return 0, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	id := m.nextID
	m.nextID++

	entry := &dictEntry{
		id:   id,
		mdx:  mdx,
		path: path,
	}

	// 尝试加载同名 MDD 文件
	mddPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".mdd"
	if _, err := os.Stat(mddPath); err == nil {
		mdd, err := libmdx.New(mddPath)
		if err == nil {
			if err := mdd.BuildIndex(); err == nil {
				entry.mdd = mdd
			}
		}
	}

	m.dicts[id] = entry
	return id, nil
}

// LoadAll 扫描目录加载所有 MDX 字典
func (m *manager) LoadAll(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var loadErrors []error
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.ToLower(filepath.Ext(entry.Name())) == ".mdx" {
			fullPath := filepath.Join(dir, entry.Name())
			if _, err := m.LoadDict(fullPath); err != nil {
				loadErrors = append(loadErrors, err)
			}
		}
	}

	if len(loadErrors) > 0 {
		return errors.Join(loadErrors...)
	}
	return nil
}

// Lookup 在指定字典中查询单词
func (m *manager) Lookup(dictID uint, word string) ([]byte, error) {
	m.mu.RLock()
	entry, ok := m.dicts[dictID]
	m.mu.RUnlock()

	if !ok {
		return nil, ErrDictNotFound
	}

	result, err := entry.mdx.Lookup(word)
	if err != nil {
		return nil, ErrWordNotFound
	}

	return result, nil
}

// Search 跨字典搜索单词
func (m *manager) Search(word string, dictIDs ...uint) []SearchResult {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []SearchResult

	// 如果没有指定字典 ID，则搜索所有字典
	if len(dictIDs) == 0 {
		for id := range m.dicts {
			dictIDs = append(dictIDs, id)
		}
	}

	for _, id := range dictIDs {
		entry, ok := m.dicts[id]
		if !ok {
			continue
		}

		result, err := entry.mdx.Lookup(word)
		if err != nil {
			continue
		}

		results = append(results, SearchResult{
			DictID:     id,
			DictName:   entry.mdx.Name(),
			DictTitle:  entry.mdx.Title(),
			Word:       word,
			Definition: string(result),
		})
	}

	return results
}

// Suggest 跨字典前缀搜索建议
func (m *manager) Suggest(prefix string, limit int) []SuggestResult {
	m.mu.RLock()
	defer m.mu.RUnlock()

	prefix = strings.ToLower(prefix)
	var results []SuggestResult
	seen := make(map[string]bool) // 去重

	for _, entry := range m.dicts {
		keywords, err := entry.mdx.GetKeyWordEntries()
		if err != nil {
			continue
		}

		for _, kw := range keywords {
			if strings.HasPrefix(strings.ToLower(kw.KeyWord), prefix) {
				// 去重：同一个词只返回一次
				if seen[kw.KeyWord] {
					continue
				}
				seen[kw.KeyWord] = true

				results = append(results, SuggestResult{
					Word:      kw.KeyWord,
					DictID:    entry.id,
					DictTitle: entry.mdx.Title(),
				})
				if len(results) >= limit {
					return results
				}
			}
		}
	}
	return results
}

// GetResource 获取 MDD 资源文件
func (m *manager) GetResource(dictID uint, path string) (io.Reader, error) {
	m.mu.RLock()
	entry, ok := m.dicts[dictID]
	m.mu.RUnlock()

	if !ok {
		return nil, ErrDictNotFound
	}

	if entry.mdd == nil {
		return nil, ErrNoMDD
	}

	// 标准化路径
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimPrefix(path, "\\")

	data, err := entry.mdd.Lookup(path)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return bytes.NewReader(data), nil
}

// ListLoaded 列出已加载的字典
func (m *manager) ListLoaded() []DictInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	infos := make([]DictInfo, 0, len(m.dicts))
	for _, entry := range m.dicts {
		info := DictInfo{
			ID:          entry.id,
			Name:        entry.mdx.Name(),
			Title:       entry.mdx.Title(),
			Description: entry.mdx.Description(),
			Path:        entry.path,
			HasMDD:      entry.mdd != nil,
			WordCount:   entry.mdx.GetKeyWordEntriesSize(),
		}
		infos = append(infos, info)
	}
	return infos
}

// Unload 卸载指定字典
func (m *manager) Unload(dictID uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.dicts[dictID]; !ok {
		return ErrDictNotFound
	}

	delete(m.dicts, dictID)
	return nil
}
