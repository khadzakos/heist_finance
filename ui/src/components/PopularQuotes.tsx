import React, { useEffect, useState, useRef } from 'react';
import { MarketData, MarketType } from '../types';
import { fetchMarketsData } from '../api/marketApi';
import QuoteChart from './QuoteChart';
import { motion, AnimatePresence } from 'framer-motion';

interface PopularQuotesProps {
  quotes?: MarketData[];
  compact?: boolean;
}

const PopularQuotes: React.FC<PopularQuotesProps> = ({ 
  quotes: initialQuotes, 
  compact = false 
}) => {
  const [quotes, setQuotes] = useState<MarketData[]>([]);
  const [loading, setLoading] = useState(true);
  const [displayedQuotes, setDisplayedQuotes] = useState<MarketData[]>([]);
  const timerRef = useRef<NodeJS.Timeout | null>(null);
  const rotationTimerRef = useRef<NodeJS.Timeout | null>(null);
  
  // Get initial data if not provided as props
  useEffect(() => {
    if (initialQuotes && initialQuotes.length > 0) {
      // Filter top 9 by volume
      const sortedQuotes = [...initialQuotes]
        .filter(q => q.volume !== undefined && q.price !== undefined)
        .sort((a, b) => (b.volume || 0) - (a.volume || 0))
        .slice(0, 9);
        
      setQuotes(sortedQuotes);
      
      // Выбираем начальный набор из 3-х случайных котировок
      selectRandomQuotes(sortedQuotes);
      
      setLoading(false);
    } else {
      const fetchData = async () => {
        try {
          setLoading(true);
          // Fetch crypto and stock data
          const [cryptoData, stockData] = await Promise.all([
            fetchMarketsData(MarketType.CRYPTO),
            fetchMarketsData(MarketType.STOCK)
          ]);
          
          // Combine and sort by volume
          const allMarkets = [...cryptoData, ...stockData]
            .filter(market => market.volume !== undefined && market.price !== undefined)
            .sort((a, b) => (b.volume || 0) - (a.volume || 0))
            .slice(0, 9);  // Get top 9 by volume
            
          if (allMarkets.length === 0) {
            console.warn('No market data available, using fallback data');
            // Fallback data if API returns no data
            const fallbackData = generateFallbackData();
            setQuotes(fallbackData);
            selectRandomQuotes(fallbackData);
          } else {
            setQuotes(allMarkets);
            selectRandomQuotes(allMarkets);
          }
        } catch (error) {
          console.error('Failed to fetch popular quotes:', error);
          // Use fallback data on error
          const fallbackData = generateFallbackData();
          setQuotes(fallbackData);
          selectRandomQuotes(fallbackData);
        } finally {
          setLoading(false);
        }
      };
      
      fetchData();
      
      // Refresh data every 20 seconds to keep it fresh
      timerRef.current = setInterval(fetchData, 20000);
      return () => {
        if (timerRef.current) clearInterval(timerRef.current);
      };
    }
  }, [initialQuotes]);
  
  // Выбор случайных котировок из доступных
  const selectRandomQuotes = (availableQuotes: MarketData[]) => {
    if (availableQuotes.length <= 3) {
      setDisplayedQuotes(availableQuotes);
      return;
    }
    
    // Создаем копию массива
    const quotesCopy = [...availableQuotes];
    const selected: MarketData[] = [];
    
    // Выбираем 3 случайные котировки
    for (let i = 0; i < 3; i++) {
      if (quotesCopy.length === 0) break;
      
      // Выбираем случайный индекс
      const randomIndex = Math.floor(Math.random() * quotesCopy.length);
      
      // Добавляем выбранную котировку в результат
      selected.push(quotesCopy[randomIndex]);
      
      // Удаляем выбранную котировку из копии, чтобы избежать дубликатов
      quotesCopy.splice(randomIndex, 1);
    }
    
    setDisplayedQuotes(selected);
  };
  
  // Ротация котировок каждые 5 секунд
  useEffect(() => {
    if (quotes.length <= 3) return; // Не нужно ротировать, если у нас всего 3 или меньше котировок
    
    // Устанавливаем интервал ротации в 5 секунд
    rotationTimerRef.current = setInterval(() => {
      selectRandomQuotes(quotes);
    }, 5000);
    
    return () => {
      if (rotationTimerRef.current) clearInterval(rotationTimerRef.current);
    };
  }, [quotes]);
  
  // Generate fallback data in case API fails
  const generateFallbackData = (): MarketData[] => {
    const symbols = ['BTC', 'ETH', 'SOL'];
    const exchanges = ['Binance', 'Coinbase', 'Kraken'];
    const result: MarketData[] = [];
    
    for (let i = 0; i < 3; i++) {
      const price = i === 0 ? 45000 : i === 1 ? 2500 : 100;
      result.push({
        exchange: exchanges[i],
        symbol: symbols[i],
        market: 'crypto',
        price: price,
        priceChangePercent: `${(Math.random() * 10 - 5).toFixed(3)}%`,
        volume: 10000000 * (3 - i), 
        high: price * 1.05,
        low: price * 0.95,
        timestamp: new Date().toISOString()
      });
    }
    
    return result;
  };
  
  // Loading indicator
  if (loading) {
    return (
      <div className="flex justify-center items-center h-full">
        <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }
  
  if (displayedQuotes.length === 0) {
    return <div className="p-4 text-center text-gray-500">No quotes available</div>;
  }

  // For compact view (usually in sidebar or smaller containers)
  if (compact) {
    return (
      <div className="space-y-2">
        {displayedQuotes.map((quote, index) => (
          <QuoteChart
            key={`${quote.exchange}:${quote.symbol}`}
            asset={quote}
            index={index}
            compact={true}
          />
        ))}
      </div>
    );
  }
  
  // Full dashboard view
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 h-full">
      {displayedQuotes.map((quote, index) => (
        <AnimatePresence key={`animate-${index}`} mode="wait">
          <motion.div
            key={`${quote.exchange}:${quote.symbol}`}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ 
              duration: 0.4, 
              delay: index * 0.2, // Увеличенная задержка между карточками (было 0.05)
              ease: "easeInOut" 
            }}
            className="h-full"
          >
            <QuoteChart
              asset={quote}
              index={index}
              compact={false}
            />
          </motion.div>
        </AnimatePresence>
      ))}
    </div>
  );
};

export default PopularQuotes; 