package okx

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

type OKXConnector struct {
	symbolChunks [][]string
}

type instrumentResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		InstID string `json:"instId"`
		State  string `json:"state"`
	} `json:"data"`
}

type StreamResponse struct {
	Arg struct {
		Channel string `json:"channel"`
		InstID  string `json:"instId"`
	} `json:"arg"`
	Data []json.RawMessage `json:"data"`
}

func NewConnector() *OKXConnector {
	return &OKXConnector{}
}

func (c *OKXConnector) Connect(ctx context.Context) error {
	resp, err := http.Get("https://www.okx.com/api/v5/public/instruments?instType=SPOT")
	if err != nil {
		return fmt.Errorf("get instruments: %w", err)
	}
	defer resp.Body.Close()

	var result instrumentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode instruments: %w", err)
	}

	if result.Code != "0" {
		return fmt.Errorf("API error: %s", result.Msg)
	}

	var all []string
	for _, inst := range result.Data {
		if inst.State == "live" {
			all = append(all, inst.InstID)
		}
	}

	log.Printf("OKX: found %d active instruments", len(all))
	c.symbolChunks = chunkStrings(all, 10)
	return nil
}

func (c *OKXConnector) SubscribeToMarketData(ctx context.Context, pub producer.MessageProducer) error {
	for _, chunk := range c.symbolChunks {
		conn, _, err := websocket.DefaultDialer.Dial("wss://ws.okx.com:8443/ws/v5/public", nil)
		if err != nil {
			log.Printf("dial error: %v", err)
			continue
		}

		var args []map[string]string
		for _, instID := range chunk {
			args = append(args, map[string]string{
				"channel": "tickers",
				"instId":  instID,
			})
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

			if streamMsg.Arg.Channel == "tickers" {
				log.Printf("streamMsg: %v", streamMsg)
				for _, data := range streamMsg.Data {
					if err := pub.Publish(data); err != nil {
						log.Printf("publish error: %v", err)
					}
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
