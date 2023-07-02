package main

import (
	"github.com/tullur/lets-go-chat/internal/config"
	"github.com/tullur/lets-go-chat/internal/server"
)

func main() {
	conf := config.New()

	server.Run(conf)
}
