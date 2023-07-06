package main

import (
	"log"

	"github.com/tullur/lets-go-chat/internal/config"
	"github.com/tullur/lets-go-chat/internal/server"
	"github.com/tullur/lets-go-chat/internal/service"
)

type Services struct {
	user *service.UserService
	chat *service.ChatService
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/package main

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
