import { Tabs, Tab, Chip } from '@heroui/react'
import { useSearchStore } from '@/stores'
import type { SearchResult } from '@/types'

// 与 DictCard 保持一致的颜色配置
const DICT_COLORS = [
  'bg-blue-500',
  'bg-emerald-500',
  'bg-violet-500',
  'bg-amber-500',
  'bg-rose-500',
]

function getColorByDictId(dictId: number) {
  return DICT_COLORS[dictId % DICT_COLORS.length]
}

interface DictFilterTabsProps {
  results: SearchResult[]
}

interface DictStat {
  dictId: number
  dictTitle: string
  count: number
}

export function DictFilterTabs({ results }: DictFilterTabsProps) {
  const { activeDictFilter, setActiveDictFilter } = useSearchStore()

  // 统计每个词典的结果数
  const dictStats: DictStat[] = []
  const dictMap = new Map<number, DictStat>()

  results.forEach((r) => {
    if (!dictMap.has(r.dict_id)) {
      const stat = {
        dictId: r.dict_id,
        dictTitle: r.dict_title,
        count: 0,
      }
      dictMap.set(r.dict_id, stat)
      dictStats.push(stat)
    }
    dictMap.get(r.dict_id)!.count++
  })

  // 如果只有一个词典，不显示筛选
  if (dictStats.length <= 1) {
    return null
  }

  const totalCount = results.length

  return (
    <div className="mb-4 overflow-x-auto pb-1 -mx-2 px-2">
      <Tabs
        selectedKey={activeDictFilter === null ? 'all' : String(activeDictFilter)}
        onSelectionChange={(key) => {
          if (key === 'all') {
            setActiveDictFilter(null)
          } else {
            setActiveDictFilter(Number(key))
          }
        }}
        variant="underlined"
        color="primary"
        size="sm"
        classNames={{
          tabList: 'gap-4 w-full relative rounded-none p-0 border-b border-divider',
          cursor: 'w-full bg-primary',
          tab: 'max-w-fit px-0 h-10',
          tabContent: 'group-data-[selected=true]:text-primary font-medium',
        }}
      >
        <Tab
          key="all"
          title={
            <div className="flex items-center gap-2">
              <span>全部</span>
              <Chip size="sm" variant="flat" className="text-xs min-w-[24px] h-5">
                {totalCount}
              </Chip>
            </div>
          }
        />
        {dictStats.map((stat) => (
          <Tab
            key={String(stat.dictId)}
            title={
              <div className="flex items-center gap-2">
                <span className={`w-2 h-2 rounded-full ${getColorByDictId(stat.dictId)}`} />
                <span className="truncate max-w-[120px]">{stat.dictTitle}</span>
                <Chip size="sm" variant="flat" className="text-xs min-w-[24px] h-5">
                  {stat.count}
                </Chip>
              </div>
            }
          />
        ))}
      </Tabs>
    </div>
  )
}
