package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	Conn *websocket.Conn
	URL  string
}

func NewWSClient(url string) *WSClient {
	return &WSClient{URL: url}
}

func (c *WSClient) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.URL, nil)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

func (c *WSClient) ReadMessage() ([]byte, error) {
	_, msg, err := c.Conn.ReadMessage()
	return msg, err
}

func (c *WSClient) WriteMessage(messageType int, data []byte) error {
	return c.Conn.WriteMessage(messageType, data)
}

func (c *WSClient) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

func (c *WSClient) Close() error {
	return c.Conn.Close()
}
