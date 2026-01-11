import { useRef, useState } from 'react'
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

// 从 HTML 内容中提取音频 URL
function extractAudioUrl(html: string): string | null {
  // 匹配 <audio src="..."> 或 src="xxx.mp3" 模式
  const audioRegex = /src=["']([^"']+\.(mp3|wav|ogg|m4a))["']/i
  const match = html.match(audioRegex)
  return match ? match[1] : null
}

interface AudioPlayerProps {
  definition: string
}

export function AudioPlayer({ definition }: AudioPlayerProps) {
  const audioUrl = extractAudioUrl(definition)
  const audioRef = useRef<HTMLAudioElement>(null)
  const [isPlaying, setIsPlaying] = useState(false)

  if (!audioUrl) return null

  const toggle = () => {
    if (audioRef.current) {
      if (isPlaying) {
        audioRef.current.pause()
      } else {
        audioRef.current.play()
      }
      setIsPlaying(!isPlaying)
    }
  }

  return (
    <>
      <audio
        ref={audioRef}
        src={audioUrl}
        onEnded={() => setIsPlaying(false)}
        onPause={() => setIsPlaying(false)}
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
