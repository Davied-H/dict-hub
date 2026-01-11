import { Card, CardHeader, CardBody, Divider } from '@heroui/react'
import { AudioPlayer } from './AudioPlayer'
import { DictContent } from './DictContent'
import type { SearchResult } from '@/types'

function BookIcon({ className }: { className?: string }) {
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
        strokeWidth={1.5}
        d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
      />
    </svg>
  )
}

interface DictCardProps {
  dictTitle: string
  results: SearchResult[]
}

export function DictCard({ dictTitle, results }: DictCardProps) {
  return (
    <Card className="dict-card overflow-hidden shadow-sm border border-default-200 dark:border-default-100">
      {/* 字典标题栏 */}
      <CardHeader className="dict-card-header bg-gradient-to-r from-primary-50 to-secondary-50 dark:from-primary-900/20 dark:to-secondary-900/20 py-3 px-5">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-primary/10 rounded-lg">
            <BookIcon className="w-5 h-5 text-primary" />
          </div>
          <div>
            <h2 className="font-semibold text-default-800 dark:text-default-200">
              {dictTitle}
            </h2>
          </div>
        </div>
      </CardHeader>

      <Divider />

      {/* 词条内容区域 */}
      <CardBody className="p-0">
        {results.map((result, idx) => (
          <article
            key={idx}
            className={`dict-entry p-5 ${
              idx !== results.length - 1
                ? 'border-b border-default-100 dark:border-default-50'
                : ''
            }`}
          >
            {/* 单词标题行 */}
            <header className="flex items-start justify-between gap-4 mb-4">
              <div className="flex items-center gap-3">
                <h3 className="text-2xl font-bold text-primary-600 dark:text-primary-400 tracking-tight">
                  {result.word}
                </h3>
                <AudioPlayer definition={result.definition} />
              </div>
            </header>

            {/* 释义内容 */}
            <div className="dict-definition">
              <DictContent html={result.definition} />
            </div>
          </article>
        ))}
      </CardBody>
    </Card>
  )
}
