import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface SearchState {
  // 搜索输入
  keyword: string
  // 已提交的搜索词（触发实际搜索）
  submittedKeyword: string
  // 最近搜索
  recentSearches: string[]
  // 是否显示建议下拉
  showSuggestions: boolean

  // Actions
  setKeyword: (keyword: string) => void
  submitSearch: (word?: string) => void
  setShowSuggestions: (show: boolean) => void
  addToRecent: (word: string) => void
  clearRecent: () => void
}

export const useSearchStore = create<SearchState>()(
  persist(
    (set, get) => ({
      keyword: '',
      submittedKeyword: '',
      recentSearches: [],
      showSuggestions: false,

      setKeyword: (keyword) => set({ keyword }),

      submitSearch: (word) => {
        const searchWord = word ?? get().keyword.trim()
        if (searchWord) {
          set({ submittedKeyword: searchWord, showSuggestions: false })
          get().addToRecent(searchWord)
        }
      },

      setShowSuggestions: (show) => set({ showSuggestions: show }),

      addToRecent: (word) => {
        const recent = get().recentSearches.filter((w) => w !== word)
        set({ recentSearches: [word, ...recent].slice(0, 20) })
      },

      clearRecent: () => set({ recentSearches: [] }),
    }),
    {
      name: 'search-storage',
      partialize: (state) => ({ recentSearches: state.recentSearches }),
    }
  )
)
