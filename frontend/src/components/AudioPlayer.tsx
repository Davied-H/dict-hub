import { useRef, useState, useCallback, useEffect } from 'react'
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

interface AudioAvailability {
  gb: boolean
  us: boolean
}

interface AudioPlayerProps {
  word: string
  definition?: string
}

export function AudioPlayer({ word }: AudioPlayerProps) {
  const [availability, setAvailability] = useState<AudioAvailability | null>(null)
  const [playingType, setPlayingType] = useState<'gb' | 'us' | null>(null)
  const [loading, setLoading] = useState(true)
  const audioRef = useRef<HTMLAudioElement>(null)

  // 检查音频可用性
  useEffect(() => {
    if (!word) {
      setLoading(false)
      return
    }

    setLoading(true)
    setAvailability(null)
    setPlayingType(null)

    fetch(`/api/v1/audio/${encodeURIComponent(word)}/availability`)
      .then(res => res.json())
      .then(data => {
        if (data.code === 0 && data.data) {
          setAvailability(data.data)
        } else {
          setAvailability({ gb: false, us: false })
        }
      })
      .catch(() => setAvailability({ gb: false, us: false }))
      .finally(() => setLoading(false))
  }, [word])

  const playAudio = useCallback((type: 'gb' | 'us') => {
    if (!audioRef.current) return

    // 如果正在播放同一类型，则暂停
    if (playingType === type) {
      audioRef.current.pause()
      setPlayingType(null)
      return
    }

    // 播放新音频
    audioRef.current.src = `/api/v1/audio/${encodeURIComponent(word)}?type=${type}`
    audioRef.current.play().catch(() => {
      // 播放失败
      setPlayingType(null)
    })
    setPlayingType(type)
  }, [word, playingType])

  const handleEnded = useCallback(() => {
    setPlayingType(null)
  }, [])

  const handlePause = useCallback(() => {
    setPlayingType(null)
  }, [])

  // 加载中或无可用音频时不显示
  if (loading) return null
  if (!availability || (!availability.gb && !availability.us)) return null

  return (
    <div className="flex items-center gap-0.5">
      <audio
        ref={audioRef}
        onEnded={handleEnded}
        onPause={handlePause}
      />

      {availability.gb && (
        <Button
          isIconOnly
          variant="light"
          size="sm"
          onPress={() => playAudio('gb')}
          title="英式发音 (British)"
          className="min-w-unit-8 h-unit-7 px-1"
        >
          <span className="text-[10px] font-medium text-default-500 mr-0.5">UK</span>
          {playingType === 'gb' ? (
            <PauseIcon className="w-3.5 h-3.5" />
          ) : (
            <SpeakerIcon className="w-3.5 h-3.5" />
          )}
        </Button>
      )}

      {availability.us && (
        <Button
          isIconOnly
          variant="light"
          size="sm"
          onPress={() => playAudio('us')}
          title="美式发音 (American)"
          className="min-w-unit-8 h-unit-7 px-1"
        >
          <span className="text-[10px] font-medium text-default-500 mr-0.5">US</span>
          {playingType === 'us' ? (
            <PauseIcon className="w-3.5 h-3.5" />
          ) : (
            <SpeakerIcon className="w-3.5 h-3.5" />
          )}
        </Button>
      )}
    </div>
  )
}
