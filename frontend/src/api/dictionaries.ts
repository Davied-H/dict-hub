import { apiClient } from './client'
import type { ApiResponse, DictSource } from '@/types'

export const dictionariesApi = {
  list: () => apiClient.get<ApiResponse<DictSource[]>>('/dictionaries'),

  add: (path: string) =>
    apiClient.post<ApiResponse<DictSource>>('/dictionaries', { path }),

  toggle: (id: number) =>
    apiClient.put<ApiResponse<DictSource>>(`/dictionaries/${id}/toggle`),

  reorder: (orders: { id: number; sort_order: number }[]) =>
    apiClient.put<ApiResponse<null>>('/dictionaries/reorder', orders),

  delete: (id: number) => apiClient.delete<ApiResponse<null>>(`/dictionaries/${id}`),
}
