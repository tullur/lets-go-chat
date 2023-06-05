package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tullur/lets-go-chat/internal/handlers"
	"github.com/tullur/lets-go-chat/internal/service"
)

func UserRoutes() http.Handler {
	r := chi.NewRouter()

	userService, err := service.NewUserService(service.WithInMemoryRepository())
	if err != nil {
		log.Fatalln(err)
	}

	r.Get("/", handlers.HandleUserList(userService))
	r.Post("/", handlers.HandleUserCreation(userService))
	r.Post("/login", handlers.HandleUserLogin(userService))

	return r
}
