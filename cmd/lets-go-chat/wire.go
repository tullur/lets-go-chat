//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/tullur/lets-go-chat/internal/service"
)

func NewUser(connection string) (*service.UserService, error) {
	wire.Build(service.WithMongoRepository, service.NewUserService)

	return nil, nil
}
