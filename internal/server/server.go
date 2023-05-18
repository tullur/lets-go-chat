package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Run(port string) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/user", UserRoutes())
	})

	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Println(err)
	}
}
