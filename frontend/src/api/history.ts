import { apiClient } from './client'
import type { PagedResponse, SearchHistory, ApiResponse } from '@/types'

export const historyApi = {
  list: (page = 1, pageSize = 20) =>
    apiClient.get<PagedResponse<SearchHistory[]>>('/history', {
      params: { page, page_size: pageSize },
    }),

  clear: () => apiClient.delete<ApiResponse<null>>('/history'),

  exportUrl: '/api/v1/history/export',
}
