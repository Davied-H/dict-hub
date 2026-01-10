import { apiClient } from './client'
import type { Word, ApiResponse, PagedResponse } from '@/types'

export const dictionaryApi = {
  getWords: async (page = 1, pageSize = 20) => {
    const response = await apiClient.get<PagedResponse<Word[]>>('/words', {
      params: { page, page_size: pageSize },
    })
    return response.data
  },

  getWord: async (id: number) => {
    const response = await apiClient.get<ApiResponse<Word>>(`/words/${id}`)
    return response.data
  },

  search: async (keyword: string) => {
    const response = await apiClient.get<ApiResponse<Word[]>>('/search', {
      params: { q: keyword },
    })
    return response.data
  },

  createWord: async (word: Omit<Word, 'id' | 'createdAt' | 'updatedAt'>) => {
    const response = await apiClient.post<ApiResponse<Word>>('/words', word)
    return response.data
  },
}
