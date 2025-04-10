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

type CoinbaseMarketData struct {
	Type        string `json:"type"`
	Sequence    int64  `json:"sequence"`
	ProductID   string `json:"product_id"`
	Price       string `json:"price"`
	Open24h     string `json:"open_24h"`
	Volume24h   string `json:"volume_24h"`
	Low24h      string `json:"low_24h"`
	High24h     string `json:"high_24h"`
	Volume30d   string `json:"volume_30d"`
	BestBid     string `json:"best_bid"`
	BestBidSize string `json:"best_bid_size"`
	BestAsk     string `json:"best_ask"`
	BestAskSize string `json:"best_ask_size"`
	Side        string `json:"side"`
	Time        string `json:"time"`
	TradeID     int64  `json:"trade_id"`
	LastSize    string `json:"last_size"`
}

type OkxMarketData struct {
	InstType  string `json:"instType"`
	InstID    string `json:"instId"`
	Last      string `json:"last"`
	LastSz    string `json:"lastSz"`
	AskPx     string `json:"askPx"`
	AskSz     string `json:"askSz"`
	BidPx     string `json:"bidPx"`
	BidSz     string `json:"bidSz"`
	Open24h   string `json:"open24h"`
	High24h   string `json:"high24h"`
	Low24h    string `json:"low24h"`
	VolCcy24h string `json:"volCcy24h"`
	Vol24h    string `json:"vol24h"`
	Ts        string `json:"ts"`
	SodUtc0   string `json:"sodUtc0"`
	SodUtc8   string `json:"sodUtc8"`
}

type MoexMarketData struct {
	// TODO: add fields
}
