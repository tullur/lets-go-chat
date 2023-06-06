package handlers

import "golang.org/x/net/websocket"

type Server struct {
	conns map[*websocket.Conn]bool
}
