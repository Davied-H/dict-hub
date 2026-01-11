import { useRef, useEffect, useCallback } from 'react'
import { Input, Button, Card, Kbd } from '@heroui/react'
import { useSearchStore } from '@/stores'
import { useSuggest } from '@/hooks'
import { HighlightText } from './HighlightText'

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

function StarIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="currentColor" viewBox="0 0 20 20">
      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
    </svg>
  )
}

export function SearchBar() {
  const containerRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)
  const { 
    keyword, 
    showSuggestions, 
    selectedSuggestionIndex,
    setKeyword, 
    submitSearch, 
    setShowSuggestions,
    setSelectedSuggestionIndex,
  } = useSearchStore()
  const { data: suggestData, isLoading: isSuggestLoading } = useSuggest(keyword)

  const suggestions = suggestData?.suggestions ?? []
  
  // 按精确匹配排序：精确匹配在前
  const sortedSuggestions = [...suggestions].sort((a, b) => {
    const aExact = a.word.toLowerCase() === keyword.toLowerCase()
    const bExact = b.word.toLowerCase() === keyword.toLowerCase()
    if (aExact && !bExact) return -1
    if (!aExact && bExact) return 1
    return 0
  })

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

  // 全局快捷键 / 聚焦搜索框
  useEffect(() => {
    const handleGlobalKeyDown = (e: KeyboardEvent) => {
      // 如果在输入框中，不处理
      if (
        e.target instanceof HTMLInputElement ||
        e.target instanceof HTMLTextAreaElement
      ) {
        return
      }
      if (e.key === '/') {
        e.preventDefault()
        inputRef.current?.focus()
      }
    }
    document.addEventListener('keydown', handleGlobalKeyDown)
    return () => document.removeEventListener('keydown', handleGlobalKeyDown)
  }, [])

  const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
    const suggestionsCount = sortedSuggestions.length

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault()
        if (showSuggestions && suggestionsCount > 0) {
          setSelectedSuggestionIndex(
            selectedSuggestionIndex < suggestionsCount - 1
              ? selectedSuggestionIndex + 1
              : 0
          )
        } else if (suggestionsCount > 0) {
          setShowSuggestions(true)
          setSelectedSuggestionIndex(0)
        }
        break
      case 'ArrowUp':
        e.preventDefault()
        if (showSuggestions && suggestionsCount > 0) {
          setSelectedSuggestionIndex(
            selectedSuggestionIndex > 0
              ? selectedSuggestionIndex - 1
              : suggestionsCount - 1
          )
        }
        break
      case 'Enter':
        e.preventDefault()
        if (selectedSuggestionIndex >= 0 && selectedSuggestionIndex < suggestionsCount) {
          handleSelectSuggestion(sortedSuggestions[selectedSuggestionIndex].word)
        } else {
          submitSearch()
        }
        break
      case 'Escape':
        setShowSuggestions(false)
        setSelectedSuggestionIndex(-1)
        break
    }
  }, [showSuggestions, selectedSuggestionIndex, sortedSuggestions, setShowSuggestions, setSelectedSuggestionIndex, submitSearch])

  const handleSelectSuggestion = (word: string) => {
    setKeyword(word)
    submitSearch(word)
  }

  const isExactMatch = (word: string) => word.toLowerCase() === keyword.toLowerCase()

  return (
    <div ref={containerRef} className="relative w-full max-w-2xl">
      <div className="flex flex-col sm:flex-row gap-2">
        <Input
          ref={inputRef}
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onKeyDown={handleKeyDown}
          onFocus={() => setShowSuggestions(true)}
          placeholder="搜索单词..."
          startContent={<SearchIcon className="w-5 h-5 text-default-400" />}
          endContent={
            <Kbd keys={['command']} className="hidden sm:flex text-xs">
              /
            </Kbd>
          }
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
          isLoading={isSuggestLoading}
        >
          搜索
        </Button>
      </div>

      {showSuggestions && sortedSuggestions.length > 0 && (
        <div className="absolute top-full left-0 right-0 z-50 mt-1">
          <Card className="py-1 shadow-lg overflow-hidden">
            <div className="max-h-80 overflow-y-auto">
              {sortedSuggestions.map((item, idx) => {
                const isExact = isExactMatch(item.word)
                const isSelected = idx === selectedSuggestionIndex

                return (
                  <button
                    key={`${item.dict_id}-${item.word}-${idx}`}
                    onClick={() => handleSelectSuggestion(item.word)}
                    onMouseEnter={() => setSelectedSuggestionIndex(idx)}
                    className={`
                      w-full px-4 py-2.5 text-left transition-colors
                      ${isSelected 
                        ? 'bg-primary-50 dark:bg-primary-900/30' 
                        : 'hover:bg-default-100 dark:hover:bg-default-50'
                      }
                      ${isExact ? 'border-l-3 border-l-primary-500' : ''}
                    `}
                  >
                    <div className="flex items-center justify-between gap-2">
                      <div className="flex items-center gap-2 min-w-0">
                        {isExact && (
                          <StarIcon className="w-4 h-4 text-warning-500 flex-shrink-0" />
                        )}
                        <HighlightText
                          text={item.word}
                          keyword={keyword}
                          className="font-semibold text-default-800 dark:text-default-200 truncate"
                        />
                        {isExact && (
                          <span className="text-xs px-1.5 py-0.5 rounded-full bg-primary-100 text-primary-700 dark:bg-primary-900/50 dark:text-primary-300 flex-shrink-0">
                            精确
                          </span>
                        )}
                      </div>
                      <span className="text-xs text-default-400 dark:text-default-500 truncate flex-shrink-0 ml-2">
                        {item.dict_title}
                      </span>
                    </div>
                  </button>
                )
              })}
            </div>
            {/* 键盘操作提示 */}
            <div className="hidden sm:flex items-center justify-between px-4 py-2 border-t border-default-100 dark:border-default-50 bg-default-50/50 dark:bg-default-100/5 text-xs text-default-400">
              <div className="flex items-center gap-3">
                <span className="flex items-center gap-1">
                  <Kbd className="text-xs">↑</Kbd>
                  <Kbd className="text-xs">↓</Kbd>
                  <span>导航</span>
                </span>
                <span className="flex items-center gap-1">
                  <Kbd className="text-xs">Enter</Kbd>
                  <span>选择</span>
                </span>
                <span className="flex items-center gap-1">
                  <Kbd className="text-xs">Esc</Kbd>
                  <span>关闭</span>
                </span>
              </div>
            </div>
          </Card>
        </div>
      )}
    </div>
  )
}
