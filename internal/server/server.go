package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	engine *gin.Engine
	http   *http.Server
	logger zerolog.Logger
}

func NewServer(port string, logger zerolog.Logger) *Server {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: engine,
	}

	return &Server{
		engine: engine,
		http:   httpServer,
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info().Msg("Starting server")
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info().Msg("Shutting down server")
	return s.http.Shutdown(ctx)
}
