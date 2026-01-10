import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface SearchState {
  keyword: string
  recentSearches: string[]
  setKeyword: (keyword: string) => void
  addToRecent: (word: string) => void
  clearRecent: () => void
}

export const useSearchStore = create<SearchState>()(
  persist(
    (set, get) => ({
      keyword: '',
      recentSearches: [],
      setKeyword: (keyword) => set({ keyword }),
      addToRecent: (word) => {
        const recent = get().recentSearches.filter((w) => w !== word)
        set({ recentSearches: [word, ...recent].slice(0, 20) })
      },
      clearRecent: () => set({ recentSearches: [] }),
    }),
    { name: 'search-storage' }
  )
)
