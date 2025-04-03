import React, { useEffect, useState, lazy, Suspense } from 'react';
import { fetchHomePageData } from '../services/api';
import { HomePageData, MarketData } from '../types';
import { Link } from 'react-router-dom';
import PopularQuotes from './PopularQuotes';

// Lazy-loaded components - specify relative paths
const MarketSection = lazy(() => import('../components/MarketSection'));
const MarketChart = lazy(() => import('../components/MarketChart'));

// Loading component
const LoadingComponent = () => (
  <div className="bg-white p-4 rounded-lg shadow-md">
    <div className="h-32 flex items-center justify-center">
      <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
    </div>
  </div>
);

const MarketDashboard: React.FC = () => {
  const [data, setData] = useState<HomePageData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [popularQuotes, setPopularQuotes] = useState<MarketData[]>([]);
  const [topGainers, setTopGainers] = useState<MarketData[]>([]);
  const [topLosers, setTopLosers] = useState<MarketData[]>([]);
  const [cryptoExchanges, setCryptoExchanges] = useState<Record<string, MarketData[]>>({});
  const [stockExchanges, setStockExchanges] = useState<Record<string, MarketData[]>>({});
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date());

  // Check if priceChangePercent is negative
  const isNegativeChange = (asset: MarketData): boolean => {
    if (!asset.priceChangePercent) return false;
    return asset.priceChangePercent.toString().startsWith('-');
  };

  // Get numeric value from priceChangePercent
  const getPercentValue = (asset: MarketData): number => {
    if (!asset.priceChangePercent) return 0;
    return parseFloat(asset.priceChangePercent.toString().replace('%', '').replace('-', ''));
  };

  // Get top popular quotes for charts
  const getPopularQuotes = (data: HomePageData): MarketData[] => {
    // Get all assets sorted by volume
    const allAssets = [...data.crypto, ...data.stock];
    
    // Remove duplicates (in case the same asset appears in both markets)
    const uniqueAssets = new Map<string, MarketData>();
    allAssets.forEach(asset => {
      const key = `${asset.exchange}:${asset.symbol}`;
      if (!uniqueAssets.has(key)) {
        uniqueAssets.set(key, asset);
      }
    });
    
    // Convert to array and sort by volume (descending)
    const sortedAssets = Array.from(uniqueAssets.values())
      .filter(asset => asset.volume !== undefined && asset.volume > 0)
      .sort((a, b) => (b.volume || 0) - (a.volume || 0));
    
    return sortedAssets.slice(0, Math.min(9, sortedAssets.length));
  };

  // Group data by exchange
  const groupByExchange = (marketData: MarketData[]): Record<string, MarketData[]> => {
    const exchangeGroups = marketData.reduce((acc, item) => {
      if (!acc[item.exchange]) {
        acc[item.exchange] = [];
      }
      acc[item.exchange].push(item);
      return acc;
    }, {} as Record<string, MarketData[]>);

    // Sort items within each exchange by volume and limit to 15 per exchange
    Object.keys(exchangeGroups).forEach(exchange => {
      exchangeGroups[exchange] = exchangeGroups[exchange]
        .filter(item => item.volume !== undefined)
        .sort((a, b) => (b.volume || 0) - (a.volume || 0))
        .slice(0, 15);
    });

    return exchangeGroups;
  };

  // Select top performing assets for lists
  const getTopPerformers = (marketData: MarketData[], limit: number = 10): MarketData[] => {
    // Create a map to track unique assets by exchange-symbol
    const uniqueAssets = new Map<string, MarketData>();
    
    // Add all assets to the map
    marketData.forEach(asset => {
      const key = `${asset.exchange}:${asset.symbol}`;
      uniqueAssets.set(key, asset);
    });
    
    // Convert to array, filter for positive changes, sort by percentage change
    return Array.from(uniqueAssets.values())
      .filter(asset => asset.priceChangePercent && !isNegativeChange(asset))
      .sort((a, b) => getPercentValue(b) - getPercentValue(a))
      .slice(0, limit);
  };

  // Select worst performing assets for lists
  const getWorstPerformers = (marketData: MarketData[], limit: number = 10): MarketData[] => {
    // Create a map to track unique assets by exchange-symbol
    const uniqueAssets = new Map<string, MarketData>();
    
    // Add all assets to the map
    marketData.forEach(asset => {
      const key = `${asset.exchange}:${asset.symbol}`;
      uniqueAssets.set(key, asset);
    });
    
    // Convert to array, filter for negative changes, sort by percentage change
    return Array.from(uniqueAssets.values())
      .filter(asset => asset.priceChangePercent && isNegativeChange(asset))
      .sort((a, b) => getPercentValue(b) - getPercentValue(a))
      .slice(0, limit);
  };

  // Initial data fetch and update every 30 seconds
  useEffect(() => {
    const mountedRef = { current: true };
    const dataRef = { current: data }; // Для отслеживания текущих данных
    let isUpdating = false; // Флаг для предотвращения параллельных запросов
    
    const fetchData = async () => {
      // Предотвращаем несколько одновременных запросов
      if (isUpdating || !mountedRef.current) return;
      
      isUpdating = true;
      
      try {
        if (mountedRef.current) setLoading(true);
        const result = await fetchHomePageData();
        
        // Сравниваем с последними данными, чтобы избежать лишних обновлений
        if (!mountedRef.current) {
          isUpdating = false;
          return;
        }
        
        // Make sure we handle empty data
        if (!result.crypto || !result.stock || 
           (result.crypto.length === 0 && result.stock.length === 0)) {
          console.warn('API returned empty data');
          if (mountedRef.current) {
            setError('No market data available at the moment');
            setLoading(false);
          }
          isUpdating = false;
          return;
        }
        
        // Обновляем данные только если компонент смонтирован
        if (mountedRef.current) {
          // Запоминаем новые данные для сравнения в будущем
          dataRef.current = result;
          
          setData(result);
          setLastUpdated(new Date());
          
          // Process data for all components
          const quotes = getPopularQuotes(result);
          setPopularQuotes(quotes);
          
          // Set top gainers/losers lists
          const allAssets = [...result.crypto, ...result.stock];
          if (allAssets.length > 0) {
            setTopGainers(getTopPerformers(allAssets, 10));
            setTopLosers(getWorstPerformers(allAssets, 10));
          }
          
          setCryptoExchanges(groupByExchange(result.crypto));
          setStockExchanges(groupByExchange(result.stock));
          
          setLoading(false);
          setError(null);
        }
      } catch (err) {
        console.error('Failed to fetch market data:', err);
        if (mountedRef.current) {
          setError('Failed to fetch market data');
          setLoading(false);
        }
      } finally {
        isUpdating = false;
      }
    };

    // Немедленный первый запрос
    fetchData();

    // Увеличиваем интервал обновления до 60 секунд для снижения нагрузки
    const intervalId = setInterval(fetchData, 60000);

    // Cleanup function to prevent updates after unmount
    return () => {
      mountedRef.current = false;
      clearInterval(intervalId);
    };
  }, []); // Empty dependency array since we don't need to re-run this effect

  // Render dashboard main content
  const renderMainContent = () => {
    // Get the top 5 exchanges for each market type based on total volume
    const getTopExchanges = (market: 'crypto' | 'stock', exchangeMap: Record<string, MarketData[]>): Record<string, MarketData[]> => {
      const exchanges = Object.keys(exchangeMap);
      if (exchanges.length <= 5) return exchangeMap;
      
      // Calculate total volume per exchange
      const volumeByExchange: Record<string, number> = {};
      exchanges.forEach(exchange => {
        volumeByExchange[exchange] = exchangeMap[exchange].reduce((sum, asset) => sum + (asset.volume || 0), 0);
      });
      
      // Sort exchanges by total volume and take top 5
      const top5Exchanges = exchanges
        .sort((a, b) => volumeByExchange[b] - volumeByExchange[a])
        .slice(0, 5);
      
      // Return filtered map with only top 5 exchanges
      const result: Record<string, MarketData[]> = {};
      top5Exchanges.forEach(exchange => {
        result[exchange] = exchangeMap[exchange];
      });
      
      return result;
    };
    
    // Get top 5 exchanges for each market
    const topCryptoExchanges = getTopExchanges('crypto', cryptoExchanges);
    const topStockExchanges = getTopExchanges('stock', stockExchanges);
    
    return (
      <>
        <div className="col-span-1 md:col-span-3 mb-6">
          <div className="flex justify-between items-center mb-3">
            <h2 className="text-xl font-bold">Popular Quotes</h2>
            <div className="text-sm text-gray-500">
              Updated: {lastUpdated.toLocaleTimeString()}
            </div>
          </div>
          <div className="h-[350px]">
            <Suspense fallback={<div className="h-full w-full flex items-center justify-center">Loading quotes...</div>}>
              <PopularQuotes quotes={popularQuotes} />
            </Suspense>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-8">
          <Suspense fallback={<LoadingComponent />}>
            <MarketSection 
              title="Crypto Market" 
              exchanges={topCryptoExchanges} 
              market="crypto" 
            />
          </Suspense>
          
          <Suspense fallback={<LoadingComponent />}>
            <MarketSection 
              title="Stock Market" 
              exchanges={topStockExchanges} 
              market="stock" 
            />
          </Suspense>
        </div>
      </>
    );
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-4 bg-red-100 text-red-700 rounded-md">
        <p>{error}</p>
      </div>
    );
  }

  if (!data) return null;

  return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Market Dashboard</h1>
        <p className="text-gray-600">
          Real-time data from multiple markets
        </p>
      </div>
      
      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {[...Array(3)].map((_, i) => (
            <LoadingComponent key={i} />
          ))}
        </div>
      ) : error ? (
        <div className="bg-red-100 p-4 rounded-md text-red-700">
          {error}
        </div>
      ) : (
        renderMainContent()
      )}
    </div>
  );
};

export default MarketDashboard; 