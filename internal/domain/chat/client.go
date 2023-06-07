package chat

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID     string
	Connection *websocket.Conn
}
