package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tullur/lets-go-chat/internal/domain/chat"
	"github.com/tullur/lets-go-chat/internal/service"
)

const messageLimit = 10

var (
	clients    = make(map[*chat.Client]bool)
	register   = make(chan *chat.Client)
	unregister = make(chan *chat.Client)
	broadcast  = make(chan chat.Message)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GetChatUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(getUserList())
	}
}

func ChatConnection(chatService *service.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if isEmpty(token) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		currentToken, err := chatService.GetToken(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := &chat.Client{UserID: currentToken.Id(), Connection: conn}

		register <- client

		msgs, err := chatService.GetMessages(messageLimit)
		if err != nil {
			log.Println(err)
		}

		for _, msg := range msgs {
			err := client.Connection.WriteMessage(websocket.TextMessage, []byte(msg.Content))
			if err != nil {
				log.Println(err)
				break
			}

			log.Println(msg.Content)
		}

		defer func() {
			unregister <- client
			chatService.RevokeToken(token)
			client.Connection.Close()
		}()

		for {
			msgType, msg, err := client.Connection.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}

			log.Printf("%s sent: %s\n", currentToken.User(), string(msg))

			currentMessage := chat.Message{Sender: client.Connection, Content: msg}

			chatService.AddMessage(currentMessage)

			broadcast <- currentMessage

			if err = conn.WriteMessage(msgType, msg); err != nil {
				log.Println(err)
				break
			}
		}
	}
}

func BroadcastMessages() {
	for {
		select {
		case client := <-register:
			clients[client] = true
			sendUserList(getUserList())
		case client := <-unregister:
			if _, ok := clients[client]; ok {
				delete(clients, client)
				client.Connection.Close()
				sendUserList(getUserList())
			}
		case message := <-broadcast:
			for client := range clients {
				if client.Connection != message.Sender {
					err := client.Connection.WriteMessage(websocket.TextMessage, message.Content)
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
		users = append(users, client.UserID)
	}

	return
}

func sendUserList(usersList []string) {
	for client := range clients {
		err := client.Connection.WriteJSON(usersList)
		if err != nil {
			log.Println(err)
		}
	}
}

func isEmpty(token string) bool {
	return token == ""
}
