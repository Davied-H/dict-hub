import { useRef, useState, useCallback, useMemo } from 'react'
import { Button } from '@heroui/react'

function SpeakerIcon({ className }: { className?: string }) {
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
        d="M15.536 8.464a5 5 0 010 7.072m2.828-9.9a9 9 0 010 12.728M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z"
      />
    </svg>
  )
}

function PauseIcon({ className }: { className?: string }) {
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
        d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
  )
}

function SpeakerOffIcon({ className }: { className?: string }) {
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
        d="M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z"
      />
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M17 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2"
      />
    </svg>
  )
}

// 从 HTML 内容中提取音频 URL
function extractAudioUrl(html: string): string | null {
  // 匹配 <audio src="..."> 或 src="xxx.mp3" 模式
  const audioRegex = /src=["']([^"']+\.(mp3|wav|ogg|m4a))["']/i
  const match = html.match(audioRegex)
  return match ? match[1] : null
}

// 生成统一音频 API URL
function getAudioApiUrl(word: string): string {
  return `/api/v1/audio/${encodeURIComponent(word)}`
}

interface AudioPlayerProps {
  definition: string
  word?: string // 可选，用于回退到统一音频 API
}

export function AudioPlayer({ definition, word }: AudioPlayerProps) {
  const extractedUrl = useMemo(() => extractAudioUrl(definition), [definition])
  const fallbackUrl = useMemo(() => word ? getAudioApiUrl(word) : null, [word])
  
  const audioRef = useRef<HTMLAudioElement>(null)
  const [isPlaying, setIsPlaying] = useState(false)
  const [currentUrl, setCurrentUrl] = useState<string | null>(extractedUrl)
  const [hasError, setHasError] = useState(false)
  const [audioUnavailable, setAudioUnavailable] = useState(false)

  // 确定最终使用的音频 URL
  const audioUrl = useMemo(() => {
    // 如果提取的 URL 失败，使用回退 URL
    if (hasError && fallbackUrl) {
      return fallbackUrl
    }
    // 优先使用提取的 URL
    if (extractedUrl) {
      return extractedUrl
    }
    // 没有提取到 URL 时，使用回退 URL
    return fallbackUrl
  }, [extractedUrl, fallbackUrl, hasError])

  // 处理音频加载错误，尝试回退
  const handleError = useCallback(() => {
    if (currentUrl === extractedUrl && fallbackUrl && !hasError) {
      // 提取的 URL 失败，尝试回退
      setHasError(true)
      setCurrentUrl(fallbackUrl)
    } else {
      // 所有音频源都失败，标记为不可用
      setAudioUnavailable(true)
      setIsPlaying(false)
    }
  }, [currentUrl, extractedUrl, fallbackUrl, hasError])

  const toggle = useCallback(() => {
    if (audioRef.current) {
      if (isPlaying) {
        audioRef.current.pause()
      } else {
        audioRef.current.play().catch(() => {
          // 播放失败时尝试回退
          handleError()
        })
      }
      setIsPlaying(!isPlaying)
    }
  }, [isPlaying, handleError])

  // 没有任何可用的音频 URL 时不显示按钮
  if (!audioUrl) return null

  // 音频不可用时显示禁用按钮
  if (audioUnavailable) {
    return (
      <Button
        isIconOnly
        variant="light"
        size="sm"
        isDisabled
        title="此单词暂无音频"
      >
        <SpeakerOffIcon className="w-4 h-4 text-default-300" />
      </Button>
    )
  }

  return (
    <>
      <audio
        ref={audioRef}
        src={audioUrl}
        onEnded={() => setIsPlaying(false)}
        onPause={() => setIsPlaying(false)}
        onError={handleError}
      />
      <Button isIconOnly variant="light" size="sm" onPress={toggle}>
        {isPlaying ? (
          <PauseIcon className="w-4 h-4" />
        ) : (
          <SpeakerIcon className="w-4 h-4" />
        )}
      </Button>
    </>
  )
}
