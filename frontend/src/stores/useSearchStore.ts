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
  // 键盘导航：当前选中的建议索引 (-1 表示无选中)
  selectedSuggestionIndex: number
  // 词典筛选：当前选中的词典ID (null 表示全部)
  activeDictFilter: number | null

  // Actions
  setKeyword: (keyword: string) => void
  submitSearch: (word?: string) => void
  setShowSuggestions: (show: boolean) => void
  setSelectedSuggestionIndex: (index: number) => void
  setActiveDictFilter: (dictId: number | null) => void
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
      selectedSuggestionIndex: -1,
      activeDictFilter: null,

      setKeyword: (keyword) => set({ keyword, selectedSuggestionIndex: -1 }),

      submitSearch: (word) => {
        const searchWord = word ?? get().keyword.trim()
        if (searchWord) {
          set({ 
            submittedKeyword: searchWord, 
            showSuggestions: false,
            selectedSuggestionIndex: -1,
            activeDictFilter: null, // 重置词典筛选
          })
          get().addToRecent(searchWord)
        }
      },

      setShowSuggestions: (show) => set({ 
        showSuggestions: show,
        selectedSuggestionIndex: show ? -1 : get().selectedSuggestionIndex,
      }),

      setSelectedSuggestionIndex: (index) => set({ selectedSuggestionIndex: index }),

      setActiveDictFilter: (dictId) => set({ activeDictFilter: dictId }),

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
