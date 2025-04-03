import React, { useState, useEffect, useRef, useCallback, useReducer, memo } from 'react';
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

// Компонент для отображения цены с анимацией изменения
const PriceDisplay = memo(({ price, direction }: { price?: number; direction: 'up' | 'down' | null }) => {
  // Форматирование цены
  const formattedPrice = formatPrice(price);
  
  // CSS классы для анимации
  const baseClass = "text-right text-xs whitespace-nowrap transition-colors duration-300";
  const directionClass = 
    direction === 'up' ? ' text-green-600' :
    direction === 'down' ? ' text-red-600' : 
    ' text-gray-700';
  
  return (
    <div className={baseClass + directionClass}>
      ${formattedPrice}
    </div>
  );
});

// Компонент для отображения процентного изменения
const PercentageDisplay = memo(({ 
  percentChange, 
  isNegative, 
  direction 
}: { 
  percentChange: string | undefined; 
  isNegative: boolean;
  direction: 'up' | 'down' | null;
}) => {
  // CSS классы для цвета и анимации
  const baseClass = "inline-block px-1.5 py-0.5 text-xs font-medium rounded-full transition-colors duration-300";
  const colorClass = isNegative ? ' text-red-600 bg-red-50' : ' text-green-600 bg-green-50';
  const pulseClass = direction ? (direction === 'up' ? ' scale-105' : ' scale-95') : '';
  
  return (
    <div className={`${baseClass}${colorClass}${pulseClass}`}>
      {percentChange || '0.00%'}
    </div>
  );
});

// Компонент строки актива с оптимизацией рендеринга
const AssetRow = memo(({ 
  asset, 
  priceDirection, 
  formattedPrice,
  onClick 
}: { 
  asset: MarketData; 
  priceDirection: 'up' | 'down' | null;
  formattedPrice: string;
  onClick: () => void;
}) => {
  const isNegative = asset.priceChangePercent?.startsWith('-') ?? false;
  
  return (
    <tr 
      className="hover:bg-gray-50 transition-colors cursor-pointer"
      onClick={onClick}
    >
      <td className="px-3 py-2 whitespace-nowrap">
        <div className="font-medium text-xs">{asset.symbol}</div>
      </td>
      <td className="px-3 py-2 whitespace-nowrap text-right text-xs">
        <PriceDisplay 
          price={asset.price} 
          direction={priceDirection} 
        />
      </td>
      <td className="px-3 py-2 whitespace-nowrap text-right">
        <PercentageDisplay 
          percentChange={asset.priceChangePercent} 
          isNegative={isNegative}
          direction={priceDirection}
        />
      </td>
    </tr>
  );
}, (prevProps, nextProps) => {
  // Проверяем, нужно ли перерисовывать компонент
  const sameAsset = prevProps.asset.symbol === nextProps.asset.symbol &&
                    prevProps.asset.exchange === nextProps.asset.exchange;
  if (!sameAsset) return false;
  
  const priceChanged = prevProps.asset.price !== nextProps.asset.price;
  const percentChanged = prevProps.asset.priceChangePercent !== nextProps.asset.priceChangePercent;
  const directionChanged = prevProps.priceDirection !== nextProps.priceDirection;
  
  return !(priceChanged || percentChanged || directionChanged);
});

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

const MarketSection: React.FC<MarketSectionProps> = ({ title, exchanges, market }) => {
  const exchangeNames = Object.keys(exchanges);
  const containerRef = useRef<HTMLDivElement>(null);
  
  // Состояние для симуляции обновлений цен
  const [simulatedExchanges, setSimulatedExchanges] = useState<Record<string, MarketData[]>>(exchanges);
  const [priceDirections, setPriceDirections] = useState<Record<string, Record<string, 'up' | 'down' | null>>>({});
  
  // Сохраняем предыдущие цены для определения направления изменения
  const prevPricesRef = useRef<Record<string, Record<string, number>>>({});
  const mountedRef = useRef(true);
  const originalExchangesRef = useRef(exchanges);
  
  // Мемоизированный обработчик навигации к активу
  const handleAssetClick = useCallback((exchange: string, symbol: string) => {
    window.location.href = `/exchange/${exchange}/asset/${symbol}`;
  }, []);
  
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
  
  // Эффект для инициализации и очистки состояния
  useEffect(() => {
    mountedRef.current = true;
    originalExchangesRef.current = exchanges;
    
    // Инициализируем начальные цены для сравнения
    const initialPrices: Record<string, Record<string, number>> = {};
    Object.keys(exchanges).forEach(exchange => {
      initialPrices[exchange] = {};
      exchanges[exchange].forEach(asset => {
        if (asset.price !== undefined) {
          initialPrices[exchange][asset.symbol] = asset.price;
        }
      });
    });
    
    prevPricesRef.current = initialPrices;
    setSimulatedExchanges(exchanges);
    
    return () => {
      mountedRef.current = false;
    };
  }, []);
  
  // Эффект для обновления при изменении внешних данных
  useEffect(() => {
    // Обновляем только если данные действительно изменились
    if (JSON.stringify(Object.keys(exchanges)) !== JSON.stringify(Object.keys(originalExchangesRef.current))) {
      originalExchangesRef.current = exchanges;
      setSimulatedExchanges(exchanges);
      
      // Сбрасываем направления цен
      setPriceDirections({});
      
      // Обновляем начальные цены
      const newPrices: Record<string, Record<string, number>> = {};
      Object.keys(exchanges).forEach(exchange => {
        newPrices[exchange] = {};
        exchanges[exchange].forEach(asset => {
          if (asset.price !== undefined) {
            newPrices[exchange][asset.symbol] = asset.price;
          }
        });
      });
      
      prevPricesRef.current = newPrices;
    } else {
      // Если только обновились данные внутри бирж, обновляем с сохранением анимации
      const newExchanges: Record<string, MarketData[]> = {};
      const newDirections: Record<string, Record<string, 'up' | 'down' | null>> = {};
      const newPrices: Record<string, Record<string, number>> = {};
      
      Object.keys(exchanges).forEach(exchange => {
        newExchanges[exchange] = [];
        newDirections[exchange] = {};
        newPrices[exchange] = {};
        
        exchanges[exchange].forEach(asset => {
          const prevPrice = prevPricesRef.current[exchange]?.[asset.symbol];
          
          if (prevPrice !== undefined && asset.price !== undefined) {
            if (asset.price > prevPrice) {
              newDirections[exchange][asset.symbol] = 'up';
            } else if (asset.price < prevPrice) {
              newDirections[exchange][asset.symbol] = 'down';
            }
          }
          
          if (asset.price !== undefined) {
            newPrices[exchange][asset.symbol] = asset.price;
          }
          
          newExchanges[exchange].push(asset);
        });
      });
      
      setSimulatedExchanges(newExchanges);
      setPriceDirections(newDirections);
      prevPricesRef.current = newPrices;
      
      // Сбрасываем индикаторы через 1 секунду
      setTimeout(() => {
        if (mountedRef.current) {
          setPriceDirections({});
        }
      }, 1000);
    }
  }, [exchanges]);
  
  // Симуляция обновления цен каждые 2.5 секунды
  useEffect(() => {
    if (Object.keys(simulatedExchanges).length === 0) return;
    
    const intervalId = setInterval(() => {
      if (!mountedRef.current) return;
      
      const newExchanges: Record<string, MarketData[]> = {};
      const newDirections: Record<string, Record<string, 'up' | 'down' | null>> = {};
      const newPrices: Record<string, Record<string, number>> = {};
      
      // Перебираем все биржи
      Object.keys(simulatedExchanges).forEach(exchange => {
        newExchanges[exchange] = [];
        newDirections[exchange] = {};
        newPrices[exchange] = {};
        
        // Обновляем активы внутри каждой биржи
        simulatedExchanges[exchange].forEach(asset => {
          const currentPrice = asset.price;
          
          if (currentPrice !== undefined) {
            // Генерируем небольшое случайное изменение цены (-0.3% до +0.3%)
            const change = (Math.random() - 0.5) * 0.006 * currentPrice;
            const newPrice = currentPrice + change;
            
            // Определяем направление изменения для анимации
            const prevPrice = prevPricesRef.current[exchange]?.[asset.symbol];
            if (prevPrice !== undefined) {
              if (newPrice > prevPrice) {
                newDirections[exchange][asset.symbol] = 'up';
              } else if (newPrice < prevPrice) {
                newDirections[exchange][asset.symbol] = 'down';
              }
            }
            
            // Сохраняем новую цену для следующего сравнения
            newPrices[exchange][asset.symbol] = newPrice;
            
            // Обновляем процент изменения
            let newPercentChange = asset.priceChangePercent;
            if (asset.priceChangePercent) {
              const originalPercent = parseFloat(asset.priceChangePercent.toString().replace('%', ''));
              const percentAdjustment = (change / currentPrice) * 100;
              const newPercent = originalPercent + percentAdjustment;
              newPercentChange = `${newPercent.toFixed(2)}%`;
            }
            
            // Создаем обновленный актив
            newExchanges[exchange].push({
              ...asset,
              price: newPrice,
              priceChangePercent: newPercentChange
            });
          } else {
            // Если цены нет, просто копируем актив
            newExchanges[exchange].push({ ...asset });
          }
        });
      });
      
      // Обновляем состояние
      setSimulatedExchanges(newExchanges);
      setPriceDirections(newDirections);
      prevPricesRef.current = newPrices;
      
      // Сбрасываем индикаторы через 1 секунду
      setTimeout(() => {
        if (mountedRef.current) {
          setPriceDirections({});
        }
      }, 1000);
    }, 2500); // Обновляем каждые 2.5 секунды
    
    return () => clearInterval(intervalId);
  }, [simulatedExchanges]);
  
  // Обработчики навигации
  const handleNext = useCallback(() => {
    dispatch({ type: 'NEXT_EXCHANGE' });
  }, []);
  
  const handlePrev = useCallback(() => {
    dispatch({ type: 'PREV_EXCHANGE' });
  }, []);
  
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
                  {simulatedExchanges[exchange]?.length || 0}
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
                  {simulatedExchanges[exchange]?.slice(0, 15).map(asset => (
                    <AssetRow 
                      key={`${exchange}-${asset.symbol}`}
                      asset={asset}
                      priceDirection={priceDirections[exchange]?.[asset.symbol] || null}
                      formattedPrice={formatPrice(asset.price)}
                      onClick={() => handleAssetClick(exchange, asset.symbol)}
                    />
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

export default memo(MarketSection); 