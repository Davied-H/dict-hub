export interface Word {
  id: number
  word: string
  phonetic?: string
  definition: string
  example?: string
  partOfSpeech?: string
  synonyms?: string
  antonyms?: string
  difficulty: number
  createdAt: string
  updatedAt: string
}

export interface SearchHistory {
  id: number
  word: string
  sessionId: string
  found: boolean
  responseTime: number
  createdAt: string
}

export interface WordFrequency {
  id: number
  word: string
  searchCount: number
  viewCount: number
  lastSearched: string
}

export interface DictSource {
  id: number
  name: string
  title: string
  description: string
  path: string
  enabled: boolean
  sort_order: number
  word_count: number
  has_mdd: boolean
  file_size: number
  created_at: string
  updated_at: string
}

// 通用 API 响应包装
export interface ApiResponse<T> {
  code: number
  message: string
  data?: T
}

export interface PagedResponse<T> {
  code: number
  message: string
  data?: T
  meta?: {
    total: number
    page: number
    pageSize: number
  }
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
