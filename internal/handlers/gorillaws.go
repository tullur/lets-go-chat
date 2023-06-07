package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tullur/lets-go-chat/internal/domain/chat"
)

var (
	clients    = make(map[*websocket.Conn]bool)
	broadcast  = make(chan chat.Message)
	register   = make(chan *websocket.Conn)
	unregister = make(chan *websocket.Conn)
	users      = make([]string, 0, len(clients))
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func authenticate(token string) bool {
	return token != ""
}

func HandleGorillaChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if !authenticate(token) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		register <- conn

		defer func() {
			unregister <- conn
			conn.Close()
		}()

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}

			log.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			broadcast <- chat.Message{Sender: conn, Content: msg}

			if err = conn.WriteMessage(msgType, msg); err != nil {
				log.Println(err)
				break
			}
		}
	}
}

func GetActiveUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(getUserList())
	}
}

func HandleMessages() {
	for {
		select {
		case client := <-register:
			clients[client] = true
			sendUserList(getUserList())
		case client := <-unregister:
			if _, ok := clients[client]; ok {
				delete(clients, client)
				client.Close()
				sendUserList(getUserList())
			}
		case message := <-broadcast:
			for client := range clients {
				if client != message.Sender {
					err := client.WriteMessage(websocket.TextMessage, message.Content)
					if err != nil {
						log.Println(err)
						break
					}
				}
			}
		}
	}
}

func getUserList() (users []string) {
	for client := range clients {
		users = append(users, client.RemoteAddr().String())
	}

	return
}

func sendUserList(usersList []string) {
	for client := range clients {
		err := client.WriteJSON(usersList)
		if err != nil {
			log.Println(err)
		}
	}
}
