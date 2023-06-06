package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/tullur/lets-go-chat/internal/handlers"
)

func Run(port string) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/", handlers.Greet)

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/user", UserRoutes())
		r.Mount("/chat", ChatRoutes())
	})

	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Println(err)
	}
}
