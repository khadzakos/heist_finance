package bybit

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

type BybitConnector struct {
	symbolChunks [][]string
}

type instrumentResponse struct {
	Result struct {
		List []struct {
			Symbol string `json:"symbol"`
			Status string `json:"status"`
		} `json:"list"`
	} `json:"result"`
}

type StreamResponse struct {
	Topic string          `json:"topic"`
	Ts    int64           `json:"ts"`
	Type  string          `json:"type"`
	CS    int64           `json:"cs"`
	Data  json.RawMessage `json:"data"`
}

func NewConnector() *BybitConnector {
	return &BybitConnector{}
}

func (c *BybitConnector) Connect(ctx context.Context) error {
	resp, err := http.Get("https://api.bybit.com/v5/market/instruments-info?category=spot")
	if err != nil {
		return fmt.Errorf("get instruments: %w", err)
	}
	defer resp.Body.Close()

	var result instrumentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode instruments: %w", err)
	}

	var all []string
	for _, s := range result.Result.List {
		if s.Status == "Trading" {
			all = append(all, s.Symbol)
		}
	}

	log.Printf("Bybit: found %d active symbols", len(all))
	c.symbolChunks = chunkStrings(all, 10)
	return nil
}

func (c *BybitConnector) SubscribeToMarketData(ctx context.Context, pub producer.MessageProducer) error {
	for _, chunk := range c.symbolChunks {
		conn, _, err := websocket.DefaultDialer.Dial("wss://stream.bybit.com/v5/public/spot", nil)
		if err != nil {
			log.Printf("dial error: %v", err)
			continue
		}

		var args []string
		for _, s := range chunk {
			args = append(args, "tickers."+s)
		}

		subMsg := map[string]interface{}{
			"op":   "subscribe",
			"args": args,
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
						"op": "ping",
					}); err != nil {
						log.Printf("ping error: %v", err)
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

			log.Printf("streamMsg: %v", streamMsg)

			if err := pub.Publish(streamMsg.Data); err != nil {
				log.Printf("publish error: %v", err)
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
