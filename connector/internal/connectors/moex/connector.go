package moex

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"connector/internal/producer"
)

type MOEXConnector struct {
	symbolChunks [][]string
}

type securityResponse struct {
	Securities struct {
		Data [][]interface{} `json:"data"`
	} `json:"securities"`
}

type MarketDataResponse struct {
	MarketData struct {
		Data [][]interface{} `json:"data"`
	} `json:"marketdata"`
}

type StreamResponse struct {
	ProductID string  `json:"product_id"`
	Price     float64 `json:"price"`
	Volume24h float64 `json:"volume_24h"`
	Low24h    float64 `json:"low_24h"`
	High24h   float64 `json:"high_24h"`
	BestBid   float64 `json:"best_bid"`
	BestAsk   float64 `json:"best_ask"`
	Time      string  `json:"time"`
}

func NewConnector() *MOEXConnector {
	return &MOEXConnector{}
}

func (c *MOEXConnector) Connect(ctx context.Context) error {
	resp, err := http.Get("https://iss.moex.com/iss/engines/stock/markets/shares/securities.json?iss.meta=off&iss.only=securities")
	if err != nil {
		return fmt.Errorf("get securities: %w", err)
	}
	defer resp.Body.Close()

	var result securityResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode securities: %w", err)
	}

	var all []string
	for _, sec := range result.Securities.Data {
		secID := sec[0].(string)
		status := sec[5].(string)
		if status == "A" { // "A" означает активный инструмент
			all = append(all, secID)
		}
	}

	log.Printf("MOEX: found %d active securities", len(all))
	c.symbolChunks = chunkStrings(all, 10)
	return nil
}

func (c *MOEXConnector) SubscribeToMarketData(ctx context.Context, pub producer.MessageProducer) error {
	for _, chunk := range c.symbolChunks {
		go func(symbols []string) {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					for _, symbol := range symbols {
						resp, err := http.Get(fmt.Sprintf("https://iss.moex.com/iss/engines/stock/markets/shares/securities/%s.json?iss.meta=off&iss.only=marketdata", symbol))
						if err != nil {
							log.Printf("get market data for %s: %v", symbol, err)
							continue
						}
						defer resp.Body.Close()

						var marketData MarketDataResponse
						if err := json.NewDecoder(resp.Body).Decode(&marketData); err != nil {
							log.Printf("decode market data for %s: %v", symbol, err)
							continue
						}

						if len(marketData.MarketData.Data) == 0 {
							continue
						}

						data := marketData.MarketData.Data[0]
						streamMsg := StreamResponse{
							ProductID: symbol,
							Price:     safeFloat64(data[12]), // LAST - последняя цена
							Volume24h: safeFloat64(data[20]), // VOLTODAY - объем за день
							Low24h:    safeFloat64(data[15]), // LOW - минимум за день
							High24h:   safeFloat64(data[14]), // HIGH - максимум за день
							BestBid:   safeFloat64(data[9]),  // BID - лучшая цена покупки
							BestAsk:   safeFloat64(data[10]), // OFFER - лучшая цена продажи
							Time:      time.Now().Format(time.RFC3339),
						}

						log.Printf("streamMsg for %s: %v", symbol, streamMsg)

						msg, err := json.Marshal(streamMsg)
						if err != nil {
							log.Printf("marshal error for %s: %v", symbol, err)
							continue
						}

						if err := pub.Publish(msg); err != nil {
							log.Printf("publish error for %s: %v", symbol, err)
						}
					}
				}
			}
		}(chunk)
	}

	<-ctx.Done()
	return ctx.Err()
}

func safeFloat64(v interface{}) float64 {
	if v == nil {
		return 0
	}
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

func chunkStrings(list []string, size int) [][]string {
	var chunks [][]string
	for size < len(list) {
		list, chunks = list[size:], append(chunks, list[0:size:size])
	}
	return append(chunks, list)
}
