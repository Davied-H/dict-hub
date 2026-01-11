import { Skeleton, Card, CardBody } from '@heroui/react'
import { DictCard } from './DictCard'
import type { SearchResult, GroupedSearchResults } from '@/types'

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
  if (isLoading) {
    return <SearchSkeleton />
  }

  if (error) {
    return <ErrorDisplay message={error.message || '搜索时发生错误'} />
  }

  if (results.length === 0) {
    return <EmptyResults keyword={keyword} />
  }

  const grouped = groupResults(results)

  return (
    <div className="space-y-6">
      {/* 结果统计 */}
      <div className="flex items-center justify-between text-sm text-default-500">
        <span>
          找到 <span className="font-semibold text-primary-500">{results.length}</span> 条结果
        </span>
      </div>
      
      {/* 字典卡片列表 */}
      {Object.entries(grouped).map(([dictId, group]) => (
        <DictCard
          key={dictId}
          dictTitle={group.dictTitle}
          results={group.results}
        />
      ))}
    </div>
  )
}
