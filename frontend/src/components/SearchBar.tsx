import { useRef, useEffect } from 'react'
import { Input, Button, Card } from '@heroui/react'
import { useSearchStore } from '@/stores'
import { useSuggest } from '@/hooks'

function SearchIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
      />
    </svg>
  )
}

export function SearchBar() {
  const containerRef = useRef<HTMLDivElement>(null)
  const { keyword, showSuggestions, setKeyword, submitSearch, setShowSuggestions } =
    useSearchStore()
  const { data: suggestData } = useSuggest(keyword)

  const suggestions = suggestData?.suggestions ?? []

  // 点击外部关闭建议下拉
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setShowSuggestions(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [setShowSuggestions])

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      submitSearch()
    } else if (e.key === 'Escape') {
      setShowSuggestions(false)
    }
  }

  const handleSelectSuggestion = (word: string) => {
    setKeyword(word)
    submitSearch(word)
  }

  return (
    <div ref={containerRef} className="relative w-full max-w-2xl">
      <div className="flex flex-col sm:flex-row gap-2">
        <Input
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onKeyDown={handleKeyDown}
          onFocus={() => setShowSuggestions(true)}
          placeholder="搜索单词..."
          startContent={<SearchIcon className="w-5 h-5 text-default-400" />}
          size="lg"
          variant="bordered"
          classNames={{
            inputWrapper: 'bg-background',
          }}
        />
        <Button
          color="primary"
          size="lg"
          onPress={() => submitSearch()}
          className="px-8 w-full sm:w-auto"
        >
          搜索
        </Button>
      </div>

      {showSuggestions && suggestions.length > 0 && (
        <div className="absolute top-full left-0 right-0 z-50 mt-1">
          <Card className="py-1 shadow-lg">
            {suggestions.map((item, idx) => (
              <button
                key={`${item.dict_id}-${item.word}-${idx}`}
                onClick={() => handleSelectSuggestion(item.word)}
                className="w-full px-4 py-2 text-left hover:bg-default-100 dark:hover:bg-default-50 transition-colors flex items-center justify-between"
              >
                <span className="font-medium">{item.word}</span>
                <span className="text-xs text-default-400">{item.dict_title}</span>
              </button>
            ))}
          </Card>
        </div>
      )}
    </div>
  )
}
