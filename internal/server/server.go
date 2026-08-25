package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Dorrrke/shortener-url-cource/internal/server/handlers"
	"github.com/Dorrrke/shortener-url-cource/internal/service/auth"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	srv *http.Server
	log *zap.Logger
}

func New(port int, s handlers.LinkService, us handlers.UserService, log *zap.Logger) *Server {
	addr := fmt.Sprintf(":%d", port)
	lh := handlers.NewLinkHandler(s, addr)
	uh := handlers.NewUserHandler(us)

	router := configureRoutes(lh, uh)
	log.Debug("configure routes completed")

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8080",
	}

	return &Server{
		srv: srv,
		log: log,
	}
}

func (s *Server) Run() error {
	s.log.Debug("start server", zap.String("addr", s.srv.Addr))
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.srv.Shutdown(context.TODO())
}

func configureRoutes(lh *handlers.LinkHandler, uh *handlers.UserHandler) *gin.Engine {
	router := gin.Default()

	users := router.Group("/users")
	{
		users.POST("/register", uh.Register)
		users.POST("/login", uh.Login)
		users.GET("/info", auth.AuthMiddleware, uh.GetInfo)
	}

	link := router.Group("/link")
	{
		link.POST("/add", auth.AuthMiddleware, lh.SaveURL)
		link.GET("/all", auth.AuthMiddleware, lh.GetAll)
		link.GET("/:id", lh.Get)
	}
	return router
}
