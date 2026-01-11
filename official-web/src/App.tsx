import { Routes, Route } from 'react-router-dom';
import HomePage from './pages/HomePage';
import DocsPage from './pages/DocsPage';
import DocsQuickstartPage from './pages/DocsQuickstartPage';
import DocsDockerPage from './pages/DocsDockerPage';
import DocsConfigurationPage from './pages/DocsConfigurationPage';
import DocsDictionaryGuidePage from './pages/DocsDictionaryGuidePage';
import DocsApiPage from './pages/DocsApiPage';
import DocsFaqPage from './pages/DocsFaqPage';

function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/docs" element={<DocsPage />} />
      <Route path="/docs/quickstart" element={<DocsQuickstartPage />} />
      <Route path="/docs/docker" element={<DocsDockerPage />} />
      <Route path="/docs/configuration" element={<DocsConfigurationPage />} />
      <Route path="/docs/dictionary-guide" element={<DocsDictionaryGuidePage />} />
      <Route path="/docs/api" element={<DocsApiPage />} />
      <Route path="/docs/faq" element={<DocsFaqPage />} />
    </Routes>
  );
}

export default App;
