package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tullur/lets-go-chat/internal/config"
	"github.com/tullur/lets-go-chat/internal/handlers"
	"github.com/tullur/lets-go-chat/internal/service"

	_ "github.com/joho/godotenv/autoload"
)

func Run(conf *config.Config) {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Get("/", handlers.Greet)
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("This is panic situation")
	})

	userService, err := service.NewUserService(service.WithMongoRepository(conf.Database.Host))
	if err != nil {
		log.Fatalln(err)
	}

	tokenService, err := service.NewTokenService(service.WithInMemoryTokenRepository())
	if err != nil {
		log.Fatalln(err)
	}

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/user", UserRoutes(userService, tokenService))
		r.Mount("/chat", ChatRoutes(tokenService))
	})

	server := &http.Server{
		Handler:      r,
		Addr:         conf.Server.Addr,
		ReadTimeout:  conf.Server.ReadTimeout,
		WriteTimeout: conf.Server.WriteTimeout,
		IdleTimeout:  conf.Server.IdleTimeout,
	}

	log.Fatalln(server.ListenAndServe())
}
