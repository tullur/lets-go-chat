package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tullur/lets-go-chat/internal/handlers"
	"github.com/tullur/lets-go-chat/internal/service"
)

func UserRoutes(userService *service.UserService, chatService *service.ChatService) http.Handler {
	r := chi.NewRouter()

	r.Get("/", handlers.GetUsers(userService))
	r.Post("/", handlers.CreateUser(userService))
	r.Post("/login", handlers.LoginUser(userService, chatService))

	return r
}

func ChatRoutes(userService *service.UserService, chatService *service.ChatService) http.Handler {
	r := chi.NewRouter()

	r.Handle("/", handlers.ChatConnection(chatService))
	r.Get("/users", handlers.GetChatUsers())

	go handlers.BroadcastMessages()

	return r
}
