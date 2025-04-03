import React, { useEffect, useState, useRef, memo, useCallback, useMemo } from 'react';
import { useParams, Link } from 'react-router-dom';
import { 
  LineChart, Line, XAxis, YAxis, CartesianGrid, 
  Tooltip, Legend, ResponsiveContainer 
} from 'recharts';
import { fetchAssetDetails } from '../services/api';
import { MarketData } from '../types';

// Function to generate historical price data points
const generateHistoricalData = (
  currentPrice: number, 
  periods: number = 24, 
  volatility: number = 0.02
): { time: string; price: number }[] => {
  const data = [];
  let price = currentPrice * 0.9; // Start 10% lower than current
  const now = new Date();
  
  for (let i = periods; i > 0; i--) {
    // Calculate time for this data point (hours ago)
    const timePoint = new Date(now);
    timePoint.setHours(now.getHours() - i);
    
    // Add some random fluctuation based on volatility parameter
    const change = (Math.random() - 0.5) * 2 * volatility * price;
    price += change;
    
    // Gradually move towards the current price
    const weight = i / periods;
    price = price * weight + currentPrice * (1 - weight);
    
    data.push({
      time: timePoint.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      price: price,
    });
  }
  
  // Add current price as the last point
  data.push({
    time: 'Now',
    price: currentPrice,
  });
  
  return data;
};

// Memoized number display component for changing prices
const NumberDisplay = memo(({
  value,
  className
}: {
  value: string | null;
  className?: string;
}) => (
  <span className={className || ""}>
    {value || "N/A"}
  </span>
));

// Memoized price display component
const PriceDisplay = memo(({
  label, 
  value, 
  className
}: {
  label: string;
  value: string | null;
  className?: string;
}) => (
  <div className="flex justify-between items-center border-b pb-2">
    <span className="font-medium text-gray-600">{label}</span>
    <NumberDisplay value={value} className={className} />
  </div>
));

// Memoized chart component
const PriceChart = memo(({
  chartData,
  chartColor
}: {
  chartData: {time: string; price: number}[];
  chartColor: string;
}) => (
  <ResponsiveContainer width="100%" height="100%">
    <LineChart
      data={chartData}
      margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
    >
      <CartesianGrid strokeDasharray="3 3" opacity={0.2} />
      <XAxis dataKey="time" />
      <YAxis 
        domain={[(dataMin: number) => dataMin * 0.995, (dataMax: number) => dataMax * 1.005]}
        tickFormatter={(tick) => `$${tick.toFixed(3)}`}
      />
      <Tooltip 
        formatter={(value: number) => [`$${value.toFixed(3)}`, 'Price']}
      />
      <Legend />
      <Line 
        type="monotone" 
        dataKey="price" 
        stroke={chartColor} 
        activeDot={{ r: 8 }} 
        name="Price"
        dot={false}
        isAnimationActive={false}
      />
    </LineChart>
  </ResponsiveContainer>
));

const AssetDetailsPage: React.FC = () => {
  const { exchange, symbol } = useParams<{ exchange: string; symbol: string }>();
  const [asset, setAsset] = useState<MarketData | null>(null);
  const [chartData, setChartData] = useState<{time: string; price: number}[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [simulatedPrice, setSimulatedPrice] = useState<number | null>(null);
  const [simulatedChange, setSimulatedChange] = useState<string>('');
  const [priceDirection, setPriceDirection] = useState<'up' | 'down' | null>(null);
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date());
  const originalPrice = useRef<number | null>(null);
  const lastPrice = useRef<number | null>(null);
  const updateInterval = useRef<NodeJS.Timeout | null>(null);
  const miniUpdateInterval = useRef<NodeJS.Timeout | null>(null);
  const chartDataRef = useRef<{time: string; price: number}[]>([]);
  // Create a ref to track if component is mounted
  const isMountedRef = useRef(true);
  // Track if initial data has been loaded
  const dataLoadedRef = useRef(false);

  // Check if price change is negative
  const isNegativeChange = useCallback((): boolean => {
    return simulatedChange.startsWith('-');
  }, [simulatedChange]);

  // Determine chart line color based on price change
  const chartColor = isNegativeChange() ? "#ef4444" : "#22c55e";

  // Format the price with currency symbol
  const formatPrice = useCallback((price: number | null | undefined): string => {
    if (price === null || price === undefined) return "N/A";
    return `$${price.toFixed(3)}`;
  }, []);

  // Optimize the fetchData function to batch updates
  const fetchData = useCallback(async () => {
    try {
      if (!exchange || !symbol || !isMountedRef.current) return;
      
      const result = await fetchAssetDetails(exchange, symbol);
      
      if (!isMountedRef.current) return;
      
      // Store original price for calculating changes
      const initialPrice = result.price || 100;
      if (!dataLoadedRef.current) {
        originalPrice.current = initialPrice;
        lastPrice.current = initialPrice;
      }
      
      // Update state directly
      setAsset(result);
      setSimulatedPrice(initialPrice);
      setSimulatedChange(result.priceChangePercent || '0.000%');
      
      // Generate initial chart data if first load
      if (!dataLoadedRef.current) {
        const initialChartData = generateHistoricalData(initialPrice);
        setChartData(initialChartData);
        chartDataRef.current = initialChartData;
        dataLoadedRef.current = true;
      }
      
      setLastUpdated(new Date());
      setLoading(false);
      
      // Setup price update intervals
      setupPriceUpdateIntervals(initialPrice);
    } catch (err) {
      if (isMountedRef.current) {
        setError(`Failed to fetch details for asset: ${symbol}`);
        setLoading(false);
      }
    }
  }, [exchange, symbol]);

  // Setup intervals for price updates
  const setupPriceUpdateIntervals = useCallback((basePrice: number) => {
    if (!isMountedRef.current) return;
    
    // Clear any existing intervals
    if (updateInterval.current) {
      clearInterval(updateInterval.current);
    }
    
    if (miniUpdateInterval.current) {
      clearInterval(miniUpdateInterval.current);
    }
    
    // Set up interval for major price updates every 2 seconds
    updateInterval.current = setInterval(() => {
      if (!isMountedRef.current) return;
      
      // Get current price
      const currentSimPrice = simulatedPrice || originalPrice.current || basePrice;
      const prevPrice = lastPrice.current || originalPrice.current || basePrice;
      
      // Random price movement (-0.5% to +0.5%)
      const movement = (Math.random() - 0.5) * 0.01 * currentSimPrice;
      const newPrice = currentSimPrice + movement;
      
      // Calculate percentage change from original price
      const origPrice = originalPrice.current || basePrice;
      const percentChange = ((newPrice - origPrice) / origPrice) * 100;
      
      // Set price direction for animation with stronger effect
      if (newPrice > prevPrice) {
        setPriceDirection('up');
      } else if (newPrice < prevPrice) {
        setPriceDirection('down');
      }
      
      // Update simulated values
      setSimulatedPrice(newPrice);
      setSimulatedChange(percentChange.toFixed(3) + '%');
      lastPrice.current = newPrice;
      setLastUpdated(new Date());
      
      // Update chart by adding a new data point
      setChartData(prev => {
        const newData = [...prev.slice(1)];
        newData.push({
          time: 'Now',
          price: newPrice
        });
        
        // Rename the previous "Now" to timestamp
        const now = new Date();
        const timeStr = now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
        if (newData.length > 1) {
          newData[newData.length - 2].time = timeStr;
        }
        
        chartDataRef.current = newData;
        return newData;
      });
    }, 2000); // Update every 2 seconds
    
    // Set up interval for minor updates and clearing directions
    miniUpdateInterval.current = setInterval(() => {
      if (!isMountedRef.current) return;
      
      // Clear price direction after 1 second
      setPriceDirection(null);
      
      if (simulatedPrice && chartData.length > 0) {
        // Apply very small random price movement (-0.1% to +0.1%)
        const tinyMovement = (Math.random() - 0.5) * 0.002 * simulatedPrice;
        const tinyNewPrice = simulatedPrice + tinyMovement;
        
        // Only update price, not the chart
        setSimulatedPrice(tinyNewPrice);
      }
    }, 1000); // Update every second
  }, [simulatedPrice, chartData.length]);

  useEffect(() => {
    isMountedRef.current = true;
    dataLoadedRef.current = false;
    
    if (!exchange || !symbol) return;

    // Clear any existing intervals when component mounts or parameters change
    if (updateInterval.current) {
      clearInterval(updateInterval.current);
    }
    
    if (miniUpdateInterval.current) {
      clearInterval(miniUpdateInterval.current);
    }

    fetchData();

    return () => {
      isMountedRef.current = false;
      
      // Clear intervals when component unmounts
      if (updateInterval.current) {
        clearInterval(updateInterval.current);
      }
      
      if (miniUpdateInterval.current) {
        clearInterval(miniUpdateInterval.current);
      }
    };
  }, [exchange, symbol, fetchData]);

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

  if (!asset) {
    return <div>Asset not found</div>;
  }

  // Get price class name for current price
  const getCurrentPriceClassName = () => {
    const baseClass = "font-bold text-lg transition-colors duration-500";
    if (priceDirection === 'up') return `${baseClass} text-green-600 bg-green-50`;
    if (priceDirection === 'down') return `${baseClass} text-red-600 bg-red-50`;
    return baseClass;
  };

  return (
    <div className="p-6 max-w-7xl mx-auto">
      <div className="mb-6 flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">{asset.symbol}</h1>
          <p className="text-gray-600">{asset.exchange}</p>
        </div>
        <div className="flex flex-col gap-2">
          <div className="flex gap-4">
            <Link 
              to={`/exchange/${exchange}`} 
              className="px-4 py-2 bg-gray-200 rounded-md hover:bg-gray-300 transition-colors"
            >
              Back to Exchange
            </Link>
            <Link 
              to="/" 
              className="px-4 py-2 bg-gray-200 rounded-md hover:bg-gray-300 transition-colors"
            >
              Back to Dashboard
            </Link>
          </div>
          <p className="text-xs text-gray-500 text-right">Last updated: {lastUpdated.toLocaleTimeString()}</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-5 gap-6">
        {/* Price Chart */}
        <div className="lg:col-span-3 bg-white p-4 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold mb-4">Price Chart (24H)</h2>
          <div className="h-80">
            {simulatedPrice ? (
              <PriceChart chartData={chartData} chartColor={chartColor} />
            ) : (
              <div className="flex items-center justify-center h-full">
                <p className="text-gray-500">No price data available for chart</p>
              </div>
            )}
          </div>
        </div>

        {/* Asset Information */}
        <div className="lg:col-span-2 bg-white p-4 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold mb-4">Market Data</h2>
          
          <div className="space-y-4">
            <PriceDisplay 
              label="Current Price" 
              value={formatPrice(simulatedPrice)}
              className={getCurrentPriceClassName()}
            />
            
            <PriceDisplay 
              label="24h High" 
              value={formatPrice(asset.high)}
              className="text-green-600"
            />
            
            <PriceDisplay 
              label="24h Low" 
              value={formatPrice(asset.low)}
              className="text-red-600"
            />
            
            <PriceDisplay 
              label="24h Volume" 
              value={asset.volume ? `$${asset.volume.toLocaleString()}` : null}
            />
            
            <PriceDisplay 
              label="Price Change" 
              value={simulatedChange || null}
              className={`font-semibold ${isNegativeChange() ? 'text-red-500' : 'text-green-500'}`}
            />
            
            <div className="flex justify-between items-center pt-2 border-t">
              <span className="font-medium text-gray-600">Last Updated</span>
              <span className="text-sm text-gray-500">
                {lastUpdated.toLocaleString()}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default memo(AssetDetailsPage); 