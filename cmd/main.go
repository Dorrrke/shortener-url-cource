package main

import (
	"context"
	"url-shortener/internal"
	"url-shortener/internal/server"
	"url-shortener/internal/service"
	"url-shortener/internal/service/user"
	"url-shortener/internal/storage/db"
	inmemory "url-shortener/internal/storage/inmemory/links"
	"url-shortener/internal/storage/inmemory/users"
	"url-shortener/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	cfg, err := internal.ReadConfig()
	if err != nil {
		panic(err)
	}

	log, _ := logger.New(cfg.Debug)
	var linkStorage service.LinkRepository
	var userStorage user.UserRepository

	dbStorage, err := db.New(context.Background(), cfg.DBDSN)
	if err == nil {
		linkStorage = dbStorage
		userStorage = dbStorage
	} else {
		log.Warn("failed to connect to db:", zap.Error(err))
		linkStorage = inmemory.New(log)
		userStorage = users.New(log)
	}

	shortener := service.New(linkStorage, log)
	userService := user.NewUserService(userStorage, log)

	srv := server.New(cfg.Port, shortener, userService, log)

	if err := srv.Run(); err != nil {
		panic(err)
	}
}
