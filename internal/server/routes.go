package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tullur/lets-go-chat/internal/domain/user/memory"
	"github.com/tullur/lets-go-chat/internal/handlers"
)

func UserRoutes() http.Handler {
	r := chi.NewRouter()
	userInMemoryStore := memory.NewInMemoryRepository()

	r.Get("/", handlers.HandleUserList(userInMemoryStore))
	r.Post("/", handlers.HandleUserCreation(userInMemoryStore))
	r.Post("/login", handlers.HandleUserLogin(userInMemoryStore))

	return r
}
