package mdx

import "io"

// DictInfo 字典元信息
type DictInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Path        string `json:"path"`
	HasMDD      bool   `json:"has_mdd"`
	WordCount   int64  `json:"word_count"`
}

// SearchResult 搜索结果
type SearchResult struct {
	DictID     uint   `json:"dict_id"`
	DictName   string `json:"dict_name"`
	DictTitle  string `json:"dict_title"`
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

// SuggestResult 搜索建议结果
type SuggestResult struct {
	Word      string `json:"word"`
	DictID    uint   `json:"dict_id"`
	DictTitle string `json:"dict_title"`
}

// DictManager 字典管理器接口
type DictManager interface {
	// LoadDict 加载单个 MDX 字典文件
	LoadDict(path string) (uint, error)

	// LoadAll 扫描目录加载所有 MDX 字典
	LoadAll(dir string) error

	// Lookup 在指定字典中查询单词
	Lookup(dictID uint, word string) ([]byte, error)

	// Search 跨字典搜索单词
	Search(word string, dictIDs ...uint) []SearchResult

	// Suggest 前缀搜索建议
	Suggest(prefix string, limit int) []SuggestResult

	// GetResource 获取 MDD 资源文件
	GetResource(dictID uint, path string) (io.Reader, error)

	// ListLoaded 列出已加载的字典
	ListLoaded() []DictInfo

	// Unload 卸载指定字典
	Unload(dictID uint) error
}
