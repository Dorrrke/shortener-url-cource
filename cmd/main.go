package main

import (
	"context"

	"github.com/Dorrrke/shortener-url-cource/internal"
	"github.com/Dorrrke/shortener-url-cource/internal/server"
	"github.com/Dorrrke/shortener-url-cource/internal/service"
	"github.com/Dorrrke/shortener-url-cource/internal/service/user"
	"github.com/Dorrrke/shortener-url-cource/internal/storage/db"
	inmemory "github.com/Dorrrke/shortener-url-cource/internal/storage/inmemory/links"
	"github.com/Dorrrke/shortener-url-cource/internal/storage/inmemory/users"
	"github.com/Dorrrke/shortener-url-cource/pkg/logger"

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
