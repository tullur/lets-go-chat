package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tullur/lets-go-chat/internal/domain/chat/token/memory"
	"github.com/tullur/lets-go-chat/internal/handlers"
	"github.com/tullur/lets-go-chat/internal/service"
)

func Run(port string) {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Get("/", handlers.Greet)
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("This is panic situation")
	})

	userService, err := service.NewUserService(service.WithInMemoryRepository())
	if err != nil {
		log.Fatalln(err)
	}

	tokenService := service.NewTokenService(memory.NewMemoryTokenRepositorysitory())

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/user", UserRoutes(userService, tokenService))
		r.Mount("/chat", ChatRoutes(userService, tokenService))
	})

	http.ListenAndServe(port, r)
}
