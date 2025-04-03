import { HomePageData, MarketData } from '../types';

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

export async function fetchHomePageData(): Promise<HomePageData> {
  try {
    const response = await fetch(`${API_BASE_URL}/`);
    if (!response.ok) {
      console.error(`API error: ${response.status} ${response.statusText}`);
      // Return empty data instead of throwing to prevent app from crashing
      return { crypto: [], stock: [] };
    }
    
    const data = await response.json();
    
    // Ensure we have both crypto and stock data, even if one is missing
    return {
      crypto: Array.isArray(data.crypto) ? data.crypto.map(mapResponseToMarketData) : [],
      stock: Array.isArray(data.stock) ? data.stock.map(mapResponseToMarketData) : [],
    };
  } catch (error) {
    console.error('Error fetching home page data:', error);
    return { crypto: [], stock: [] };
  }
}

export async function fetchCryptoMarketData(): Promise<MarketData[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/crypto-market`);
    if (!response.ok) {
      console.error(`API error: ${response.status} ${response.statusText}`);
      return [];
    }
    
    const data = await response.json();
    return Array.isArray(data) ? data.map(mapResponseToMarketData) : [];
  } catch (error) {
    console.error('Error fetching crypto market data:', error);
    return [];
  }
}

export async function fetchStockMarketData(): Promise<MarketData[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/stock-market`);
    if (!response.ok) {
      console.error(`API error: ${response.status} ${response.statusText}`);
      return [];
    }
    
    const data = await response.json();
    return Array.isArray(data) ? data.map(mapResponseToMarketData) : [];
  } catch (error) {
    console.error('Error fetching stock market data:', error);
    return [];
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
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    return Array.isArray(data) ? data.map(mapResponseToMarketData) : [];
  } catch (error) {
    console.error(`Error fetching exchange data for ${exchange}:`, error);
    console.log('Falling back to mock data');
    return [];
  }
}

export async function fetchAssetDetails(exchange: string, symbol: string): Promise<MarketData> {
  if (!exchange || !symbol) {
    return { exchange: exchange || '', symbol: symbol || '', market: '' };
  }
  
  try {
    const response = await fetch(`${API_BASE_URL}/exchange/${exchange}/asset/${symbol}`);
    if (!response.ok) {
      console.error(`API error: ${response.status} ${response.statusText}`);
      return { exchange, symbol, market: '' };
    }
    
    const data = await response.json();
    return mapResponseToMarketData(data);
  } catch (error) {
    console.error(`Error fetching asset details for ${symbol}:`, error);
    return { exchange, symbol, market: '' };
  }
} 