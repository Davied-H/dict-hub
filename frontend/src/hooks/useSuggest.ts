import { useQuery } from '@tanstack/react-query'
import { searchApi } from '@/api'
import { useDebounce } from './useDebounce'

export function useSuggest(query: string) {
  const debouncedQuery = useDebounce(query, 300)

  return useQuery({
    queryKey: ['suggest', debouncedQuery],
    queryFn: () => searchApi.suggest(debouncedQuery),
    enabled: debouncedQuery.trim().length >= 1,
    staleTime: 1000 * 60, // 1 分钟缓存
  })
}
