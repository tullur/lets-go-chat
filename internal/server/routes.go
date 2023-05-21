package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tullur/lets-go-chat/internal/domain/user/memory"
	"github.com/tullur/lets-go-chat/internal/handlers"
	"github.com/tullur/lets-go-chat/internal/service"
)

func UserRoutes() http.Handler {
	r := chi.NewRouter()

	userService := service.NewUser(memory.NewInMemoryRepository())

	r.Get("/", handlers.HandleUserList(userService))
	r.Post("/", handlers.HandleUserCreation(userService))
	r.Post("/login", handlers.HandleUserLogin(userService))

	return r
}
