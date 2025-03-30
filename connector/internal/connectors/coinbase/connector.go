package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"connector/internal/producer"

	"github.com/gorilla/websocket"
)

type CoinbaseConnector struct {
	symbolChunks [][]string
}

type productResponse []struct {
	ProductID string `json:"id"`
	Status    string `json:"status"`
}

type StreamResponse struct {
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

func NewConnector() *CoinbaseConnector {
	return &CoinbaseConnector{}
}

func (c *CoinbaseConnector) Connect(ctx context.Context) error {
	resp, err := http.Get("https://api.exchange.coinbase.com/products")
	if err != nil {
		return fmt.Errorf("get products: %w", err)
	}
	defer resp.Body.Close()

	var result productResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode products: %w", err)
	}

	var all []string
	for _, p := range result {
		if p.Status == "online" {
			all = append(all, p.ProductID)
		}
	}

	log.Printf("Coinbase: found %d active products", len(all))
	c.symbolChunks = chunkStrings(all, 10)
	return nil
}

func (c *CoinbaseConnector) SubscribeToMarketData(ctx context.Context, pub producer.MessageProducer) error {
	for _, chunk := range c.symbolChunks {
		conn, _, err := websocket.DefaultDialer.Dial("wss://ws-feed.exchange.coinbase.com", nil)
		if err != nil {
			log.Printf("dial error: %v", err)
			continue
		}

		subMsg := map[string]interface{}{
			"type":        "subscribe",
			"channels":    []string{"ticker"},
			"product_ids": chunk,
		}

		if err := conn.WriteJSON(subMsg); err != nil {
			log.Printf("subscribe error: %v", err)
			continue
		}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read subscription response error: %v", err)
			continue
		}
		log.Printf("subscription response: %s", string(msg))

		go func() {
			ticker := time.NewTicker(15 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if err := conn.WriteJSON(map[string]interface{}{
						"type": "heartbeat",
						"on":   true,
					}); err != nil {
						log.Printf("heartbeat error: %v", err)
						return
					}
				}
			}
		}()

		go handleConnection(ctx, conn, pub)
	}

	<-ctx.Done()
	return ctx.Err()
}

func handleConnection(ctx context.Context, conn *websocket.Conn, pub producer.MessageProducer) {
	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read error: %v", err)
				return
			}

			var streamMsg StreamResponse
			if err := json.Unmarshal(msg, &streamMsg); err != nil {
				log.Printf("unmarshal error: %v", err)
				continue
			}

			if streamMsg.Type == "ticker" {
				log.Printf("streamMsg: %v", streamMsg)

				msg, err := json.Marshal(streamMsg)
				if err != nil {
					log.Printf("marshal error: %v", err)
					continue
				}

				if err := pub.Publish(msg); err != nil {
					log.Printf("publish error: %v", err)
				}
			}
		}
	}
}

func chunkStrings(list []string, size int) [][]string {
	var chunks [][]string
	for size < len(list) {
		list, chunks = list[size:], append(chunks, list[0:size:size])
	}
	return append(chunks, list)
}
