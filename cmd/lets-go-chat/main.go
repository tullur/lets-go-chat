package main

import (
	"github.com/tullur/lets-go-chat/internal/config"
	"github.com/tullur/lets-go-chat/internal/server"
)

func main() {
	c := config.New()

	server.Run(c)
}
