import React, { useEffect, useState, useRef } from 'react';
import { useParams, Link } from 'react-router-dom';
import { fetchExchangeData } from '../services/api';
import { MarketData } from '../types';

const ExchangeDetailsPage: React.FC = () => {
  const { exchange } = useParams<{ exchange: string }>();
  const [assets, setAssets] = useState<MarketData[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [simulatedPrices, setSimulatedPrices] = useState<Map<string, number>>(new Map());
  const [simulatedChanges, setSimulatedChanges] = useState<Map<string, string>>(new Map());
  const [priceDirections, setPriceDirections] = useState<Map<string, 'up' | 'down' | null>>(new Map());
  const originalPrices = useRef<Map<string, number>>(new Map());
  const previousPrices = useRef<Map<string, number>>(new Map());
  const updateInterval = useRef<NodeJS.Timeout | null>(null);
  const miniUpdateInterval = useRef<NodeJS.Timeout | null>(null);
  
  // Top performers lists
  const [topGainers, setTopGainers] = useState<MarketData[]>([]);
  const [topLosers, setTopLosers] = useState<MarketData[]>([]);

  useEffect(() => {
    if (!exchange) return;

    // Clear any existing intervals when component mounts or parameters change
    if (updateInterval.current) {
      clearInterval(updateInterval.current);
    }
    
    if (miniUpdateInterval.current) {
      clearInterval(miniUpdateInterval.current);
    }

    const fetchData = async () => {
      try {
        const result = await fetchExchangeData(exchange);
        setAssets(result);
        
        // Update top performers lists
        updateTopPerformers(result);
        
        // Store initial prices and set up simulated values
        const priceMap = new Map<string, number>();
        const changeMap = new Map<string, string>();
        const origPriceMap = new Map<string, number>();
        const directionMap = new Map<string, 'up' | 'down' | null>();
        
        result.forEach((asset: MarketData) => {
          const price = asset.price || 100;
          const symbol = asset.symbol;
          priceMap.set(symbol, price);
          changeMap.set(symbol, asset.priceChangePercent || '0.000%');
          origPriceMap.set(symbol, price);
          directionMap.set(symbol, null);
        });
        
        setSimulatedPrices(priceMap);
        setSimulatedChanges(changeMap);
        setPriceDirections(directionMap);
        originalPrices.current = origPriceMap;
        previousPrices.current = new Map(priceMap);
        
        setLoading(false);
        
        // Set up interval for significant price updates every 2 seconds
        updateInterval.current = setInterval(() => {
          setSimulatedPrices(prev => {
            const newPrices = new Map(prev);
            const newChanges = new Map(simulatedChanges);
            const newDirections = new Map(priceDirections);
            
            // Update each asset's price
            result.forEach((asset: MarketData) => {
              const symbol = asset.symbol;
              const currentPrice = newPrices.get(symbol) || 100;
              const prevPrice = previousPrices.current.get(symbol) || currentPrice;
              const origPrice = originalPrices.current.get(symbol) || currentPrice;
              
              // Random price movement (-0.5% to +0.5%)
              const movement = (Math.random() - 0.5) * 0.01 * currentPrice;
              const newPrice = currentPrice + movement;
              
              // Set price direction for animation
              if (newPrice > prevPrice) {
                newDirections.set(symbol, 'up');
              } else if (newPrice < prevPrice) {
                newDirections.set(symbol, 'down');
              }
              
              // Calculate percentage change from original price
              const percentChange = ((newPrice - origPrice) / origPrice) * 100;
              
              // Update maps
              newPrices.set(symbol, newPrice);
              newChanges.set(symbol, percentChange.toFixed(3) + '%');
            });
            
            // Update top performers lists with new data
            updateTopPerformersWithSimulatedData(result, newChanges);
            
            // Store previous prices for next comparison
            previousPrices.current = new Map(newPrices);
            
            // Update directions and changes state
            setPriceDirections(newDirections);
            setSimulatedChanges(newChanges);
            
            return newPrices;
          });
        }, 2000); // Update every 2 seconds
        
        // Set up interval for minor price updates and clearing directions
        miniUpdateInterval.current = setInterval(() => {
          // Clear price direction highlights after 1 second
          setPriceDirections(prev => {
            const newDirections = new Map(prev);
            
            prev.forEach((direction, symbol) => {
              if (direction !== null) {
                newDirections.set(symbol, null);
              }
            });
            
            return newDirections;
          });
        }, 1000); // Update every second
      } catch (err) {
        setError(`Failed to fetch details for exchange: ${exchange}`);
        setLoading(false);
      }
    };

    fetchData();

    return () => {
      // Clear intervals when component unmounts
      if (updateInterval.current) {
        clearInterval(updateInterval.current);
      }
      
      if (miniUpdateInterval.current) {
        clearInterval(miniUpdateInterval.current);
      }
    };
  }, [exchange]);
  
  // Update top performers lists
  const updateTopPerformers = (assets: MarketData[]) => {
    // Get top gainers (positive change, sorted descending)
    const gainers = assets
      .filter(asset => asset.priceChangePercent && !isNegativeChange(asset.priceChangePercent))
      .sort((a, b) => getPercentValue(b.priceChangePercent) - getPercentValue(a.priceChangePercent))
      .slice(0, 5);
    
    // Get top losers (negative change, sorted descending by absolute value)
    const losers = assets
      .filter(asset => asset.priceChangePercent && isNegativeChange(asset.priceChangePercent))
      .sort((a, b) => getPercentValue(a.priceChangePercent) - getPercentValue(b.priceChangePercent))
      .slice(0, 5);
    
    setTopGainers(gainers);
    setTopLosers(losers);
  };
  
  // Update top performers using simulated data
  const updateTopPerformersWithSimulatedData = (assets: MarketData[], changes: Map<string, string>) => {
    // Create copies with updated changes
    const assetsWithUpdatedChanges = assets.map(asset => {
      const newAsset = {...asset};
      const symbol = asset.symbol;
      const change = changes.get(symbol);
      if (change) {
        newAsset.priceChangePercent = change;
      }
      return newAsset;
    });
    
    // Update top performers lists
    updateTopPerformers(assetsWithUpdatedChanges);
  };
  
  // Helper to extract numeric value from percentage string
  const getPercentValue = (percentStr?: string): number => {
    if (!percentStr) return 0;
    return Math.abs(parseFloat(percentStr.replace('%', '')));
  };

  // Helper function to check if change is negative
  const isNegativeChange = (change: string): boolean => {
    return change.startsWith('-');
  };

  // Helper to get current price and format it
  const getCurrentPrice = (symbol: string): string => {
    const price = simulatedPrices.get(symbol);
    return price ? price.toFixed(3) : 'N/A';
  };

  // Helper to get current change
  const getCurrentChange = (symbol: string): string => {
    return simulatedChanges.get(symbol) || 'N/A';
  };
  
  // Helper to get price direction for highlighting
  const getPriceDirection = (symbol: string): 'up' | 'down' | null => {
    return priceDirections.get(symbol) || null;
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

  return (
    <div className="p-6 max-w-7xl mx-auto">
      <div className="mb-6 flex justify-between items-center">
        <h1 className="text-3xl font-bold">{exchange} Exchange</h1>
        <Link 
          to="/" 
          className="px-4 py-2 bg-gray-200 rounded-md hover:bg-gray-300 transition-colors"
        >
          Back to Dashboard
        </Link>
      </div>
      
      {/* Top gainers and losers sections */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        {/* Top Gainers */}
        <div className="bg-white p-4 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold mb-4">Top Gainers</h2>
          {topGainers.length > 0 ? (
            <div className="divide-y divide-gray-100">
              {topGainers.map(asset => (
                <Link 
                  key={asset.symbol}
                  to={`/exchange/${exchange}/asset/${asset.symbol}`} 
                  className="py-3 flex justify-between hover:bg-gray-50 transition-colors"
                >
                  <div className="font-medium">{asset.symbol}</div>
                  <div className="flex flex-col items-end">
                    <div className={`font-medium transition-colors duration-300 ${
                      getPriceDirection(asset.symbol) === 'up' ? 'text-green-500' : 
                      getPriceDirection(asset.symbol) === 'down' ? 'text-red-500' : ''
                    }`}>
                      ${getCurrentPrice(asset.symbol)}
                    </div>
                    <div className="text-green-500">{getCurrentChange(asset.symbol)}</div>
                  </div>
                </Link>
              ))}
            </div>
          ) : (
            <div className="text-center py-4 text-gray-500">No gainers available</div>
          )}
        </div>
        
        {/* Top Losers */}
        <div className="bg-white p-4 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold mb-4">Top Losers</h2>
          {topLosers.length > 0 ? (
            <div className="divide-y divide-gray-100">
              {topLosers.map(asset => (
                <Link 
                  key={asset.symbol}
                  to={`/exchange/${exchange}/asset/${asset.symbol}`} 
                  className="py-3 flex justify-between hover:bg-gray-50 transition-colors"
                >
                  <div className="font-medium">{asset.symbol}</div>
                  <div className="flex flex-col items-end">
                    <div className={`font-medium transition-colors duration-300 ${
                      getPriceDirection(asset.symbol) === 'up' ? 'text-green-500' : 
                      getPriceDirection(asset.symbol) === 'down' ? 'text-red-500' : ''
                    }`}>
                      ${getCurrentPrice(asset.symbol)}
                    </div>
                    <div className="text-red-500">{getCurrentChange(asset.symbol)}</div>
                  </div>
                </Link>
              ))}
            </div>
          ) : (
            <div className="text-center py-4 text-gray-500">No losers available</div>
          )}
        </div>
      </div>
      
      {/* Assets table */}
      <div className="bg-white p-4 rounded-lg shadow-md mb-6">
        <h2 className="text-xl font-semibold mb-4">All Assets</h2>
        
        {assets.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Symbol
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Price
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    24h Change
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    24h High
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    24h Low
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Volume
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {assets.map((asset) => (
                  <tr key={asset.symbol} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="font-medium">{asset.symbol}</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className={`font-medium transition-colors duration-300 ${
                        getPriceDirection(asset.symbol) === 'up' ? 'text-green-500' : 
                        getPriceDirection(asset.symbol) === 'down' ? 'text-red-500' : ''
                      }`}>
                        ${getCurrentPrice(asset.symbol)}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className={`font-medium ${
                        isNegativeChange(getCurrentChange(asset.symbol)) 
                          ? 'text-red-500' 
                          : 'text-green-500'
                      }`}>
                        {getCurrentChange(asset.symbol)}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div>${asset.high?.toFixed(3) || 'N/A'}</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div>${asset.low?.toFixed(3) || 'N/A'}</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div>{asset.volume?.toLocaleString() || 'N/A'}</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <Link 
                        to={`/exchange/${exchange}/asset/${asset.symbol}`} 
                        className="text-blue-600 hover:text-blue-900"
                      >
                        View Details
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="text-center p-4 text-gray-500">
            No assets available for this exchange
          </div>
        )}
        
        <div className="mt-4 text-sm text-gray-500 text-right">
          Last updated: {new Date().toLocaleString()}
        </div>
      </div>
    </div>
  );
};

export default ExchangeDetailsPage; 