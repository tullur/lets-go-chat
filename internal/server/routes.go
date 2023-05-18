package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tullur/lets-go-chat/internal/domain/user"
	"github.com/tullur/lets-go-chat/internal/handlers"
)

func UserRoutes() http.Handler {
	r := chi.NewRouter()
	userStore := user.New()

	r.Get("/", handlers.HandleUserList(userStore))
	r.Post("/", handlers.HandleUserCreation(userStore))
	r.Post("/login", handlers.HandleUserLogin(userStore))

	return r
}
