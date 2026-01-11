import { Card, CardHeader, CardBody } from '@heroui/react'
import { AudioPlayer } from './AudioPlayer'
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
        strokeWidth={2}
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
    <Card className="overflow-hidden">
      <CardHeader className="bg-default-100 dark:bg-default-50 py-3">
        <div className="flex items-center gap-2">
          <BookIcon className="w-5 h-5 text-default-500" />
          <span className="font-medium">{dictTitle}</span>
        </div>
      </CardHeader>
      <CardBody className="p-0">
        {results.map((result, idx) => (
          <div
            key={idx}
            className="p-4 border-b last:border-b-0 border-default-100 dark:border-default-50"
          >
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-lg font-semibold">{result.word}</h3>
              <AudioPlayer definition={result.definition} />
            </div>
            <div
              className="dict-content prose dark:prose-invert max-w-none text-sm"
              dangerouslySetInnerHTML={{ __html: result.definition }}
            />
          </div>
        ))}
      </CardBody>
    </Card>
  )
}
