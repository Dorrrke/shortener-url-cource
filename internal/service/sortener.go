package service

import (
	"errors"
	"math/rand"
	"strings"
	domainErrors "url-shortener/internal/domain/errors"
	"url-shortener/internal/domain/models"

	"go.uber.org/zap"
)

type LinkRepository interface {
	Save(models.Link) error
	Get(string) (string, error)
	Delete(string) error
	GetAll(string) (map[string]models.Link, error)
}

type Shortener struct {
	storage LinkRepository
	log     *zap.Logger
}

func New(stor LinkRepository, log *zap.Logger) *Shortener {
	return &Shortener{
		storage: stor,
		log:     log,
	}
}

func (s *Shortener) SaveURL(url string, uid string) (string, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "", domainErrors.ErrInvalidLink
	}

	lID := generateCode()

	link := models.Link{
		Short:    lID,
		Original: url,
		UserID:   uid,
	}

	err := s.storage.Save(link)
	if err == nil {
		return lID, nil
	}

	if errors.Is(err, domainErrors.ErrCodeAlreadyExists) {
		lID = generateCode()
		link.Short = lID
		err := s.storage.Save(link)
		if err == nil {
			return lID, nil
		}
	}

	return "", err
}

func (s *Shortener) Get(code string) (string, error) {
	link, err := s.storage.Get(code)
	if err != nil {
		return "", err
	}

	return link, nil
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
