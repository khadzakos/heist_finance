package wsclient

import (
	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	conn    *websocket.Conn
	url     string
	tickers []string
}
