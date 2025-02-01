package wsclient

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Stream string `json:"stream"`
	Data   struct {
		Price string `json:"p"`
		Time  int64  `json:"T"`
		Bid   string `json:"b"`
		Ask   string `json:"a"`
	} `json:"data"`
}

type WebSocketClient struct {
	conn    *websocket.Conn
	url     string
	tickers []string
}

func NewWebSocketClient(wsURL string, tickers []string) *WebSocketClient {
	return &WebSocketClient{
		url:     wsURL,
		tickers: tickers,
	}
}

func (c *WebSocketClient) Connect() error {
	u, err := url.Parse(c.url)
	if err != nil {
		return err
	}

	streams := make([]string, len(c.tickers))
	for i, ticker := range c.tickers {
		streams[i] = ticker + "@trade"
	}

	q := u.Query()
	q.Set("streams", strings.Join(streams, "/"))
	u.RawQuery = q.Encode()

	c.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer c.conn.Close()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Ошибка чтения WebSocket:", err)
			time.Sleep(5 * time.Second)
			c.Reconnect()
			continue
		}
		fmt.Printf("Received: %s\n", message)
	}
}

func (c *WebSocketClient) Listen() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Ошибка чтения WebSocket:", err)
			time.Sleep(5 * time.Second)
			c.Reconnect()
			continue
		}

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("Binance: Ошибка разбора JSON:", err)
			continue
		}
	}
}

func (c *WebSocketClient) Reconnect() {
	log.Println("Binance: Переподключение...")
	c.Connect()
}
