package main

import (
	"url-shortener/internal"
	"url-shortener/internal/server"
	"url-shortener/internal/service"
	"url-shortener/internal/service/user"
	inmemory "url-shortener/internal/storage/inmemory/links"
	"url-shortener/internal/storage/inmemory/users"
	"url-shortener/pkg/logger"
)

func main() {
	cfg, err := internal.ReadConfig()
	if err != nil {
		panic(err)
	}

	log, _ := logger.New(cfg.Debug)

	linkStorage := inmemory.New(log)
	shortener := service.New(linkStorage, log)

	userStorage := users.New(log)
	userService := user.NewUserService(userStorage, log)

	srv := server.New(cfg.Port, shortener, userService, log)

	if err := srv.Run(); err != nil {
		panic(err)
	}
}
