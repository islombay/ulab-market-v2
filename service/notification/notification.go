package notification

import "github.com/gorilla/websocket"

type SocketConnection struct {
	Connection *websocket.Conn
	Type       string
}
