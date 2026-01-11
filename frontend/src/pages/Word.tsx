import { useParams, useNavigate } from 'react-router-dom'
import { Button } from '@heroui/react'
import { useSearch } from '@/hooks'
import { SearchResults } from '@/components'

function BackIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
    </svg>
  )
}

export default function Word() {
  const { word } = useParams<{ word: string }>()
  const navigate = useNavigate()
  const decodedWord = word ? decodeURIComponent(word) : ''
  
  const { data, isLoading, error } = useSearch(decodedWord)

  const handleBack = () => {
    // 如果有历史记录则返回，否则导航到首页
    if (window.history.length > 1) {
      navigate(-1)
    } else {
      navigate('/')
    }
  }

  return (
    <div className="flex flex-col min-h-[calc(100vh-4rem)]">
      {/* 顶部导航栏 */}
      <div className="sticky top-0 z-10 bg-background/80 backdrop-blur-md border-b border-default-200 dark:border-default-100 -mx-4 px-4 py-3 mb-4">
        <div className="flex items-center gap-3">
          <Button
            isIconOnly
            variant="light"
            size="sm"
            onPress={handleBack}
            aria-label="返回"
            className="shrink-0"
          >
            <BackIcon className="w-5 h-5" />
          </Button>
          <div className="flex-1 min-w-0">
            <h1 className="text-lg font-bold text-default-800 dark:text-default-200 truncate">
              {decodedWord}
            </h1>
          </div>
        </div>
      </div>

      {/* 搜索结果 */}
      <div className="flex-1 pb-8">
        <SearchResults
          results={data?.results ?? []}
          keyword={decodedWord}
          isLoading={isLoading}
          error={error}
        />
      </div>
    </div>
  )
}
