package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"connector/internal/producer"

	"github.com/gorilla/websocket"
)

type BinanceConnector struct {
	tickers [][]string // разбиение на чанки
}

type ExchangeInfo struct {
	Symbols []struct {
		Symbol string `json:"symbol"`
		Status string `json:"status"`
	} `json:"symbols"`
}

type StreamResponse struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
}

func NewConnector() *BinanceConnector {
	return &BinanceConnector{}
}

func (c *BinanceConnector) Connect(ctx context.Context) error {
	resp, err := http.Get("https://api.binance.com/api/v3/exchangeInfo")
	if err != nil {
		return fmt.Errorf("failed to get exchangeInfo: %w", err)
	}
	defer resp.Body.Close()

	var info ExchangeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return fmt.Errorf("decode exchange info: %w", err)
	}

	var all []string
	for _, s := range info.Symbols {
		if s.Status == "TRADING" {
			all = append(all, strings.ToLower(s.Symbol)+"@ticker")
		}
	}

	log.Printf("Binance: found %d active trading pairs", len(all))
	c.tickers = chunkTickers(all, 200)
	return nil
}

func (c *BinanceConnector) SubscribeToMarketData(ctx context.Context, pub producer.MessageProducer) error {
	for _, chunk := range c.tickers {
		streamURL := "wss://stream.binance.com:9443/stream?streams=" + strings.Join(chunk, "/")

		conn, _, err := websocket.DefaultDialer.Dial(streamURL, nil)
		if err != nil {
			log.Printf("websocket dial failed: %v", err)
			continue
		}

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

			if err := pub.Publish(streamMsg.Data); err != nil {
				log.Printf("publish error: %v", err)
			}
		}
	}
}

func chunkTickers(tickers []string, size int) [][]string {
	var chunks [][]string
	for size < len(tickers) {
		tickers, chunks = tickers[size:], append(chunks, tickers[0:size:size])
	}
	return append(chunks, tickers)
}
