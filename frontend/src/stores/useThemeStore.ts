import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type Theme = 'light' | 'dark' | 'system'

interface ThemeState {
  theme: Theme
  setTheme: (theme: Theme) => void
}

// 应用主题到 document
function applyTheme(theme: Theme) {
  const root = document.documentElement
  const isDark =
    theme === 'dark' ||
    (theme === 'system' &&
      window.matchMedia('(prefers-color-scheme: dark)').matches)

  root.classList.toggle('dark', isDark)
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set) => ({
      theme: 'system',
      setTheme: (theme) => {
        set({ theme })
        applyTheme(theme)
      },
    }),
    { name: 'theme-storage' }
  )
)

// 初始化主题（在 main.tsx 调用）
export function initTheme() {
  const stored = localStorage.getItem('theme-storage')
  const theme: Theme = stored
    ? JSON.parse(stored).state?.theme ?? 'system'
    : 'system'
  applyTheme(theme)

  // 监听系统主题变化
  window
    .matchMedia('(prefers-color-scheme: dark)')
    .addEventListener('change', () => {
      const currentTheme = useThemeStore.getState().theme
      if (currentTheme === 'system') {
        applyTheme('system')
      }
    })
}
