import { MarketData, MarketType } from '../types';

const API_BASE_URL = 'http://localhost:8081';

// Helper function to convert API response to our frontend model
const mapResponseToMarketData = (data: any): MarketData => {
  return {
    exchange: data.Exchange || '',
    symbol: data.Symbol || '',
    market: data.Market || '',
    price: typeof data.Price === 'number' ? data.Price : undefined,
    volume: typeof data.Volume === 'number' ? data.Volume : undefined,
    high: typeof data.High === 'number' ? data.High : undefined,
    low: typeof data.Low === 'number' ? data.Low : undefined,
    priceChangePercent: data.PriceChangePercent || '',
    timestamp: data.Timestamp || '',
  };
};

/**
 * Fetches market data for a specific market type (crypto, stock, forex)
 */
export async function fetchMarketsData(marketType: MarketType): Promise<MarketData[]> {
  try {
    const endpoint = marketType === MarketType.CRYPTO 
      ? '/crypto-market' 
      : marketType === MarketType.STOCK 
        ? '/stock-market' 
        : '/forex-market';
    
    const response = await fetch(`${API_BASE_URL}${endpoint}`);
    if (!response.ok) {
      console.error(`API error: ${response.status} ${response.statusText}`);
      return []; // Return empty array on error
    }
    
    const data = await response.json();
    return Array.isArray(data) ? data.map(mapResponseToMarketData) : [];
  } catch (error) {
    console.error(`Error fetching ${marketType} market data:`, error);
    return []; // Return empty array on error
  }
}

/**
 * Fetches data for a specific exchange
 */
export async function fetchExchangeData(exchange: string): Promise<MarketData[]> {
  if (!exchange) return [];
  
  try {
    console.log(`Fetching exchange data for ${exchange}...`);
    const response = await fetch(`${API_BASE_URL}/exchange/${exchange}`);
    if (!response.ok) {
      console.error(`API error: ${response.status} ${response.statusText}`);
      return [];
    }
    
    const data = await response.json();
    return Array.isArray(data) ? data.map(mapResponseToMarketData) : [];
  } catch (error) {
    console.error(`Error fetching exchange data for ${exchange}:`, error);
    return [];
  }
}

/**
 * Fetches data for a specific asset
 */
export async function fetchAssetDetails(exchange: string, symbol: string): Promise<MarketData | null> {
  if (!exchange || !symbol) return null;
  
  try {
    const response = await fetch(`${API_BASE_URL}/exchange/${exchange}/asset/${symbol}`);
    if (!response.ok) {
      console.error(`API error: ${response.status} ${response.statusText}`);
      return null;
    }
    
    const data = await response.json();
    return mapResponseToMarketData(data);
  } catch (error) {
    console.error(`Error fetching asset details for ${exchange}/${symbol}:`, error);
    return null;
  }
} 