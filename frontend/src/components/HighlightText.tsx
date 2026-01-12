import { useMemo } from 'react'

interface HighlightTextProps {
  text: string
  keyword: string
  className?: string
  highlightClassName?: string
}

/**
 * 高亮文本中的关键词
 * 支持大小写不敏感匹配
 */
export function HighlightText({
  text,
  keyword,
  className = '',
  highlightClassName = 'bg-warning-200/60 dark:bg-warning-700/40 text-warning-900 dark:text-warning-100 rounded px-0.5',
}: HighlightTextProps) {
  const parts = useMemo(() => {
    if (!keyword.trim()) {
      return [{ text, highlight: false }]
    }

    const regex = new RegExp(`(${escapeRegExp(keyword)})`, 'gi')
    const splitParts = text.split(regex)

    return splitParts
      .filter((part) => part !== '')
      .map((part) => ({
        text: part,
        highlight: part.toLowerCase() === keyword.toLowerCase(),
      }))
  }, [text, keyword])

  return (
    <span className={className}>
      {parts.map((part, index) =>
        part.highlight ? (
          <mark key={index} className={highlightClassName}>
            {part.text}
          </mark>
        ) : (
          <span key={index}>{part.text}</span>
        )
      )}
    </span>
  )
}

// 转义正则特殊字符
function escapeRegExp(string: string): string {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}
