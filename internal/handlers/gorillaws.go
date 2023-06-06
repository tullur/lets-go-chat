package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	clients    = make(map[*websocket.Conn]bool)
	broadcast  = make(chan Message)
	register   = make(chan *websocket.Conn)
	unregister = make(chan *websocket.Conn)
	users      = make(chan []string)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Sender  *websocket.Conn
	Content []byte
}

func HandleGorillaChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

			broadcast <- Message{Sender: conn, Content: msg}

			if err = conn.WriteMessage(msgType, msg); err != nil {
				log.Println(err)
				break
			}
		}
	}
}

func HandleGetActiveUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(getUserList())
	}
}

func HandleMessages() {
	for {
		select {
		case client := <-register:
			clients[client] = true
			usersList := getUserList()
			sendUserList(usersList)
		case client := <-unregister:
			if _, ok := clients[client]; ok {
				delete(clients, client)
				client.Close()
				usersList := getUserList()
				sendUserList(usersList)

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

func getUserList() []string {
	usersList := make([]string, 0, len(clients))
	for client := range clients {
		usersList = append(usersList, client.RemoteAddr().String())
	}

	return usersList
}

func sendUserList(usersList []string) {
	for client := range clients {
		err := client.WriteJSON(usersList)
		if err != nil {
			log.Println(err)
		}
	}
}
