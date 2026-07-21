package users

import (
	"url-shortener/internal/domain/errors"
	"url-shortener/internal/domain/models"

	"go.uber.org/zap"
)

type InMemoryUserStorage struct {
	storage map[string]models.User
	log     *zap.Logger
}

func New(log *zap.Logger) *InMemoryUserStorage {
	return &InMemoryUserStorage{
		storage: make(map[string]models.User),
		log:     log,
	}
}

func (ims *InMemoryUserStorage) SaveUser(user models.User) error {
	for uid, us := range ims.storage {
		if uid == user.ID {
			return errors.ErrUserIDAlreadyExists
		}
		if us.Email == user.Email {
			return errors.ErrUserAlreadyExists
		}
	}

	ims.storage[user.ID] = user
	return nil
}

func (ims *InMemoryUserStorage) GetUser(email string) (models.User, error) {
	for _, us := range ims.storage {
		if us.Email == email {
			return us, nil
		}
	}

	return models.User{}, errors.ErrUserNotFound
}
