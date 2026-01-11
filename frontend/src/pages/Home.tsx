import { Button, Chip } from '@heroui/react'
import { useSearchStore } from '@/stores'
import { useSearch } from '@/hooks'
import { ThemeToggle, SearchBar, SearchResults } from '@/components'

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
    <div className="min-h-screen bg-background transition-colors duration-300">
      {/* 主题切换按钮 - 固定右上角 */}
      <div className="fixed top-4 right-4 z-50">
        <ThemeToggle />
      </div>

      {/* 主内容区 */}
      <main
        className={`flex flex-col items-center px-4 transition-all duration-500 ${
          hasResults ? 'pt-8' : 'pt-[30vh]'
        }`}
      >
        {/* Logo/品牌 */}
        <h1
          className={`font-bold text-default-800 transition-all duration-500 mb-8 ${
            hasResults ? 'text-2xl' : 'text-4xl'
          }`}
        >
          Dict Hub
        </h1>

        {/* 搜索栏 */}
        <SearchBar />

        {/* 最近搜索 - 仅在无结果时显示 */}
        {!hasResults && recentSearches.length > 0 && (
          <div className="mt-8 w-full max-w-2xl">
            <div className="flex justify-between items-center mb-3">
              <span className="text-sm text-default-500">最近搜索</span>
              <Button size="sm" variant="light" onPress={clearRecent}>
                清除
              </Button>
            </div>
            <div className="flex flex-wrap gap-2">
              {recentSearches.slice(0, 8).map((word) => (
                <Chip
                  key={word}
                  variant="flat"
                  className="cursor-pointer"
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

        {/* 搜索结果区域 */}
        {submittedKeyword && (
          <div className="w-full max-w-4xl mt-8 pb-16">
            <SearchResults
              results={data?.results ?? []}
              keyword={submittedKeyword}
              isLoading={isLoading}
              error={error}
            />
          </div>
        )}
      </main>
    </div>
  )
}
