import React, { useEffect, useState, useCallback, useRef, memo, useMemo } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { fetchExchangeData } from '../services/api';
import { MarketData } from '../types';

// Статичный компонент для отображения цены с обновляемым значением
const StaticPriceDisplay = memo(({ 
  price, 
  direction
}: { 
  price: number | undefined; 
  direction: 'up' | 'down' | null;
}) => {
  // CSS класс на основе направления изменения цены
  const baseClass = "transition-colors duration-300 font-medium";
  const directionClass = direction === 'up' 
    ? ' text-green-600' 
    : direction === 'down' 
      ? ' text-red-600' 
      : ' text-gray-900';
  
  return (
    <span className={baseClass + directionClass}>
      ${price?.toFixed(3) || "N/A"}
    </span>
  );
});

// Статичный компонент для отображения процентного изменения
const StaticPercentageDisplay = memo(({
  percentChange,
  isNegative
}: {
  percentChange: string | undefined;
  isNegative: boolean;
}) => {
  const className = `inline-block px-2 py-1 text-xs font-medium rounded-full ${
    isNegative 
      ? 'text-red-600 bg-red-50' 
      : 'text-green-600 bg-green-50'
  }`;
  
  return (
    <span className={className}>
      {percentChange || "N/A"}
    </span>
  );
});

// Строка таблицы активов - используем memo с функцией сравнения для оптимизации рендеринга
const AssetRow = memo(({ 
  asset, 
  direction,
  isNegativeChange,
  onClick 
}: { 
  asset: MarketData; 
  direction: 'up' | 'down' | null;
  isNegativeChange: boolean;
  onClick: () => void;
}) => {
  return (
    <tr 
      className="hover:bg-gray-50 cursor-pointer transition-colors border-b border-gray-100"
      onClick={onClick}
    >
      <td className="px-6 py-4 whitespace-nowrap">
        <div className="font-medium text-gray-900">{asset.symbol}</div>
      </td>
      <td className="px-6 py-4 whitespace-nowrap">
        <StaticPriceDisplay 
          price={asset.price} 
          direction={direction}
        />
      </td>
      <td className="px-6 py-4 whitespace-nowrap text-gray-700">
        {asset.volume 
          ? `$${asset.volume >= 1000000 
              ? `${(asset.volume / 1000000).toFixed(2)}M` 
              : asset.volume.toLocaleString()}`
          : "N/A"}
      </td>
      <td className="px-6 py-4 whitespace-nowrap text-gray-700">${asset.high?.toFixed(3) || "N/A"}</td>
      <td className="px-6 py-4 whitespace-nowrap text-gray-700">${asset.low?.toFixed(3) || "N/A"}</td>
      <td className="px-6 py-4 whitespace-nowrap">
        <StaticPercentageDisplay 
          percentChange={asset.priceChangePercent} 
          isNegative={isNegativeChange}
        />
      </td>
    </tr>
  );
}, (prevProps, nextProps) => {
  // Оптимизация рендеринга: перерисовываем только при изменении важных свойств
  // Если символ и биржа не изменились и это тот же самый актив - считаем что это один и тот же актив
  const sameAsset = prevProps.asset.symbol === nextProps.asset.symbol && 
                    prevProps.asset.exchange === nextProps.asset.exchange;
                    
  if (!sameAsset) return false; // Если активы разные, нужно перерисовать
  
  // Проверяем изменились ли значимые данные
  const priceChanged = prevProps.asset.price !== nextProps.asset.price;
  const directionChanged = prevProps.direction !== nextProps.direction;
  const percentChanged = prevProps.asset.priceChangePercent !== nextProps.asset.priceChangePercent;
  const negativeChanged = prevProps.isNegativeChange !== nextProps.isNegativeChange;
  
  // Если ничего не изменилось, не перерисовываем
  return !(priceChanged || directionChanged || percentChanged || negativeChanged);
});

// Заголовок таблицы
const TableHeader = memo(({ 
  onSort, 
  sortField, 
  sortDirection 
}: {
  onSort: (field: keyof MarketData) => void;
  sortField: keyof MarketData;
  sortDirection: 'asc' | 'desc';
}) => {
  // Компонент для отображения иконки сортировки
  const SortIcon = ({ field }: { field: keyof MarketData }) => {
    if (field !== sortField) return null;
    return (
      <span className="ml-1">
        {sortDirection === 'asc' ? '▲' : '▼'}
      </span>
    );
  };

  const headerClasses = "px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer transition-colors hover:bg-gray-100";

  return (
    <thead className="bg-gray-50">
      <tr>
        <th 
          className={headerClasses}
          onClick={() => onSort('symbol')}
        >
          Symbol <SortIcon field="symbol" />
        </th>
        <th 
          className={headerClasses}
          onClick={() => onSort('price')}
        >
          Price <SortIcon field="price" />
        </th>
        <th 
          className={headerClasses}
          onClick={() => onSort('volume')}
        >
          24h Volume <SortIcon field="volume" />
        </th>
        <th 
          className={headerClasses}
          onClick={() => onSort('high')}
        >
          24h High <SortIcon field="high" />
        </th>
        <th 
          className={headerClasses}
          onClick={() => onSort('low')}
        >
          24h Low <SortIcon field="low" />
        </th>
        <th 
          className={headerClasses}
          onClick={() => onSort('priceChangePercent')}
        >
          Change (24h) <SortIcon field="priceChangePercent" />
        </th>
      </tr>
    </thead>
  );
});

// Таблица активов - вынесена в отдельный мемоизированный компонент
const AssetTable = memo(({
  sortedData,
  priceDirections,
  handleSort,
  sortField,
  sortDirection,
  handleAssetClick,
  isNegativeChange
}: {
  sortedData: MarketData[];
  priceDirections: Record<string, 'up' | 'down' | null>;
  handleSort: (field: keyof MarketData) => void;
  sortField: keyof MarketData;
  sortDirection: 'asc' | 'desc';
  handleAssetClick: (symbol: string) => void;
  isNegativeChange: (asset: MarketData) => boolean;
}) => {
  if (!sortedData.length) {
    return (
      <div className="bg-white shadow-md rounded-lg p-6 text-center text-gray-500">
        No assets available for this exchange
      </div>
    );
  }
  
  return (
    <div className="bg-white shadow-md rounded-lg overflow-hidden border border-gray-200">
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <TableHeader 
            onSort={handleSort} 
            sortField={sortField} 
            sortDirection={sortDirection} 
          />
          <tbody className="bg-white divide-y divide-gray-100">
            {sortedData.map((asset) => (
              <AssetRow 
                key={`${asset.exchange}-${asset.symbol}`}
                asset={asset}
                direction={priceDirections[asset.symbol] || null}
                isNegativeChange={isNegativeChange(asset)}
                onClick={() => handleAssetClick(asset.symbol)}
              />
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
});

// Компонент обновления времени - выносим его отдельно, чтобы не перерисовывать всю страницу
const LastUpdatedDisplay = memo(({ lastUpdated }: { lastUpdated: Date }) => {
  return (
    <p className="text-xs text-gray-500 mt-1">Last updated: {lastUpdated.toLocaleTimeString()}</p>
  );
});

const ExchangePage: React.FC = () => {
  const { exchange } = useParams<{ exchange: string }>();
  const navigate = useNavigate();
  const [data, setData] = useState<MarketData[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [sortField, setSortField] = useState<keyof MarketData>('symbol');
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc');
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date());
  
  // Храним предыдущие цены для определения направления изменения
  const prevPricesRef = useRef<Record<string, number>>({});
  // Направления изменения цен для визуальных индикаторов
  const [priceDirections, setPriceDirections] = useState<Record<string, 'up' | 'down' | null>>({});
  
  // Флаг монтирования компонента
  const mountedRef = useRef(true);
  // Идентификатор интервала обновления
  const intervalIdRef = useRef<NodeJS.Timeout | null>(null);
  // Идентификатор таймаута сброса индикаторов
  const timeoutIdRef = useRef<NodeJS.Timeout | null>(null);
  // Флаг, показывающий, что запрос к API в процессе
  const isLoadingRef = useRef(false);
  // Стабильные ссылки на данные
  const dataRef = useRef<MarketData[]>([]);
  
  // Проверка отрицательного изменения цены - мемоизируем функцию
  const isNegativeChange = useCallback((asset: MarketData): boolean => {
    if (!asset.priceChangePercent) return false;
    return asset.priceChangePercent.toString().startsWith('-');
  }, []);
  
  // Получение числового значения из строки процентного изменения - мемоизируем функцию
  const getPercentValue = useCallback((asset: MarketData): number => {
    if (!asset.priceChangePercent) return 0;
    return parseFloat(asset.priceChangePercent.toString().replace('%', '').replace('-', ''));
  }, []);
  
  // Обработчик клика по активу - мемоизируем функцию
  const handleAssetClick = useCallback((symbol: string) => {
    if (!exchange) return;
    navigate(`/exchange/${exchange}/asset/${symbol}`);
  }, [exchange, navigate]);
  
  // Обработчик сортировки таблицы - мемоизируем функцию
  const handleSort = useCallback((field: keyof MarketData) => {
    setSortField(prev => {
      if (prev === field) {
        setSortDirection(prevDir => prevDir === 'asc' ? 'desc' : 'asc');
        return prev;
      } else {
        setSortDirection('asc');
        return field;
      }
    });
  }, []);
  
  // Функция сортировки данных - мемоизируем результат
  const sortedData = useMemo(() => {
    if (!data.length) return [];
    
    return [...data].sort((a, b) => {
      let valueA = a[sortField];
      let valueB = b[sortField];
      
      // Особая обработка для процентного изменения
      if (sortField === 'priceChangePercent') {
        valueA = getPercentValue(a) * (isNegativeChange(a) ? -1 : 1);
        valueB = getPercentValue(b) * (isNegativeChange(b) ? -1 : 1);
      }
      
      // Безопасное сравнение
      if (valueA === undefined && valueB === undefined) return 0;
      if (valueA === undefined) return 1;
      if (valueB === undefined) return -1;
      
      // Сравнение строк и чисел
      if (typeof valueA === 'string' && typeof valueB === 'string') {
        return sortDirection === 'asc' 
          ? valueA.localeCompare(valueB) 
          : valueB.localeCompare(valueA);
      }
      
      if (valueA < valueB) return sortDirection === 'asc' ? -1 : 1;
      if (valueA > valueB) return sortDirection === 'asc' ? 1 : -1;
      return 0;
    });
  }, [data, sortField, sortDirection, getPercentValue, isNegativeChange]);
  
  // Обновляем только цены, не перерисовывая компоненты
  const updatePricesOnly = useCallback((newData: MarketData[]) => {
    // Создаем карту для быстрого поиска по символу
    const dataMap = new Map<string, MarketData>();
    newData.forEach(asset => {
      dataMap.set(asset.symbol, asset);
    });
    
    // Обновляем только цены в существующих данных
    const updatedData = dataRef.current.map(oldAsset => {
      const newAsset = dataMap.get(oldAsset.symbol);
      if (!newAsset) return oldAsset; // Если нет новых данных, используем старые
      
      return {
        ...oldAsset,
        price: newAsset.price,
        high: newAsset.high,
        low: newAsset.low,
        volume: newAsset.volume,
        priceChangePercent: newAsset.priceChangePercent
      };
    });
    
    // Добавляем новые активы, которых не было раньше
    newData.forEach(newAsset => {
      const exists = updatedData.some(a => a.symbol === newAsset.symbol);
      if (!exists) {
        updatedData.push(newAsset);
      }
    });
    
    return updatedData;
  }, []);
  
  // Функция для загрузки данных и обновления цен
  const fetchData = useCallback(async () => {
    if (!exchange || !mountedRef.current || isLoadingRef.current) return;
    
    // Устанавливаем флаг загрузки, чтобы избежать параллельных запросов
    isLoadingRef.current = true;
    
    try {
      const result = await fetchExchangeData(exchange);
      
      if (!mountedRef.current) {
        isLoadingRef.current = false;
        return;
      }
      
      // Проверяем, получили ли мы какие-то валидные данные
      if (!Array.isArray(result) || result.length === 0) {
        console.warn('Received empty or invalid data for', exchange);
        // Не обновляем state если получили пустые данные
        isLoadingRef.current = false;
        return;
      }
      
      // Обновляем индикаторы изменения цен
      const newDirections: Record<string, 'up' | 'down' | null> = {};
      const newPrices: Record<string, number> = {};
      
      result.forEach(asset => {
        if (asset.price !== undefined) {
          const symbol = asset.symbol;
          newPrices[symbol] = asset.price;
          
          // Определяем направление изменения цены
          if (prevPricesRef.current[symbol] !== undefined) {
            if (asset.price > prevPricesRef.current[symbol]) {
              newDirections[symbol] = 'up';
            } else if (asset.price < prevPricesRef.current[symbol]) {
              newDirections[symbol] = 'down';
            }
          }
        }
      });
      
      // Устанавливаем данные и отключаем состояние загрузки
      setData(prevData => {
        // Если это первая загрузка, просто используем новые данные
        if (prevData.length === 0) {
          dataRef.current = result;
          return result;
        } else {
          // Иначе обновляем только цены, сохраняя существующие объекты
          const updatedData = updatePricesOnly(result);
          dataRef.current = updatedData;
          return updatedData;
        }
      });
      
      // Обновляем индикаторы направления цен
      setPriceDirections(newDirections);
      
      // Сохраняем текущие цены для следующего сравнения
      prevPricesRef.current = newPrices;
      
      // Обновляем время последнего обновления
      setLastUpdated(new Date());
      
      // Отключаем состояние загрузки
      setLoading(false);
      
      // Отключаем индикаторы через 1 секунду
      if (timeoutIdRef.current) {
        clearTimeout(timeoutIdRef.current);
      }
      
      timeoutIdRef.current = setTimeout(() => {
        if (mountedRef.current) {
          setPriceDirections({});
        }
      }, 1000);
      
    } catch (err) {
      console.error("Error fetching exchange data:", err);
      if (mountedRef.current) {
        setError(`Failed to fetch data for exchange: ${exchange}`);
        setLoading(false);
      }
    } finally {
      // Сбрасываем флаг загрузки в любом случае
      isLoadingRef.current = false;
    }
  }, [exchange, updatePricesOnly]);
  
  // Настройка загрузки данных при монтировании компонента
  useEffect(() => {
    // Устанавливаем флаги
    mountedRef.current = true;
    isLoadingRef.current = false;
    
    // Очищаем предыдущие интервалы и таймауты
    if (intervalIdRef.current) {
      clearInterval(intervalIdRef.current);
      intervalIdRef.current = null;
    }
    
    if (timeoutIdRef.current) {
      clearTimeout(timeoutIdRef.current);
      timeoutIdRef.current = null;
    }
    
    // Сбрасываем состояние для нового обмена
    setData([]);
    dataRef.current = [];
    setPriceDirections({});
    prevPricesRef.current = {};
    
    // Начальная загрузка
    setLoading(true);
    fetchData();
    
    // Настраиваем интервал обновления данных каждые 10 секунд через API
    intervalIdRef.current = setInterval(() => {
      if (mountedRef.current && !isLoadingRef.current) {
        fetchData();
      }
    }, 10000);
    
    // Настраиваем частые мелкие обновления с симуляцией изменения цен
    const simulationIntervalId = setInterval(() => {
      if (!mountedRef.current || dataRef.current.length === 0) return;
      
      setData(prevData => {
        if (prevData.length === 0) return prevData;
        
        // Создаем новые направления для анимации
        const newDirections: Record<string, 'up' | 'down' | null> = {};
        const newPrices: Record<string, number> = {};
        
        // Создаем новые данные с имитацией изменения цен
        const updatedData = prevData.map(asset => {
          if (asset.price === undefined) return asset;
          
          const symbol = asset.symbol;
          const currentPrice = asset.price;
          const oldPrice = prevPricesRef.current[symbol] || currentPrice;
          
          // Применяем случайное изменение цены (-0.2% до +0.2%)
          const movement = (Math.random() - 0.5) * 0.004 * currentPrice;
          const newPrice = currentPrice + movement;
          
          // Сохраняем новую цену и направление изменения
          newPrices[symbol] = newPrice;
          
          if (newPrice > oldPrice) {
            newDirections[symbol] = 'up';
          } else if (newPrice < oldPrice) {
            newDirections[symbol] = 'down';
          }
          
          // Пересчитываем процент изменения, если было первоначальное значение
          let newPercentChange = asset.priceChangePercent;
          if (asset.priceChangePercent) {
            const originalPercent = parseFloat(asset.priceChangePercent.toString().replace('%', ''));
            const percentAdjustment = (movement / currentPrice) * 100;
            const newPercent = originalPercent + percentAdjustment;
            newPercentChange = `${newPercent.toFixed(3)}%`;
          }
          
          // Возвращаем обновленный актив
          return {
            ...asset,
            price: newPrice,
            priceChangePercent: newPercentChange
          };
        });
        
        // Обновляем индикаторы направления и сохраняем новые цены
        setPriceDirections(newDirections);
        prevPricesRef.current = newPrices;
        
        // Обновляем время
        setLastUpdated(new Date());
        
        // Сбрасываем направления через 1 секунду
        if (timeoutIdRef.current) {
          clearTimeout(timeoutIdRef.current);
        }
        
        timeoutIdRef.current = setTimeout(() => {
          if (mountedRef.current) {
            setPriceDirections({});
          }
        }, 1000);
        
        // Обновляем ссылку на данные
        dataRef.current = updatedData;
        return updatedData;
      });
    }, 2000); // Обновляем каждые 2 секунды
    
    // Очистка при размонтировании
    return () => {
      mountedRef.current = false;
      
      if (intervalIdRef.current) {
        clearInterval(intervalIdRef.current);
        intervalIdRef.current = null;
      }
      
      if (timeoutIdRef.current) {
        clearTimeout(timeoutIdRef.current);
        timeoutIdRef.current = null;
      }
      
      clearInterval(simulationIntervalId);
    };
  }, [exchange, fetchData]);
  
  // Показываем загрузку
  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }
  
  // Показываем ошибку
  if (error) {
    return (
      <div className="p-4 bg-red-100 text-red-700 rounded-md">
        <p>{error}</p>
        <button 
          onClick={() => {setError(null); setLoading(true); fetchData();}}
          className="mt-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Retry
        </button>
      </div>
    );
  }
  
  // Проверяем наличие параметра биржи
  if (!exchange) {
    return <div>Exchange not specified</div>;
  }
  
  // Основной рендеринг страницы
  return (
    <div className="p-6 max-w-7xl mx-auto">
      <div className="mb-6 flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">{exchange}</h1>
          <p className="text-gray-600">{data.length} assets available</p>
        </div>
        <div className="flex flex-col items-end">
          <Link to="/" className="px-4 py-2 bg-gray-200 rounded-md hover:bg-gray-300 transition-colors">
            Back to Dashboard
          </Link>
          <LastUpdatedDisplay lastUpdated={lastUpdated} />
        </div>
      </div>

      <AssetTable
        sortedData={sortedData}
        priceDirections={priceDirections}
        handleSort={handleSort}
        sortField={sortField}
        sortDirection={sortDirection}
        handleAssetClick={handleAssetClick}
        isNegativeChange={isNegativeChange}
      />
    </div>
  );
};

// Экспортируем компонент с мемоизацией для предотвращения ненужных рендеров
export default memo(ExchangePage); 