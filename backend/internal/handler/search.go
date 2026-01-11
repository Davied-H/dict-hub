package handler

import (
	"regexp"
	"sort"
	"strconv"
	"time"

	"dict-hub/internal/cache"
	"dict-hub/internal/model"
	"dict-hub/internal/service/mdx"
	"dict-hub/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	SearchCacheTTL  = 30 * time.Minute
	SuggestCacheTTL = 5 * time.Minute
	DefaultLimit    = 10
)

// SearchHandler 搜索处理器
type SearchHandler struct {
	manager mdx.DictManager
	cache   *cache.Cache
	db      *gorm.DB
}

// NewSearchHandler 创建搜索处理器
func NewSearchHandler(manager mdx.DictManager, cache *cache.Cache, db *gorm.DB) *SearchHandler {
	return &SearchHandler{manager: manager, cache: cache, db: db}
}

// Search 跨字典搜索
// GET /api/v1/search?word=xxx
func (h *SearchHandler) Search(c *gin.Context) {
	word := c.Query("word")
	if word == "" {
		response.BadRequest(c, "word parameter is required")
		return
	}

	// 检查缓存
	cacheKey := "search:" + word
	if cached, ok := h.cache.Get(cacheKey); ok {
		response.Success(c, cached)
		return
	}

	results := h.manager.Search(word)

	// URL 重写
	for i := range results {
		results[i].Definition = rewriteResourceURLs(results[i].Definition, results[i].DictID)
	}

	// 词频排序
	h.sortByFrequency(results)

	// 更新词频（异步）
	go h.updateFrequency(word)

	// 构建响应数据
	data := gin.H{
		"results": results,
	}

	// 缓存结果
	h.cache.Set(cacheKey, data, SearchCacheTTL)

	response.Success(c, data)
}

// Suggest 搜索建议
// GET /api/v1/search/suggest?q=xxx
func (h *SearchHandler) Suggest(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.BadRequest(c, "q parameter is required")
		return
	}

	limitStr := c.DefaultQuery("limit", strconv.Itoa(DefaultLimit))
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 50 {
		limit = DefaultLimit
	}

	// 检查缓存
	cacheKey := "suggest:" + q
	if cached, ok := h.cache.Get(cacheKey); ok {
		response.Success(c, cached)
		return
	}

	results := h.manager.Suggest(q, limit)

	// 构建响应数据
	data := gin.H{
		"query":       q,
		"suggestions": results,
	}

	// 缓存结果
	h.cache.Set(cacheKey, data, SuggestCacheTTL)

	response.Success(c, data)
}

// sortByFrequency 按词频排序搜索结果
func (h *SearchHandler) sortByFrequency(results []mdx.SearchResult) {
	if len(results) <= 1 {
		return
	}

	// 提取所有单词
	words := make([]string, len(results))
	for i, r := range results {
		words[i] = r.Word
	}

	// 批量查询词频
	var frequencies []model.WordFrequency
	h.db.Where("word IN ?", words).Find(&frequencies)

	freqMap := make(map[string]int64)
	for _, f := range frequencies {
		freqMap[f.Word] = f.SearchCount
	}

	// 按词频降序排序
	sort.Slice(results, func(i, j int) bool {
		return freqMap[results[i].Word] > freqMap[results[j].Word]
	})
}

// updateFrequency 异步更新词频
func (h *SearchHandler) updateFrequency(word string) {
	h.db.Exec(`
		INSERT INTO word_frequencies (word, search_count, last_searched, created_at, updated_at)
		VALUES (?, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT(word) DO UPDATE SET
			search_count = search_count + 1,
			last_searched = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
	`, word)
}

// rewriteResourceURLs 重写相对路径资源 URL
// 将 src="style.css" 或 href="image.png" 重写为 /api/v1/resources/{dictID}/path
func rewriteResourceURLs(definition string, dictID uint) string {
	// 匹配 src="..." 和 href="..." 中的相对路径
	// 排除以 / 或 http:// 或 https:// 开头的路径
	re := regexp.MustCompile(`(src|href)=["']([^"':/][^"']*)["']`)
	return re.ReplaceAllStringFunc(definition, func(match string) string {
		// 提取属性和路径
		parts := re.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}
		attr, path := parts[1], parts[2]
		return attr + `="/api/v1/resources/` + strconv.FormatUint(uint64(dictID), 10) + `/` + path + `"`
	})
}
