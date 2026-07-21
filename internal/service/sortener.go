package service

import (
	"math/rand"
	"strings"
	"url-shortener/internal/domain/errors"
	"url-shortener/internal/domain/models"

	"go.uber.org/zap"
)

type Storage interface {
	Save(models.Link) error
	Get(string) (models.Link, error)
	Delete(string) error
	GetAll(string) (map[string]models.Link, error)
}

type Shortener struct {
	storage Storage
	log     *zap.Logger
}

func New(stor Storage, log *zap.Logger) *Shortener {
	return &Shortener{
		storage: stor,
		log:     log,
	}
}

func (s *Shortener) SaveURL(url string, uid string) (string, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", errors.ErrInvalidLink
	}

	lID := generateCode()

	link := models.Link{
		Short:    lID,
		Original: url,
		UserID:   uid,
	}

	err := s.storage.Save(link)
	if err != nil {
		return "", err
	}

	return lID, nil
}

func (s *Shortener) Get(code string) (models.Link, error) {
	url, err := s.storage.Get(code)
	if err != nil {
		return models.Link{}, err
	}

	return url, nil
}

func (s *Shortener) Delete(code string) error {
	return s.storage.Delete(code)
}

func (s *Shortener) GetAll(uid string) (map[string]models.Link, error) {
	return s.storage.GetAll(uid)
}

func generateCode() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

	var result strings.Builder

	for i := 0; i < 6; i++ {
		result.WriteByte(
			chars[rand.Intn(len(chars))],
		)
	}

	return result.String()
}
