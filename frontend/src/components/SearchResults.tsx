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
    <div className="space-y-6">
      {[1, 2].map((i) => (
        <Card key={i}>
          <CardBody className="space-y-3">
            <Skeleton className="w-1/4 h-6 rounded-lg" />
            <Skeleton className="w-full h-4 rounded-lg" />
            <Skeleton className="w-3/4 h-4 rounded-lg" />
            <Skeleton className="w-1/2 h-4 rounded-lg" />
          </CardBody>
        </Card>
      ))}
    </div>
  )
}

function ErrorDisplay({ message }: { message: string }) {
  return (
    <Card className="border-danger-200 bg-danger-50 dark:bg-danger-900/20">
      <CardBody className="text-center py-8">
        <p className="text-danger">{message}</p>
      </CardBody>
    </Card>
  )
}

function EmptyResults({ keyword }: { keyword: string }) {
  return (
    <Card>
      <CardBody className="text-center py-12">
        <p className="text-default-500">
          未找到 "<span className="font-medium">{keyword}</span>" 的相关结果
        </p>
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
