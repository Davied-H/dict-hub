import { apiClient } from './client'
import type { ApiResponse, SearchResponse, SuggestResponse } from '@/types'

export const searchApi = {
  // 跨字典搜索
  search: async (word: string): Promise<SearchResponse> => {
    const response = await apiClient.get<ApiResponse<SearchResponse>>(
      '/search',
      { params: { word } }
    )
    return response.data.data!
  },

  // 搜索建议（自动补全）
  suggest: async (q: string, limit = 10): Promise<SuggestResponse> => {
    const response = await apiClient.get<ApiResponse<SuggestResponse>>(
      '/search/suggest',
      { params: { q, limit } }
    )
    return response.data.data!
  },
}
