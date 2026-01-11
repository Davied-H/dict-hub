import { Routes, Route } from 'react-router-dom'
import Layout from '@/components/Layout'
import { Home, Settings, History, Word } from '@/pages'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<Home />} />
        <Route path="settings" element={<Settings />} />
        <Route path="history" element={<History />} />
        <Route path="word/:word" element={<Word />} />
      </Route>
    </Routes>
  )
}

export default App
