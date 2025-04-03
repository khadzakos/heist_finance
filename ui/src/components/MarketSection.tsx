import React, { useState, useEffect, useRef, useCallback, useReducer } from 'react';
import { Link } from 'react-router-dom';
import { MarketData } from '../types';

interface MarketSectionProps {
  title: string;
  exchanges: Record<string, MarketData[]>;
  market: string;
}

// Определяем тип состояния
interface MarketSectionState {
  currentExchangeIndex: number;
  visibleExchanges: string[];
  columnsCount: number;
}

// Определяем типы действий
type MarketSectionAction = 
  | { type: 'NEXT_EXCHANGE' }
  | { type: 'PREV_EXCHANGE' }
  | { type: 'SET_COLUMNS', count: number }
  | { type: 'UPDATE_VISIBLE', exchanges: string[] };

// Reducer для управления состоянием
const marketSectionReducer = (state: MarketSectionState, action: MarketSectionAction): MarketSectionState => {
  switch (action.type) {
    case 'NEXT_EXCHANGE':
      return {
        ...state,
        currentExchangeIndex: state.currentExchangeIndex + 1
      };
    case 'PREV_EXCHANGE':
      return {
        ...state,
        currentExchangeIndex: state.currentExchangeIndex - 1
      };
    case 'SET_COLUMNS':
      return {
        ...state,
        columnsCount: action.count
      };
    case 'UPDATE_VISIBLE':
      return {
        ...state,
        visibleExchanges: action.exchanges
      };
    default:
      return state;
  }
};

const MarketSection: React.FC<MarketSectionProps> = ({ title, exchanges, market }) => {
  const exchangeNames = Object.keys(exchanges);
  const containerRef = useRef<HTMLDivElement>(null);
  
  // Инициализация состояния с useReducer
  const [state, dispatch] = useReducer(marketSectionReducer, {
    currentExchangeIndex: 0,
    visibleExchanges: [],
    columnsCount: 1
  });
  
  const { currentExchangeIndex, visibleExchanges, columnsCount } = state;
  
  // Обновляем видимые биржи при изменении размера или индекса
  useEffect(() => {
    // Защита от пустых данных
    if (!exchangeNames.length) return;
    
    const updateVisibleColumns = () => {
      const container = containerRef.current;
      if (!container) return;
      
      // Рассчитываем количество видимых колонок
      const containerWidth = container.offsetWidth;
      const estimatedColumnsCount = Math.max(1, Math.floor(containerWidth / 250));
      const maxVisible = Math.min(exchangeNames.length, estimatedColumnsCount);
      
      // Обновляем счетчик колонок только если он изменился
      if (maxVisible !== columnsCount) {
        dispatch({ type: 'SET_COLUMNS', count: maxVisible });
      }
      
      // Определяем какие биржи показывать
      const newVisibleExchanges: string[] = [];
      for (let i = 0; i < maxVisible; i++) {
        const index = (currentExchangeIndex + i) % exchangeNames.length;
        newVisibleExchanges.push(exchangeNames[index]);
      }
      
      // Проверяем, изменились ли видимые биржи
      const exchangesChanged = 
        newVisibleExchanges.length !== visibleExchanges.length ||
        newVisibleExchanges.some((ex, idx) => ex !== visibleExchanges[idx]);
      
      if (exchangesChanged) {
        dispatch({ type: 'UPDATE_VISIBLE', exchanges: newVisibleExchanges });
      }
    };
    
    // Первоначальное обновление
    updateVisibleColumns();
    
    // Добавляем слушатель изменения размера окна
    window.addEventListener('resize', updateVisibleColumns);
    return () => window.removeEventListener('resize', updateVisibleColumns);
  }, [exchangeNames, currentExchangeIndex, columnsCount, visibleExchanges]);
  
  // Обработчики навигации
  const handleNext = useCallback(() => {
    dispatch({ type: 'NEXT_EXCHANGE' });
  }, []);
  
  const handlePrev = useCallback(() => {
    dispatch({ type: 'PREV_EXCHANGE' });
  }, []);
  
  // Format volume for display
  const formatVolume = (volume?: number): string => {
    if (volume === undefined) return 'N/A';
    
    if (volume >= 1000000000) {
      return `$${(volume / 1000000000).toFixed(2)}B`;
    } else if (volume >= 1000000) {
      return `$${(volume / 1000000).toFixed(2)}M`;
    } else if (volume >= 1000) {
      return `$${(volume / 1000).toFixed(2)}K`;
    }
    
    return `$${volume.toLocaleString()}`;
  };
  
  // Format price for display
  const formatPrice = (price?: number): string => {
    if (price === undefined) return 'N/A';
    
    if (price < 0.01) {
      return price.toFixed(6);
    } else if (price < 1) {
      return price.toFixed(4);
    } else if (price < 1000) {
      return price.toFixed(2);
    }
    
    return price.toLocaleString(undefined, { maximumFractionDigits: 2 });
  };
  
  if (exchangeNames.length === 0) {
    return (
      <div className="bg-white p-6 rounded-lg shadow-md">
        <h2 className="text-xl font-bold mb-4">{title}</h2>
        <p className="text-gray-500">No exchanges available</p>
      </div>
    );
  }
  
  const needsNavigation = exchangeNames.length > columnsCount;
  
  return (
    <section className="mb-8">
      <div className="mb-4 flex justify-between items-center">
        <h2 className="text-xl font-bold">
          {title}
          <span className="ml-2 text-sm font-normal text-gray-500 inline-block py-1 px-2 rounded-full bg-gray-100">
            {exchangeNames.length} exchanges
          </span>
        </h2>
        
        {needsNavigation && (
          <div className="flex items-center">
            <button 
              onClick={handlePrev}
              className="p-1 rounded-full hover:bg-gray-100 mr-2"
              aria-label="Previous exchanges"
            >
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            
            <button 
              onClick={handleNext}
              className="p-1 rounded-full hover:bg-gray-100 ml-2"
              aria-label="Next exchanges"
            >
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
        )}
      </div>
      
      <div ref={containerRef} className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4">
        {visibleExchanges.map(exchange => (
          <div 
            key={exchange} 
            className="bg-white rounded-lg shadow-md flex flex-col h-full border border-gray-200 hover:shadow-lg transition-shadow duration-300"
          >
            <div className="px-3 py-2 border-b flex justify-between items-center bg-gray-50 rounded-t-lg">
              <h3 className="text-sm font-semibold text-gray-800">
                <span className="mr-2">{exchange}</span>
                <span className="text-xs font-normal text-gray-500 inline-block py-0.5 px-1.5 rounded-full bg-gray-100">
                  {exchanges[exchange].length}
                </span>
              </h3>
              <Link 
                to={`/exchange/${exchange}`}
                className="text-blue-600 hover:text-blue-800 transition-colors text-xs font-medium flex items-center rounded-md px-2 py-1 hover:bg-blue-50"
              >
                View All
                <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3 ml-1" viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M10.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L12.586 11H5a1 1 0 110-2h7.586l-2.293-2.293a1 1 0 010-1.414z" clipRule="evenodd" />
                </svg>
              </Link>
            </div>
            
            <div className="overflow-y-hidden" style={{ overflowY: 'hidden' }}>
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50 sticky top-0 z-10">
                  <tr>
                    <th className="px-3 py-2 text-left text-xs font-medium text-gray-500 tracking-wider">Symbol</th>
                    <th className="px-3 py-2 text-right text-xs font-medium text-gray-500 tracking-wider">Price</th>
                    <th className="px-3 py-2 text-right text-xs font-medium text-gray-500 tracking-wider">24h</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-100">
                  {exchanges[exchange].slice(0, 15).map(asset => (
                    <tr 
                      key={asset.symbol}
                      className="hover:bg-gray-50 transition-colors cursor-pointer"
                      onClick={() => window.location.href = `/exchange/${asset.exchange}/asset/${asset.symbol}`}
                    >
                      <td className="px-3 py-2 whitespace-nowrap">
                        <div className="font-medium text-xs">{asset.symbol}</div>
                      </td>
                      <td className="px-3 py-2 whitespace-nowrap text-right text-xs">
                        ${formatPrice(asset.price)}
                      </td>
                      <td className="px-3 py-2 whitespace-nowrap text-right">
                        <div className={`inline-block px-1.5 py-0.5 text-xs font-medium rounded-full ${
                          asset.priceChangePercent?.startsWith('-') 
                            ? 'text-red-600 bg-red-50' 
                            : 'text-green-600 bg-green-50'
                        }`}>
                          {asset.priceChangePercent || '0.00%'}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        ))}
      </div>
    </section>
  );
};

export default React.memo(MarketSection); 