import React, { useState, useEffect, useMemo, useCallback, memo, useRef } from 'react';
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';
import { MarketData } from '../types';
import { Link } from 'react-router-dom';

// Function to generate historical data for the chart
const generateChartData = (
  currentPrice: number,
  periods: number = 50,
  volatility: number = 0.01
) => {
  const data = [];
  let price = currentPrice * 0.97; // Start slightly lower
  
  for (let i = 0; i < periods; i++) {
    // Add random fluctuation
    const change = (Math.random() - 0.5) * 2 * volatility * price;
    price += change;
    
    data.push({
      time: i,
      price: price
    });
  }
  
  // Ensure the last point is the current price
  data.push({
    time: periods,
    price: currentPrice
  });
  
  return data;
};

// Функция форматирования цены
const formatPrice = (price: number, symbol: string): string => {
  // Криптовалюты типа BTC могут иметь большие значения цен
  if (symbol === 'BTC' || symbol === 'ETH') {
    return price > 1000 ? price.toLocaleString(undefined, { maximumFractionDigits: 2 }) : price.toFixed(2);
  }
  
  // Для альткоинов с маленькой ценой показываем больше знаков после запятой
  if (price < 0.1) {
    return price.toFixed(6);
  } else if (price < 1) {
    return price.toFixed(4);
  } else if (price < 10) {
    return price.toFixed(3);
  } else if (price < 1000) {
    return price.toFixed(2);
  }
  
  // Форматируем большие числа
  return price.toLocaleString(undefined, { maximumFractionDigits: 2 });
};

// Component for a single quote chart
const QuoteChart = memo(({ 
  asset, 
  index,
  compact = false
}: { 
  asset: MarketData, 
  index: number,
  compact?: boolean
}) => {
  const [chartData, setChartData] = useState<any[]>([]);
  const price = asset.price || 0;
  
  // Calculate the percentage change from API data
  const percentChange = asset.priceChangePercent || '0.000%';
  const isNegative = percentChange.startsWith('-');
  
  // Determine chart line color based on overall trend
  const chartColor = isNegative ? "#ef4444" : "#22c55e";
  
  // Используем useRef для предотвращения постоянных обновлений
  const prevSymbolRef = useRef(asset.symbol);
  const chartDataRef = useRef(chartData);
  
  // Initialize chart data with useMemo
  useEffect(() => {
    // Только обновляем данные графика, если символ изменился
    // или данные еще не были сгенерированы
    if (prevSymbolRef.current !== asset.symbol || chartDataRef.current.length === 0) {
      const initialData = generateChartData(price, compact ? 30 : 50, 0.01);
      setChartData(initialData);
      chartDataRef.current = initialData;
      prevSymbolRef.current = asset.symbol;
    }
  }, [asset.symbol, compact]); // Зависимости только символ и размер графика

  // Memoize price display for better performance
  const formattedPrice = useMemo(() => formatPrice(price, asset.symbol), [price, asset.symbol]);
  const formattedHigh = useMemo(() => asset.high ? formatPrice(asset.high, asset.symbol) : "N/A", [asset.high, asset.symbol]);
  const formattedLow = useMemo(() => asset.low ? formatPrice(asset.low, asset.symbol) : "N/A", [asset.low, asset.symbol]);
  
  if (compact) {
    return (
      <div className="flex items-center justify-between bg-white h-full px-3 w-full animate-fadeIn">
        <div className="flex items-center">
          <div className="text-sm font-medium">{asset.symbol}</div>
          <div className="text-xs text-gray-500 ml-1">({asset.exchange})</div>
        </div>
        <div className="flex items-center space-x-2">
          <div className="font-medium text-sm">
            ${formattedPrice}
          </div>
          <div className={`text-xs ${isNegative ? 'text-red-500' : 'text-green-500'}`}>
            {percentChange}
          </div>
        </div>
      </div>
    );
  }
  
  // Memoize the tooltip formatter for chart
  const tooltipFormatter = useCallback((value: number) => [`$${formatPrice(value, asset.symbol)}`, 'Price'], [asset.symbol]);
  
  return (
    <div className="bg-white p-4 rounded-lg shadow-md hover:shadow-lg transition-all duration-300 border border-gray-200 h-full">
      <div className="flex justify-between items-start mb-3">
        <div>
          <h3 className="text-lg font-bold mb-1">
            <span className="mr-2">{asset.symbol}</span>
            <span className="text-sm text-gray-500">({asset.exchange})</span>
          </h3>
          <div className="mb-2">
            <span className="text-gray-500 text-sm mr-2">Volume:</span>
            <span className="text-gray-800 text-sm font-medium">${asset.volume?.toLocaleString() || "N/A"}</span>
          </div>
        </div>
        <div className="text-right">
          <div className="font-bold text-lg">
            ${formattedPrice}
          </div>
          <div className={`inline-block px-2 py-1 text-xs font-medium rounded-full ${isNegative ? 'text-red-600 bg-red-50' : 'text-green-600 bg-green-50'}`}>
            {percentChange}
          </div>
        </div>
      </div>
      
      {chartData.length > 0 && (
        <div className="h-40 mt-2">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={chartData} margin={{ top: 5, right: 5, bottom: 5, left: 5 }}>
              <YAxis 
                domain={['dataMin - 1', 'dataMax + 1']} 
                hide 
              />
              <XAxis dataKey="time" hide />
              <Tooltip 
                formatter={tooltipFormatter}
                labelFormatter={() => ''}
                contentStyle={{
                  backgroundColor: 'rgba(255, 255, 255, 0.8)',
                  border: '1px solid #ccc',
                  borderRadius: '4px'
                }}
              />
              <Line 
                type="monotone" 
                dataKey="price" 
                stroke={chartColor} 
                strokeWidth={2}
                dot={false}
                activeDot={{ r: 5 }}
                isAnimationActive={false}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
      )}
      
      <div className="flex justify-between mt-4 pt-2 border-t border-gray-100">
        <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-sm w-full">
          <div>High: <span className="font-medium">${formattedHigh}</span></div>
          <div>Low: <span className="font-medium">${formattedLow}</span></div>
        </div>
        <Link 
          to={`/exchange/${asset.exchange}/asset/${asset.symbol}`}
          className="text-blue-600 hover:text-blue-800 transition-colors text-sm flex items-center"
        >
          Details
          <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 ml-1" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M10.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L12.586 11H5a1 1 0 110-2h7.586l-2.293-2.293a1 1 0 010-1.414z" clipRule="evenodd" />
          </svg>
        </Link>
      </div>
    </div>
  );
}, (prevProps, nextProps) => {
  // Оптимизированное сравнение для memo
  // Сравниваем только важные свойства
  return (
    prevProps.asset.symbol === nextProps.asset.symbol &&
    prevProps.compact === nextProps.compact &&
    prevProps.asset.price === nextProps.asset.price &&
    prevProps.asset.priceChangePercent === nextProps.asset.priceChangePercent
  );
});

export default QuoteChart; 