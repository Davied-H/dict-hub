import { useMemo } from 'react'
import { Skeleton, Card, CardBody, Chip } from '@heroui/react'
import { useSearchStore } from '@/stores'
import { DictCard } from './DictCard'
import { DictFilterTabs } from './DictFilterTabs'
import type { SearchResult, GroupedSearchResults } from '@/types'

// 判断是否为精确匹配
function isExactMatch(word: string, keyword: string): boolean {
  return word.toLowerCase() === keyword.toLowerCase()
}

// 按字典分组搜索结果
function groupResults(results: SearchResult[]): GroupedSearchResults {
  return results.reduce((acc, result) => {
    if (!acc[result.dict_id]) {
      acc[result.dict_id] = {
        dictName: result.dict_name,
        dictTitle: result.dict_title,
        results: [],
      }
    }
    acc[result.dict_id].results.push(result)
    return acc
  }, {} as GroupedSearchResults)
}

// 分离精确匹配和相关词条
function separateResults(
  results: SearchResult[],
  keyword: string
): { exactMatches: SearchResult[]; relatedMatches: SearchResult[] } {
  const exactMatches: SearchResult[] = []
  const relatedMatches: SearchResult[] = []

  results.forEach((result) => {
    if (isExactMatch(result.word, keyword)) {
      exactMatches.push(result)
    } else {
      relatedMatches.push(result)
    }
  })

  return { exactMatches, relatedMatches }
}

function SearchSkeleton() {
  return (
    <div className="space-y-6 animate-pulse">
      {[1, 2].map((i) => (
        <Card key={i} className="overflow-hidden">
          <div className="bg-gradient-to-r from-primary-50 to-secondary-50 dark:from-primary-900/20 dark:to-secondary-900/20 p-4">
            <div className="flex items-center gap-3">
              <Skeleton className="w-10 h-10 rounded-lg" />
              <Skeleton className="w-48 h-5 rounded-lg" />
            </div>
          </div>
          <CardBody className="space-y-4 p-5">
            <div className="flex items-center gap-3">
              <Skeleton className="w-32 h-8 rounded-lg" />
              <Skeleton className="w-8 h-8 rounded-full" />
            </div>
            <Skeleton className="w-full h-4 rounded-lg" />
            <Skeleton className="w-4/5 h-4 rounded-lg" />
            <Skeleton className="w-3/5 h-4 rounded-lg" />
            <div className="bg-default-50 dark:bg-default-100/5 rounded-lg p-3">
              <Skeleton className="w-full h-4 rounded-lg" />
            </div>
          </CardBody>
        </Card>
      ))}
    </div>
  )
}

function SearchIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
    </svg>
  )
}

function ErrorIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
    </svg>
  )
}

function StarIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="currentColor" viewBox="0 0 20 20">
      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
    </svg>
  )
}

function LinkIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
    </svg>
  )
}

function ErrorDisplay({ message }: { message: string }) {
  return (
    <Card className="border-danger-200 dark:border-danger-800 bg-danger-50 dark:bg-danger-900/20">
      <CardBody className="text-center py-12">
        <div className="flex flex-col items-center gap-4">
          <div className="p-4 bg-danger-100 dark:bg-danger-900/30 rounded-full">
            <ErrorIcon className="w-8 h-8 text-danger-500" />
          </div>
          <div>
            <p className="text-danger-600 dark:text-danger-400 font-medium">出错了</p>
            <p className="text-danger-500 dark:text-danger-500 text-sm mt-1">{message}</p>
          </div>
        </div>
      </CardBody>
    </Card>
  )
}

function EmptyResults({ keyword }: { keyword: string }) {
  return (
    <Card className="border-default-200 dark:border-default-100">
      <CardBody className="text-center py-16">
        <div className="flex flex-col items-center gap-4">
          <div className="p-4 bg-default-100 dark:bg-default-100/10 rounded-full">
            <SearchIcon className="w-10 h-10 text-default-400" />
          </div>
          <div>
            <p className="text-lg font-medium text-default-600 dark:text-default-400">
              未找到结果
            </p>
            <p className="text-default-400 dark:text-default-500 mt-2">
              没有找到 "<span className="font-semibold text-primary-500">{keyword}</span>" 的相关释义
            </p>
          </div>
          <div className="mt-4 text-sm text-default-400 dark:text-default-500">
            <p>建议：</p>
            <ul className="mt-2 space-y-1">
              <li>• 检查拼写是否正确</li>
              <li>• 尝试搜索单词的原形</li>
              <li>• 使用更简单的词汇</li>
            </ul>
          </div>
        </div>
      </CardBody>
    </Card>
  )
}

// 结果分组区块标题
function ResultSectionHeader({ 
  title, 
  count, 
  icon: Icon,
  variant = 'exact',
}: { 
  title: string
  count: number
  icon: React.ComponentType<{ className?: string }>
  variant?: 'exact' | 'related'
}) {
  const colorClasses = variant === 'exact' 
    ? 'text-primary-600 dark:text-primary-400 bg-primary-50 dark:bg-primary-900/30'
    : 'text-secondary-600 dark:text-secondary-400 bg-secondary-50 dark:bg-secondary-900/30'

  return (
    <div className="flex items-center gap-2 mb-4">
      <div className={`flex items-center gap-2 px-3 py-1.5 rounded-full ${colorClasses}`}>
        <Icon className="w-4 h-4" />
        <span className="text-sm font-semibold">{title}</span>
      </div>
      <Chip size="sm" variant="flat" className="text-xs">
        {count} 条
      </Chip>
    </div>
  )
}

interface SearchResultsProps {
  results: SearchResult[]
  keyword: string
  isLoading: boolean
  error: Error | null
}

export function SearchResults({
  results,
  keyword,
  isLoading,
  error,
}: SearchResultsProps) {
  const { activeDictFilter } = useSearchStore()

  // 按词典筛选
  const filteredResults = useMemo(() => {
    if (activeDictFilter === null) return results
    return results.filter((r) => r.dict_id === activeDictFilter)
  }, [results, activeDictFilter])

  // 分离精确匹配和相关词条
  const { exactMatches, relatedMatches } = useMemo(
    () => separateResults(filteredResults, keyword),
    [filteredResults, keyword]
  )

  // 按词典分组
  const groupedExact = useMemo(() => groupResults(exactMatches), [exactMatches])
  const groupedRelated = useMemo(() => groupResults(relatedMatches), [relatedMatches])

  // 统计词典数量
  const dictCount = useMemo(() => {
    const dictIds = new Set(results.map((r) => r.dict_id))
    return dictIds.size
  }, [results])

  if (isLoading) {
    return <SearchSkeleton />
  }

  if (error) {
    return <ErrorDisplay message={error.message || '搜索时发生错误'} />
  }

  if (results.length === 0) {
    return <EmptyResults keyword={keyword} />
  }

  return (
    <div className="space-y-6">
      {/* 词典筛选 Tabs */}
      <DictFilterTabs results={results} />

      {/* 结果统计 */}
      <div className="flex items-center justify-between text-sm">
        <div className="flex items-center gap-2 text-default-500">
          <span>
            找到 <span className="font-semibold text-primary-500">{filteredResults.length}</span> 条结果
          </span>
          {activeDictFilter === null && dictCount > 1 && (
            <span className="text-default-400">
              来自 {dictCount} 本词典
            </span>
          )}
        </div>
        {exactMatches.length > 0 && (
          <Chip 
            size="sm" 
            variant="flat" 
            color="primary"
            startContent={<StarIcon className="w-3 h-3" />}
          >
            {exactMatches.length} 精确匹配
          </Chip>
        )}
      </div>
      
      {/* 精确匹配区块 */}
      {exactMatches.length > 0 && (
        <section>
          <ResultSectionHeader 
            title="精确匹配" 
            count={exactMatches.length}
            icon={StarIcon}
            variant="exact"
          />
          <div className="space-y-4">
            {Object.entries(groupedExact).map(([dictId, group]) => (
              <DictCard
                key={dictId}
                dictTitle={group.dictTitle}
                results={group.results}
                keyword={keyword}
                isExactSection
              />
            ))}
          </div>
        </section>
      )}

      {/* 相关词条区块 */}
      {relatedMatches.length > 0 && (
        <section>
          <ResultSectionHeader 
            title="相关词条" 
            count={relatedMatches.length}
            icon={LinkIcon}
            variant="related"
          />
          <div className="space-y-4">
            {Object.entries(groupedRelated).map(([dictId, group]) => (
              <DictCard
                key={dictId}
                dictTitle={group.dictTitle}
                results={group.results}
                keyword={keyword}
              />
            ))}
          </div>
        </section>
      )}
    </div>
  )
}
