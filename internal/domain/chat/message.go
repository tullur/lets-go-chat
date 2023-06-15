package chat

import "github.com/gorilla/websocket"

type Message struct {
	Sender  *websocket.Conn
	Content []byte
}
