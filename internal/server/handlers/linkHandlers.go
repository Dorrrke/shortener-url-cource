package handlers

import (
	"errors"
	"fmt"
	"net/http"
	domainErrors "url-shortener/internal/domain/errors"
	"url-shortener/internal/domain/models"

	"github.com/gin-gonic/gin"
)

type LinkService interface {
	SaveURL(string, string) (string, error)
	Get(string) (string, error)
	Delete(string) error
	GetAll(string) (map[string]models.Link, error)
}

type LinkHandler struct {
	service LinkService
	domain  string
}

func NewLinkHandler(s LinkService, domain string) *LinkHandler {
	return &LinkHandler{service: s, domain: domain}
}

func (uh *LinkHandler) SaveURL(ctx *gin.Context) {
	uid := ctx.MustGet("userID").(string)

	var req struct {
		URL string `json:"url" binding:"required,url"`
	}

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lID, err := uh.service.SaveURL(req.URL, uid)
	if err != nil {
		if errors.Is(err, domainErrors.ErrInvalidLink) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainErrors.ErrLinkAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Вы уже сокращали эту ссылку"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	shortLink := fmt.Sprintf("http://%s/link/%s", uh.domain, lID)

	ctx.JSON(http.StatusOK, gin.H{"short_link": shortLink})
}

func (uh *LinkHandler) Get(ctx *gin.Context) {
	lID := ctx.Param("id")

	link, err := uh.service.Get(lID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrLinkNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, link)
}

func (uh *LinkHandler) GetAll(ctx *gin.Context) {
	uid := ctx.MustGet("userID").(string)

	urls, err := uh.service.GetAll(uid)
	if err != nil {
		if errors.Is(err, domainErrors.ErrLinkNotFound) ||
			errors.Is(err, domainErrors.ErrStorageIsEmpty) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"urls": urls})
}
