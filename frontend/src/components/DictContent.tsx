import { useMemo, useEffect, useRef, useCallback, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useSearchStore } from '@/stores'
import { useIsMobile } from '@/hooks'
import { FloatingWordPreview } from './FloatingWordPreview'

interface DictContentProps {
  html: string
  dictId?: number
}

// 判断是否为单词链接
function isWordLink(href: string): boolean {
  // 已知的单词协议
  if (href.startsWith('entry://')) return true
  if (href.includes('@@@LINK=')) return true
  
  // 检查 URL 中是否包含 entry:// (被浏览器解析为相对路径的情况)
  // 例如: localhost:3000/api/v1/resources/2/entry://specialty
  if (href.includes('/entry://')) return true
  
  // 排除外部链接
  if (href.startsWith('http://') || href.startsWith('https://')) return false
  
  // 排除资源链接
  if (href.startsWith('/dict-assets/') || href.startsWith('/api/')) return false
  if (href.endsWith('.css') || href.endsWith('.js') || href.endsWith('.png') || 
      href.endsWith('.jpg') || href.endsWith('.gif') || href.endsWith('.svg') ||
      href.endsWith('.mp3') || href.endsWith('.wav') || href.endsWith('.ogg')) return false
  
  // 排除锚点链接
  if (href.startsWith('#')) return false
  
  // 排除 sound:// 和其他特殊协议
  if (href.includes('://') && !href.startsWith('entry://')) return false
  
  // 简单单词链接：字母开头，可包含字母、数字、空格、连字符、撇号
  const cleanHref = href.replace(/^\/+/, '') // 移除开头的斜杠
  return /^[a-zA-Z][a-zA-Z0-9\s'\-]*$/.test(cleanHref)
}

// 从链接提取单词
function extractWordFromHref(href: string): string | null {
  if (href.startsWith('entry://')) {
    return decodeURIComponent(href.replace('entry://', ''))
  }
  // 处理被浏览器解析为相对路径的情况: /api/.../entry://word
  // 例如: localhost:3000/api/v1/resources/2/entry://specialty
  if (href.includes('/entry://')) {
    const match = href.match(/\/entry:\/\/(.+)$/)
    if (match) {
      return decodeURIComponent(match[1])
    }
  }
  if (href.includes('@@@LINK=')) {
    return decodeURIComponent(href.split('@@@LINK=')[1])
  }
  // 普通链接
  const cleanHref = href.replace(/^\/+/, '')
  return cleanHref || null
}

/**
 * 解析并美化字典 HTML 内容
 * 支持常见 MDX 字典的 CSS 类名
 * 保留字典自带的 CSS 样式（通过 /dict-assets/ 路径加载）
 * 支持多种格式的单词跳转链接
 */
export function DictContent({ html, dictId }: DictContentProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const { setKeyword, submitSearch } = useSearchStore()
  const navigate = useNavigate()
  const isMobile = useIsMobile()
  
  // 悬停预览状态
  const [previewState, setPreviewState] = useState<{
    word: string
    rect: DOMRect
  } | null>(null)
  
  // 防抖计时器
  const hoverTimerRef = useRef<number | null>(null)

  // 处理单词点击 - 桌面端触发搜索，移动端导航
  const handleWordClick = useCallback((word: string) => {
    if (isMobile) {
      // 移动端导航到单词页面
      navigate(`/word/${encodeURIComponent(word)}`)
    } else {
      // 桌面端触发搜索
      setKeyword(word)
      submitSearch(word)
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }
    setPreviewState(null)
  }, [isMobile, navigate, setKeyword, submitSearch])

  // 关闭预览
  const closePreview = useCallback(() => {
    if (hoverTimerRef.current) {
      clearTimeout(hoverTimerRef.current)
      hoverTimerRef.current = null
    }
    setPreviewState(null)
  }, [])

  // 提取字典 CSS 链接并动态加载
  useEffect(() => {
    if (!html) return

    // 提取所有 CSS 链接
    const linkRegex = /<link[^>]*href=["']([^"']*\.css)["'][^>]*>/gi
    const cssLinks: string[] = []
    let match

    while ((match = linkRegex.exec(html)) !== null) {
      const href = match[1]
      // 只处理 /dict-assets/ 开头的 CSS（由后端重写）
      if (href.startsWith('/dict-assets/') || href.startsWith('/api/')) {
        cssLinks.push(href)
      }
    }

    // 动态加载 CSS
    cssLinks.forEach((href) => {
      // 检查是否已加载
      const existingLink = document.querySelector(`link[href="${href}"]`)
      if (existingLink) return

      const link = document.createElement('link')
      link.rel = 'stylesheet'
      link.href = href
      link.setAttribute('data-dict-css', dictId?.toString() || 'unknown')
      document.head.appendChild(link)
    })
  }, [html, dictId])

  // 处理链接事件
  useEffect(() => {
    const container = containerRef.current
    if (!container) return

    // 为所有单词链接添加样式标记
    const anchors = container.querySelectorAll('a')
    anchors.forEach((anchor) => {
      const href = anchor.getAttribute('href')
      if (href && isWordLink(href)) {
        anchor.classList.add('dict-word-link')
        anchor.setAttribute('data-word-link', 'true')
      }
    })

    // 点击事件处理
    const handleClick = (e: MouseEvent) => {
      const target = e.target as HTMLElement
      const anchor = target.closest('a')
      if (!anchor) return

      const href = anchor.getAttribute('href')
      if (!href) return

      // 检查是否为单词链接
      if (isWordLink(href)) {
        e.preventDefault()
        e.stopPropagation()
        const word = extractWordFromHref(href)
        if (word) {
          handleWordClick(word)
        }
      }
    }

    // 鼠标进入事件 - 桌面端显示预览（带延迟）
    const handleMouseEnter = (e: MouseEvent) => {
      if (isMobile) return
      
      const target = e.target as HTMLElement
      const anchor = target.closest('a') as HTMLAnchorElement | null
      if (!anchor) return

      const href = anchor.getAttribute('href')
      if (!href || !isWordLink(href)) return

      const word = extractWordFromHref(href)
      if (!word) return

      // 清除之前的计时器
      if (hoverTimerRef.current) {
        clearTimeout(hoverTimerRef.current)
      }

      // 延迟300ms显示预览，避免快速滑过时频繁触发
      hoverTimerRef.current = window.setTimeout(() => {
        const rect = anchor.getBoundingClientRect()
        setPreviewState({ word, rect })
      }, 300)
    }

    // 鼠标离开事件
    const handleMouseLeave = (e: MouseEvent) => {
      const target = e.target as HTMLElement
      const anchor = target.closest('a')
      if (!anchor) return

      // 清除计时器
      if (hoverTimerRef.current) {
        clearTimeout(hoverTimerRef.current)
        hoverTimerRef.current = null
      }

      // 检查是否移动到预览弹窗
      const relatedTarget = e.relatedTarget as HTMLElement
      if (relatedTarget?.closest('.word-preview-popover')) return
      
      setPreviewState(null)
    }

    container.addEventListener('click', handleClick)
    container.addEventListener('mouseenter', handleMouseEnter, true)
    container.addEventListener('mouseleave', handleMouseLeave, true)
    
    return () => {
      container.removeEventListener('click', handleClick)
      container.removeEventListener('mouseenter', handleMouseEnter, true)
      container.removeEventListener('mouseleave', handleMouseLeave, true)
      if (hoverTimerRef.current) {
        clearTimeout(hoverTimerRef.current)
      }
    }
  }, [handleWordClick, isMobile])

  // 处理 HTML 内容
  const processedHtml = useMemo(() => {
    let content = html

    // 移除 link 标签（CSS 已通过 useEffect 动态加载）
    content = content.replace(/<link[^>]*>/gi, '')

    // 处理音频标签，添加类名以便隐藏
    content = content.replace(/<audio/gi, '<audio class="dict-audio"')

    return content
  }, [html])

  return (
    <>
      <div
        ref={containerRef}
        className="dict-content"
        dangerouslySetInnerHTML={{ __html: processedHtml }}
      />
      
      {/* 桌面端悬停预览 */}
      {previewState && !isMobile && (
        <FloatingWordPreview
          word={previewState.word}
          anchorRect={previewState.rect}
          onClose={closePreview}
          onWordClick={handleWordClick}
        />
      )}
      
      <style>{`
        .dict-word-link {
          color: hsl(var(--heroui-primary-600)) !important;
          cursor: pointer !important;
          text-decoration: underline !important;
          text-decoration-style: dotted !important;
          text-underline-offset: 2px !important;
          transition: all 0.2s ease !important;
        }
        .dict-word-link:hover {
          color: hsl(var(--heroui-primary-700)) !important;
          text-decoration-style: solid !important;
        }
        .dark .dict-word-link {
          color: hsl(var(--heroui-primary-400)) !important;
        }
        .dark .dict-word-link:hover {
          color: hsl(var(--heroui-primary-300)) !important;
        }
        @keyframes fade-in {
          from { opacity: 0; transform: translateY(4px); }
          to { opacity: 1; transform: translateY(0); }
        }
        .animate-fade-in {
          animation: fade-in 0.15s ease-out;
        }
      `}</style>
    </>
  )
}
