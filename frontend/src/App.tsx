import { Button, Input, Card, CardBody, CardHeader } from '@heroui/react'
import { useSearchStore } from '@/stores'

function App() {
  const { keyword, setKeyword, recentSearches, addToRecent, clearRecent } =
    useSearchStore()

  const handleSearch = () => {
    if (keyword.trim()) {
      addToRecent(keyword.trim())
    }
  }

  return (
    <main className="min-h-screen bg-background p-8">
      <div className="max-w-4xl mx-auto">
        <Card className="mb-8">
          <CardHeader className="flex-col items-start gap-2">
            <h1 className="text-2xl font-bold">Dict Hub</h1>
            <p className="text-default-500">Your Dictionary Application</p>
          </CardHeader>
          <CardBody className="gap-4">
            <div className="flex gap-2">
              <Input
                placeholder="Search for a word..."
                value={keyword}
                onChange={(e) => setKeyword(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                className="flex-1"
              />
              <Button color="primary" onPress={handleSearch}>
                Search
              </Button>
            </div>
          </CardBody>
        </Card>

        {recentSearches.length > 0 && (
          <Card>
            <CardHeader className="flex justify-between items-center">
              <h2 className="text-lg font-semibold">Recent Searches</h2>
              <Button size="sm" variant="light" onPress={clearRecent}>
                Clear
              </Button>
            </CardHeader>
            <CardBody>
              <div className="flex flex-wrap gap-2">
                {recentSearches.map((word) => (
                  <Button
                    key={word}
                    size="sm"
                    variant="flat"
                    onPress={() => setKeyword(word)}
                  >
                    {word}
                  </Button>
                ))}
              </div>
            </CardBody>
          </Card>
        )}
      </div>
    </main>
  )
}

export default App
