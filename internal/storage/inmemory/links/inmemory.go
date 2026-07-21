package inmemory

import (
	"url-shortener/internal/domain/errors"
	"url-shortener/internal/domain/models"

	"go.uber.org/zap"
)

type InMemoryLinkStorage struct {
	storage map[string]models.Link
	log     *zap.Logger
}

func New(log *zap.Logger) *InMemoryLinkStorage {
	return &InMemoryLinkStorage{
		storage: make(map[string]models.Link),
		log:     log,
	}
}

func (ims *InMemoryLinkStorage) Save(link models.Link) error {
	for code, dblink := range ims.storage {
		if code == link.Short {
			return errors.ErrCodeAlreadyExists
		}
		if dblink.Original == link.Original && dblink.UserID == link.UserID {
			return errors.ErrLinkAlreadyExists
		}
	}

	ims.storage[link.Short] = link
	return nil
}

func (ims *InMemoryLinkStorage) Get(code string) (models.Link, error) {
	link, ok := ims.storage[code]
	if !ok {
		return models.Link{}, errors.ErrLinkNotFound
	}

	return link, nil
}

func (ims *InMemoryLinkStorage) Delete(code string) error {
	_, ok := ims.storage[code]
	if !ok {
		return errors.ErrLinkNotFound
	}
	delete(ims.storage, code)
	return nil
}

func (ims *InMemoryLinkStorage) GetAll(uid string) (map[string]models.Link, error) {
	if len(ims.storage) == 0 {
		return nil, errors.ErrStorageIsEmpty
	}

	userMap := make(map[string]models.Link)

	for code, link := range ims.storage {
		if link.UserID == uid {
			userMap[code] = link
		}
	}

	if len(userMap) == 0 {
		return nil, errors.ErrLinkNotFound
	}

	return userMap, nil
}
