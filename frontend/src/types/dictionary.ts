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
