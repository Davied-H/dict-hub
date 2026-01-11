import { useQuery } from '@tanstack/react-query'
import { searchApi } from '@/api'

export function useSearch(word: string) {
  return useQuery({
    queryKey: ['search', word],
    queryFn: () => searchApi.search(word),
    enabled: word.trim().length > 0,
    staleTime: 1000 * 60 * 5, // 5 分钟缓存
  })
}
