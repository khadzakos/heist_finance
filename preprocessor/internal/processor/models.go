package processor

type GenericMessage interface{}

type BinanceMarketData struct {
	Event                       string `json:"e"`
	EventTime                   int64  `json:"E"`
	Symbol                      string `json:"s"`
	PriceChange                 string `json:"p"`
	PriceChangePercent          string `json:"P"`
	WeightedAveragePrice        string `json:"w"`
	FirstTradePrice             string `json:"x"`
	LastPrice                   string `json:"c"`
	LastQuantity                string `json:"Q"`
	BestBidPrice                string `json:"b"`
	BestBidQuantity             string `json:"B"`
	BestAskPrice                string `json:"a"`
	BestAskQuantity             string `json:"A"`
	OpenPrice                   string `json:"o"`
	HighPrice                   string `json:"h"`
	LowPrice                    string `json:"l"`
	TotalTradedBaseAssetVolume  string `json:"v"`
	TotalTradedQuoteAssetVolume string `json:"q"`
	StatisticsOpenTime          int64  `json:"O"`
	StatisticsCloseTime         int64  `json:"C"`
	FirstTradeID                int64  `json:"F"`
	LastTradeID                 int64  `json:"L"`
	TradeCount                  int64  `json:"n"`
}

type BybitMarketData struct {
	Symbol        string `json:"symbol"`
	LastPrice     string `json:"lastPrice"`
	HighPrice24h  string `json:"highPrice24h"`
	LowPrice24h   string `json:"lowPrice24h"`
	PrevPrice24h  string `json:"prevPrice24h"`
	Volume24h     string `json:"volume24h"`
	Turnover24h   string `json:"turnover24h"`
	Price24hPcnt  string `json:"price24hPcnt"`
	UsdIndexPrice string `json:"usdIndexPrice"`
}
