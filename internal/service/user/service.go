package user

import (
	"errors"
	domainErrors "url-shortener/internal/domain/errors"
	"url-shortener/internal/domain/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	SaveUser(models.User) error
	GetUser(string) (models.User, error)
	GetUserInfo(string) (models.User, []models.Link, error)
}

type UserService struct {
	UserRepository UserRepository
	log            *zap.Logger
}

func NewUserService(ur UserRepository, log *zap.Logger) *UserService {
	return &UserService{
		UserRepository: ur,
		log:            log,
	}
}

func (u *UserService) SaveUser(regReq models.RegisterRequest) (string, error) {
	userID := uuid.New().String()

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(regReq.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	user := models.User{
		ID:       userID,
		Name:     regReq.Name,
		Email:    regReq.Email,
		Password: string(hash),
	}

	err = u.UserRepository.SaveUser(user)
	if err == nil {
		return userID, nil
	}

	if errors.Is(err, domainErrors.ErrUserIDAlreadyExists) {
		userID = uuid.New().String()
		user.ID = userID
		err = u.UserRepository.SaveUser(user)
		if err == nil {
			return userID, err
		}
	}

	return "", err
}

func (u *UserService) ValidationUser(loginReq models.LoginRequest) (string, error) {
	user, err := u.UserRepository.GetUser(loginReq.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(loginReq.Password),
	)
	if err != nil {
		return "", domainErrors.ErrInvalidCredentials
	}

	return user.ID, nil
}

func (u *UserService) GetUserInfo(userID string) (models.UserInfo, error) {
	user, links, err := u.UserRepository.GetUserInfo(userID)
	if err != nil {
		return models.UserInfo{}, err
	}

	u.log.Debug("User info retrieved", zap.String("userID", userID), zap.Int("linksCount", len(links)))
	u.log.Debug("Links:", zap.Any("links", links))

	return models.UserInfo{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Links:     links,
	}, nil
}
