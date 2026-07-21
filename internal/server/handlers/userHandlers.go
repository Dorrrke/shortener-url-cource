package handlers

import (
	"net/http"
	"time"
	"url-shortener/internal/domain/models"
	"url-shortener/internal/service/auth"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	SaveUser(models.RegisterRequest) (string, error)
	ValidationUser(models.LoginRequest) (string, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(s UserService) *UserHandler { return &UserHandler{service: s} }

func (uh *UserHandler) Register(ctx *gin.Context) {
	var regReq models.RegisterRequest
	err := ctx.ShouldBindBodyWithJSON(&regReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uh.service.SaveUser(regReq)
	if err != nil {
		//TODO: добавить валидацию ошибки
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user_id": userID})
}

func (uh *UserHandler) Login(ctx *gin.Context) {
	var loginReq models.LoginRequest
	err := ctx.ShouldBindBodyWithJSON(&loginReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uh.service.ValidationUser(loginReq)
	if err != nil {
		//TODO: добавить валидацию ошибки
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.GenerateToken(userID, 15*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Authorization", "Bearer "+token)
	ctx.JSON(http.StatusOK, gin.H{"user_id": userID})
}
