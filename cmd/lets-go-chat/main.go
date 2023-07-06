package main

import (
	"github.com/tullur/lets-go-chat/internal/config"
	"github.com/tullur/lets-go-chat/internal/server"
	"github.com/tullur/lets-go-chat/internal/service"
	"log"
)

type Services struct {
	user *service.UserService
	chat *service.ChatService
}

func main() {
	conf := config.New()

	us, err := NewUser(conf.Database.Host)
	if err != nil {
		log.Fatalln(err)
	}

	ch, err := NewChat(conf.Database.Host)
	if err != nil {
		log.Fatalln(err)
	}

	Services{us, ch}

	server.Run(conf)
}
