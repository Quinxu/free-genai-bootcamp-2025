package server

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"lang-portal/internal/api/handlers"
	"lang-portal/internal/service"
)

// Config holds server configuration
type Config struct {
	Port int
}

// Server represents the HTTP server
type Server struct {
	router  *gin.Engine
	config  Config
	service *service.Services
}

// NewServer creates a new server instance
func NewServer(config Config, services *service.Services) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Use logger and recovery middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	return &Server{
		router:  router,
		config:  config,
		service: services,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Register routes
	s.registerRoutes()

	// Start server
	addr := fmt.Sprintf(":%d", s.config.Port)
	return s.router.Run(addr)
}

// registerRoutes registers all API routes
func (s *Server) registerRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	api := s.router.Group("/api")
	{
		// Register handlers
		wordHandler := handlers.NewWordHandler(s.service.Word)
		groupHandler := handlers.NewGroupHandler(s.service.Group)
		studyHandler := handlers.NewStudyHandler(s.service.Study)

		// Register routes
		wordHandler.RegisterRoutes(api)
		groupHandler.RegisterRoutes(api)
		studyHandler.RegisterRoutes(api)
	}
}
