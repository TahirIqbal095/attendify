package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/tahiriqbal095/attendify/internal/db"
)

type Server struct {
	engine *gin.Engine
	http   *http.Server
	logger zerolog.Logger
	pool   *db.Pool
}

func NewServer(port string, logger zerolog.Logger, pool *db.Pool) *Server {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.GET("/health", func(c *gin.Context) {
		// Check database connectivity
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"error":   "database unavailable",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"status": "ok",
			},
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
		pool:   pool,
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
