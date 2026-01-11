import { Button, Chip } from '@heroui/react'
import { useSearchStore } from '@/stores'
import { useSearch } from '@/hooks'
import { SearchBar, SearchResults } from '@/components'

function BookStackIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={1.5}
        d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
      />
    </svg>
  )
}

export default function Home() {
  const {
    submittedKeyword,
    recentSearches,
    setKeyword,
    submitSearch,
    clearRecent,
  } = useSearchStore()
  const { data, isLoading, error } = useSearch(submittedKeyword)

  const hasResults = submittedKeyword && (data?.results?.length ?? 0) > 0

  return (
    <div
      className={`flex flex-col items-center px-4 transition-all duration-500 ease-out ${
        hasResults ? 'pt-6' : 'pt-[15vh]'
      }`}
    >
      {/* Logo/品牌 */}
      <div
        className={`flex items-center gap-3 mb-6 transition-all duration-500 ${
          hasResults ? 'scale-90' : 'scale-100'
        }`}
      >
        <div
          className={`p-3 bg-gradient-to-br from-primary-500 to-secondary-500 rounded-2xl shadow-lg transition-all duration-500 ${
            hasResults ? 'p-2' : 'p-3'
          }`}
        >
          <BookStackIcon
            className={`text-white transition-all duration-500 ${
              hasResults ? 'w-6 h-6' : 'w-8 h-8'
            }`}
          />
        </div>
        <div>
          <h1
            className={`font-bold bg-gradient-to-r from-primary-600 to-secondary-600 dark:from-primary-400 dark:to-secondary-400 bg-clip-text text-transparent transition-all duration-500 ${
              hasResults ? 'text-2xl' : 'text-4xl'
            }`}
          >
            Dict Hub
          </h1>
          {!hasResults && (
            <p className="text-sm text-default-500 mt-1">
              智能词典，助你高效学习英语
            </p>
          )}
        </div>
      </div>

      {/* 搜索栏 */}
      <SearchBar />

      {/* 最近搜索 - 仅在无结果时显示 */}
      {!hasResults && recentSearches.length > 0 && (
        <div className="mt-8 w-full max-w-2xl animate-fade-in">
          <div className="flex justify-between items-center mb-3">
            <span className="text-sm font-medium text-default-600 dark:text-default-400">
              最近搜索
            </span>
            <Button
              size="sm"
              variant="light"
              className="text-default-400 hover:text-default-600"
              onPress={clearRecent}
            >
              清除
            </Button>
          </div>
          <div className="flex flex-wrap gap-2">
            {recentSearches.slice(0, 8).map((word) => (
              <Chip
                key={word}
                variant="flat"
                className="cursor-pointer hover:bg-primary-100 dark:hover:bg-primary-900/30 transition-colors"
                onClick={() => {
                  setKeyword(word)
                  submitSearch(word)
                }}
              >
                {word}
              </Chip>
            ))}
          </div>
        </div>
      )}

      {/* 功能提示 - 仅在首页无搜索时显示 */}
      {!submittedKeyword && recentSearches.length === 0 && (
        <div className="mt-12 grid grid-cols-1 md:grid-cols-3 gap-6 w-full max-w-3xl animate-fade-in-delay">
          <FeatureCard
            icon="📚"
            title="多字典支持"
            description="同时查询多本权威词典"
          />
          <FeatureCard
            icon="🎯"
            title="精准释义"
            description="专业词条，例句丰富"
          />
          <FeatureCard
            icon="🔊"
            title="发音播放"
            description="原生发音，纠正口语"
          />
        </div>
      )}

      {/* 搜索结果区域 */}
      {submittedKeyword && (
        <div className="w-full max-w-4xl mt-6 pb-16 animate-fade-in">
          <SearchResults
            results={data?.results ?? []}
            keyword={submittedKeyword}
            isLoading={isLoading}
            error={error}
          />
        </div>
      )}
    </div>
  )
}

function FeatureCard({
  icon,
  title,
  description,
}: {
  icon: string
  title: string
  description: string
}) {
  return (
    <div className="flex flex-col items-center p-6 rounded-2xl bg-default-50 dark:bg-default-100/5 border border-default-200 dark:border-default-100 hover:border-primary-300 dark:hover:border-primary-700 transition-colors">
      <span className="text-3xl mb-3">{icon}</span>
      <h3 className="font-semibold text-default-700 dark:text-default-300">
        {title}
      </h3>
      <p className="text-sm text-default-500 dark:text-default-400 text-center mt-1">
        {description}
      </p>
    </div>
  )
}
