import { useState, useMemo } from 'react'
import { Card, CardHeader, CardBody, Divider, Button, Chip } from '@heroui/react'
import { AudioPlayer } from './AudioPlayer'
import { DictContent } from './DictContent'
import { HighlightText } from './HighlightText'
import type { SearchResult } from '@/types'

// 词典主题色配置 - 根据 dict_id 轮换颜色
const DICT_THEMES = [
  { 
    name: 'blue',
    border: 'border-l-blue-500',
    headerBg: 'from-blue-50 to-sky-50 dark:from-blue-900/30 dark:to-sky-900/30',
    iconBg: 'bg-blue-100 dark:bg-blue-900/50',
    iconColor: 'text-blue-600 dark:text-blue-400',
  },
  { 
    name: 'emerald',
    border: 'border-l-emerald-500',
    headerBg: 'from-emerald-50 to-teal-50 dark:from-emerald-900/30 dark:to-teal-900/30',
    iconBg: 'bg-emerald-100 dark:bg-emerald-900/50',
    iconColor: 'text-emerald-600 dark:text-emerald-400',
  },
  { 
    name: 'violet',
    border: 'border-l-violet-500',
    headerBg: 'from-violet-50 to-purple-50 dark:from-violet-900/30 dark:to-purple-900/30',
    iconBg: 'bg-violet-100 dark:bg-violet-900/50',
    iconColor: 'text-violet-600 dark:text-violet-400',
  },
  { 
    name: 'amber',
    border: 'border-l-amber-500',
    headerBg: 'from-amber-50 to-orange-50 dark:from-amber-900/30 dark:to-orange-900/30',
    iconBg: 'bg-amber-100 dark:bg-amber-900/50',
    iconColor: 'text-amber-600 dark:text-amber-400',
  },
  { 
    name: 'rose',
    border: 'border-l-rose-500',
    headerBg: 'from-rose-50 to-pink-50 dark:from-rose-900/30 dark:to-pink-900/30',
    iconBg: 'bg-rose-100 dark:bg-rose-900/50',
    iconColor: 'text-rose-600 dark:text-rose-400',
  },
]

function getThemeByDictId(dictId: number) {
  return DICT_THEMES[dictId % DICT_THEMES.length]
}

function BookIcon({ className }: { className?: string }) {
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
        strokeWidth={1.5}
        d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
      />
    </svg>
  )
}

function ChevronDownIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
    </svg>
  )
}

function ChevronUpIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
    </svg>
  )
}

// 从 HTML 中提取纯文本预览
function extractTextPreview(html: string, maxLength = 120): string {
  const tempDiv = document.createElement('div')
  tempDiv.innerHTML = html
  const text = tempDiv.textContent || tempDiv.innerText || ''
  const cleaned = text.replace(/\s+/g, ' ').trim()
  if (cleaned.length <= maxLength) return cleaned
  return cleaned.slice(0, maxLength).trim() + '...'
}

// 尝试从 HTML 定义中提取音标
function extractPhonetic(html: string): string | null {
  const phoneticMatch = html.match(/[\/\[][^\]\/]{1,40}[\/\]]/)
  return phoneticMatch ? phoneticMatch[0] : null
}

interface DictCardProps {
  dictTitle: string
  results: SearchResult[]
  keyword?: string
  isExactSection?: boolean
}

export function DictCard({ dictTitle, results, keyword = '', isExactSection = false }: DictCardProps) {
  // 获取词典主题色 (使用第一个结果的 dict_id)
  const dictId = results[0]?.dict_id ?? 0
  const theme = getThemeByDictId(dictId)
  // 每个词条的展开状态
  const [expandedEntries, setExpandedEntries] = useState<Set<number>>(() => {
    // 默认展开第一条精确匹配
    if (isExactSection && results.length > 0) {
      return new Set([0])
    }
    return new Set()
  })

  const toggleEntry = (index: number) => {
    setExpandedEntries((prev) => {
      const next = new Set(prev)
      if (next.has(index)) {
        next.delete(index)
      } else {
        next.add(index)
      }
      return next
    })
  }

  const expandAll = () => {
    setExpandedEntries(new Set(results.map((_, i) => i)))
  }

  const collapseAll = () => {
    setExpandedEntries(new Set())
  }

  const allExpanded = expandedEntries.size === results.length
  const anyExpanded = expandedEntries.size > 0

  return (
    <Card 
      className={`
        dict-card overflow-hidden shadow-sm border-l-4
        ${theme.border}
        ${isExactSection 
          ? 'border-y border-r border-primary-200 dark:border-primary-800/50 ring-1 ring-primary-100 dark:ring-primary-900/30' 
          : 'border-y border-r border-default-200 dark:border-default-100'
        }
      `}
    >
      {/* 字典标题栏 */}
      <CardHeader 
        className={`
          dict-card-header py-2.5 sm:py-3 px-3 sm:px-5
          bg-gradient-to-r ${theme.headerBg}
        `}
      >
        <div className="flex items-center justify-between w-full">
          <div className="flex items-center gap-2 sm:gap-3">
            <div className={`p-1.5 sm:p-2 rounded-lg ${theme.iconBg}`}>
              <BookIcon className={`w-4 h-4 sm:w-5 sm:h-5 ${theme.iconColor}`} />
            </div>
            <div className="flex items-center gap-2">
              <h2 className="font-semibold text-sm sm:text-base text-default-800 dark:text-default-200">
                {dictTitle}
              </h2>
              <Chip size="sm" variant="flat" className="text-xs">
                {results.length} 条
              </Chip>
            </div>
          </div>
          
          {/* 展开/收起全部按钮 */}
          {results.length > 1 && (
            <Button
              size="sm"
              variant="light"
              onPress={allExpanded ? collapseAll : expandAll}
              className="text-xs text-default-500"
            >
              {allExpanded ? '收起全部' : '展开全部'}
            </Button>
          )}
        </div>
      </CardHeader>

      <Divider />

      {/* 词条内容区域 */}
      <CardBody className="p-0">
        {results.map((result, idx) => {
          const isExpanded = expandedEntries.has(idx)
          const preview = useMemo(() => extractTextPreview(result.definition), [result.definition])
          const phonetic = useMemo(() => extractPhonetic(result.definition), [result.definition])

          return (
            <article
              key={idx}
              className={`
                dict-entry transition-colors duration-200
                ${idx !== results.length - 1
                  ? 'border-b border-default-100 dark:border-default-50'
                  : ''
                }
                ${isExpanded 
                  ? 'bg-default-50/30 dark:bg-default-100/5' 
                  : ''
                }
              `}
            >
              {/* 词条标题行 - 可点击展开/收起 */}
              <header 
                className="flex items-start justify-between gap-2 sm:gap-4 p-3 sm:p-5 cursor-pointer hover:bg-default-50 dark:hover:bg-default-100/5 transition-colors"
                onClick={() => toggleEntry(idx)}
              >
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 sm:gap-3 flex-wrap">
                    <h3 className="text-xl sm:text-2xl font-bold text-primary-600 dark:text-primary-400 tracking-tight">
                      {keyword ? (
                        <HighlightText text={result.word} keyword={keyword} />
                      ) : (
                        result.word
                      )}
                    </h3>
                    {phonetic && (
                      <span className="text-sm text-default-400 dark:text-default-500 font-mono">
                        {phonetic}
                      </span>
                    )}
                    <AudioPlayer definition={result.definition} word={result.word} />
                  </div>
                  
                  {/* 摘要预览 - 未展开时显示 */}
                  {!isExpanded && (
                    <p className="mt-2 text-sm text-default-500 dark:text-default-400 line-clamp-2">
                      {keyword ? (
                        <HighlightText text={preview} keyword={keyword} />
                      ) : (
                        preview
                      )}
                    </p>
                  )}
                </div>
                
                {/* 展开/收起指示器 */}
                <button
                  className="flex-shrink-0 p-1.5 rounded-full hover:bg-default-100 dark:hover:bg-default-100/10 transition-colors"
                  onClick={(e) => {
                    e.stopPropagation()
                    toggleEntry(idx)
                  }}
                >
                  {isExpanded ? (
                    <ChevronUpIcon className="w-5 h-5 text-default-400" />
                  ) : (
                    <ChevronDownIcon className="w-5 h-5 text-default-400" />
                  )}
                </button>
              </header>

              {/* 完整释义内容 - 展开时显示 */}
              {isExpanded && (
                <div className="px-3 sm:px-5 pb-3 sm:pb-5 animate-expand">
                  <div className="dict-definition pt-2 border-t border-default-100 dark:border-default-50">
                    <DictContent html={result.definition} dictId={result.dict_id} />
                  </div>
                </div>
              )}
            </article>
          )
        })}
      </CardBody>
    </Card>
  )
}
