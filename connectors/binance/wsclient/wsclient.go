package wsclient

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type TransactionData struct {
	Symbol   string `json:"s"`
	Price    string `json:"p"`
	Quantity string `json:"q"`
	Time     int64  `json:"T"`
}

type Transaction struct {
	Stream string          `json:"stream"`
	Data   TransactionData `json:"data"`
}

type WebSocketClient struct {
	conn    *websocket.Conn
	url     string
	tickers []string
	trades  chan<- TransactionData
}

func NewWebSocketClient(wsURL string, tickers []string, trades chan<- TransactionData) *WebSocketClient {
	return &WebSocketClient{
		url:     wsURL,
		tickers: tickers,
		trades:  trades,
	}
}

func (c *WebSocketClient) ConnectWS() {
	for {
		u, err := url.Parse(c.url)
		if err != nil {
			log.Printf("Error parsing URL: %v", err)
			continue
		}

		streams := make([]string, len(c.tickers))
		for i, ticker := range c.tickers {
			streams[i] = ticker + "@aggTrade"
		}

		q := u.Query()
		q.Set("streams", strings.Join(streams, "/"))
		u.RawQuery = q.Encode()

		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Printf("WebSocket connection error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		c.conn = conn
		c.listen()
	}
}

func (c *WebSocketClient) listen() {
	defer c.conn.Close()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			return
		}

		var msg Transaction
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("JSON parse error: %v", err)
			continue
		}
		log.Println("Message: ", msg)
		c.trades <- msg.Data
	}
}
