// 通用 API 响应包装
export interface ApiResponse<T> {
  code: number
  message: string
  data?: T
}

// 搜索结果项（与后端 SearchResult 对应）
export interface SearchResult {
  dict_id: number
  dict_name: string
  dict_title: string
  word: string
  definition: string // HTML 内容
}

// 搜索 API 响应
export interface SearchResponse {
  results: SearchResult[]
}

// 搜索建议项（与后端 SuggestResult 对应）
export interface SuggestResult {
  word: string
  dict_id: number
  dict_title: string
}

// 建议 API 响应
export interface SuggestResponse {
  query: string
  suggestions: SuggestResult[]
}

// 按字典分组的搜索结果
export interface GroupedSearchResults {
  [dictId: number]: {
    dictName: string
    dictTitle: string
    results: SearchResult[]
  }
}
