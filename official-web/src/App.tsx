import Navbar from './components/Navbar';
import Hero from './components/Hero';
import Features from './components/Features';
import TechStack from './components/TechStack';
import QuickStart from './components/QuickStart';
import CTA from './components/CTA';
import Footer from './components/Footer';

function App() {
  return (
    <div className="min-h-screen bg-white">
      <Navbar />
      <main>
        <Hero />
        <Features />
        <TechStack />
        <QuickStart />
        <CTA />
      </main>
      <Footer />
    </div>
  );
}

export default App;
