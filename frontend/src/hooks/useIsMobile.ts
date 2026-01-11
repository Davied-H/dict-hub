import { useState, useEffect } from 'react'

const MOBILE_BREAKPOINT = 640 // sm breakpoint in Tailwind

export function useIsMobile() {
  const [isMobile, setIsMobile] = useState(() => {
    if (typeof window === 'undefined') return false
    return window.innerWidth < MOBILE_BREAKPOINT
  })

  useEffect(() => {
    const mediaQuery = window.matchMedia(`(max-width: ${MOBILE_BREAKPOINT - 1}px)`)
    
    const handleChange = (e: MediaQueryListEvent) => {
      setIsMobile(e.matches)
    }

    // 设置初始值
    setIsMobile(mediaQuery.matches)

    // 监听变化
    mediaQuery.addEventListener('change', handleChange)
    return () => mediaQuery.removeEventListener('change', handleChange)
  }, [])

  return isMobile
}
