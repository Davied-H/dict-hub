import { useEffect, useRef } from 'react'
import { Spinner } from '@heroui/react'
import { useQuery } from '@tanstack/react-query'
import { createPortal } from 'react-dom'
import { searchApi } from '@/api'

interface FloatingWordPreviewProps {
  word: string
  anchorRect: DOMRect
  onClose: () => void
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
  // 取前150个字符
  const brief = text.slice(0, 150).trim()
  return brief.length < text.length ? brief + '...' : brief
}

// 从定义中尝试提取音标
function extractPhonetic(definition: string): string | null {
  const phoneticMatch = definition.match(/[\/\[][^\]\/]{1,30}[\/\]]/)
  return phoneticMatch ? phoneticMatch[0] : null
}

export function FloatingWordPreview({ word, anchorRect, onClose, onWordClick }: FloatingWordPreviewProps) {
  const containerRef = useRef<HTMLDivElement>(null)

  const { data, isLoading, error } = useQuery({
    queryKey: ['preview', word],
    queryFn: () => searchApi.search(word),
    enabled: word.trim().length > 0,
    staleTime: 1000 * 60 * 10, // 10 分钟缓存
  })

  // 点击外部关闭
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        onClose()
      }
    }

    // 延迟添加事件监听，避免立即触发
    const timer = setTimeout(() => {
      document.addEventListener('click', handleClickOutside)
    }, 100)

    return () => {
      clearTimeout(timer)
      document.removeEventListener('click', handleClickOutside)
    }
  }, [onClose])

  const handleClick = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    onWordClick(word)
  }

  // 计算位置 - 优先显示在上方
  const calculatePosition = () => {
    const padding = 8
    const previewWidth = 280
    const previewHeight = 150 // 估算高度

    let left = anchorRect.left + anchorRect.width / 2 - previewWidth / 2
    let top = anchorRect.top - previewHeight - padding

    // 确保不超出左边界
    if (left < padding) {
      left = padding
    }
    // 确保不超出右边界
    if (left + previewWidth > window.innerWidth - padding) {
      left = window.innerWidth - previewWidth - padding
    }
    // 如果上方空间不足，显示在下方
    if (top < padding) {
      top = anchorRect.bottom + padding
    }

    return { left, top }
  }

  const position = calculatePosition()
  const firstResult = data?.results?.[0]
  const phonetic = firstResult ? extractPhonetic(firstResult.definition) : null
  const briefDef = firstResult ? extractBriefDefinition(firstResult.definition) : null

  return createPortal(
    <div
      ref={containerRef}
      className="word-preview-popover fixed z-[9999] animate-fade-in"
      style={{
        left: position.left,
        top: position.top,
      }}
      onMouseLeave={onClose}
    >
      <div className="bg-content1 border border-default-200 dark:border-default-100 rounded-xl shadow-lg p-3 min-w-[200px] max-w-[280px]">
        {/* 小箭头 */}
        <div 
          className="absolute w-3 h-3 bg-content1 border-b border-r border-default-200 dark:border-default-100 rotate-45"
          style={{
            left: '50%',
            bottom: -6,
            marginLeft: -6,
          }}
        />
        
        {isLoading ? (
          <div className="flex items-center justify-center py-4">
            <Spinner size="sm" />
            <span className="ml-2 text-sm text-default-500">加载中...</span>
          </div>
        ) : error ? (
          <div className="text-sm text-danger-500 py-2">加载失败</div>
        ) : !firstResult ? (
          <div className="text-sm text-default-500 py-2">未找到 "{word}" 的释义</div>
        ) : (
          <>
            <div className="flex items-center gap-2 mb-2">
              <h4 className="font-bold text-base text-default-800 dark:text-default-200">
                {word}
              </h4>
              {phonetic && (
                <span className="text-sm text-default-400">{phonetic}</span>
              )}
            </div>
            <p className="text-sm text-default-600 dark:text-default-400 leading-relaxed">
              {briefDef}
            </p>
            <div className="mt-3 pt-2 border-t border-default-100 dark:border-default-50">
              <button
                onClick={handleClick}
                className="text-xs text-primary-500 hover:text-primary-600 dark:hover:text-primary-400 transition-colors font-medium"
              >
                查看完整释义 →
              </button>
            </div>
          </>
        )}
      </div>
    </div>,
    document.body
  )
}
