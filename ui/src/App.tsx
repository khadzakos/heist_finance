import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import MarketDashboard from "./components/MarketDashboard";
import ExchangePage from "./pages/ExchangePage";
import AssetDetailsPage from "./pages/AssetDetailsPage";
import NotFoundPage from "./pages/NotFoundPage";

const App: React.FC = () => {
  return (
    <Router>
      <div className="min-h-screen bg-gray-100 text-gray-800">
        <header className="bg-white shadow">
          <div className="max-w-7xl mx-auto py-4 px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between items-center">
              <h1 className="text-3xl font-bold text-gray-900">Heist Finance</h1>
              <nav className="flex space-x-4">
                <a href="/" className="text-gray-700 hover:text-gray-900">Dashboard</a>
              </nav>
            </div>
          </div>
        </header>
        <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
          <Routes>
            <Route path="/" element={<MarketDashboard />} />
            <Route path="/exchange/:exchange" element={<ExchangePage />} />
            <Route path="/exchange/:exchange/asset/:symbol" element={<AssetDetailsPage />} />
            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
};

export default App;