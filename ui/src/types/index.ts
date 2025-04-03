export interface MarketData {
  exchange: string;
  symbol: string;
  market: string;
  price?: number;
  volume?: number;
  high?: number;
  low?: number;
  priceChangePercent?: string;
  timestamp?: string;
}

export interface HomePageData {
  crypto: MarketData[];
  stock: MarketData[];
}

export interface ExchangeData {
  name: string;
  assets: MarketData[];
}

export type Market = 'crypto' | 'stock'; 

export enum MarketType {
  CRYPTO = 'crypto',
  STOCK = 'stock',
} 