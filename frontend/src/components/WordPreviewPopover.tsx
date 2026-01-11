import { useState, useCallback } from 'react'
import { Popover, PopoverTrigger, PopoverContent, Spinner } from '@heroui/react'
import { useQuery } from '@tanstack/react-query'
import { searchApi } from '@/api'

interface WordPreviewPopoverProps {
  word: string
  children: React.ReactNode
  onWordClick: (word: string) => void
}

// 从 HTML 中提取纯文本
function extractTextFromHtml(html: string): string {
  const tempDiv = document.createElement('div')
  tempDiv.innerHTML = html
  return tempDiv.textContent || tempDiv.innerText || ''
}

// 从定义中提取简短释义
function extractBriefDefinition(definition: string): string {
  const text = extractTextFromHtml(definition)
  // 取前100个字符，避免过长
  const brief = text.slice(0, 150)
  return brief.length < text.length ? brief + '...' : brief
}

// 从定义中尝试提取音标
function extractPhonetic(definition: string): string | null {
  // 尝试匹配常见音标格式 /.../ 或 [...]
  const phoneticMatch = definition.match(/[\/\[][^\]\/]+[\/\]]/)
  return phoneticMatch ? phoneticMatch[0] : null
}

export function WordPreviewPopover({ word, children, onWordClick }: WordPreviewPopoverProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [shouldFetch, setShouldFetch] = useState(false)

  const { data, isLoading, error } = useQuery({
    queryKey: ['preview', word],
    queryFn: () => searchApi.search(word),
    enabled: shouldFetch && word.trim().length > 0,
    staleTime: 1000 * 60 * 10, // 10 分钟缓存
  })

  const handleOpenChange = useCallback((open: boolean) => {
    setIsOpen(open)
    if (open && !shouldFetch) {
      setShouldFetch(true)
    }
  }, [shouldFetch])

  const handleClick = useCallback((e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsOpen(false)
    onWordClick(word)
  }, [word, onWordClick])

  // 获取第一个结果的简短信息
  const firstResult = data?.results?.[0]
  const phonetic = firstResult ? extractPhonetic(firstResult.definition) : null
  const briefDef = firstResult ? extractBriefDefinition(firstResult.definition) : null

  return (
    <Popover
      isOpen={isOpen}
      onOpenChange={handleOpenChange}
      placement="top"
      showArrow
      offset={10}
      triggerType="grid"
    >
      <PopoverTrigger>
        <span
          className="word-link text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 cursor-pointer underline decoration-dotted underline-offset-2 hover:decoration-solid transition-colors"
          onClick={handleClick}
        >
          {children}
        </span>
      </PopoverTrigger>
      <PopoverContent>
        <div className="p-3 max-w-xs min-w-[200px]">
          {isLoading ? (
            <div className="flex items-center justify-center py-2">
              <Spinner size="sm" />
              <span className="ml-2 text-sm text-default-500">加载中...</span>
            </div>
          ) : error ? (
            <div className="text-sm text-danger-500">加载失败</div>
          ) : !firstResult ? (
            <div className="text-sm text-default-500">未找到释义</div>
          ) : (
            <>
              <div className="flex items-center gap-2 mb-2">
                <h4 className="font-bold text-base text-default-800 dark:text-default-200">
                  {word}
                </h4>
                {phonetic && (
                  <span className="text-sm text-default-500">{phonetic}</span>
                )}
              </div>
              <p className="text-sm text-default-600 dark:text-default-400 leading-relaxed">
                {briefDef}
              </p>
              <div className="mt-2 pt-2 border-t border-default-200 dark:border-default-100">
                <button
                  onClick={handleClick}
                  className="text-xs text-primary-500 hover:text-primary-600 dark:hover:text-primary-400 transition-colors"
                >
                  点击查看完整释义 →
                </button>
              </div>
            </>
          )}
        </div>
      </PopoverContent>
    </Popover>
  )
}
