package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"lang-portal/internal/api/handlers"
	"lang-portal/internal/service"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Initialize database
	db, err := sql.Open("sqlite3", "./data/lang_portal.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Initialize services with database connection
	studyService := service.NewStudyService(db)
	groupService := service.NewGroupService(db)
	wordService := service.NewWordService(db)

	// Initialize handlers
	studyHandler := handlers.NewStudyHandler(studyService)
	groupHandler := handlers.NewGroupHandler(groupService)
	wordHandler := handlers.NewWordHandler(wordService)

	// Setup route groups
	api := r.Group("/api")
	{
		// Register routes using handler methods
		studyHandler.RegisterRoutes(api)
		groupHandler.RegisterRoutes(api)
		wordHandler.RegisterRoutes(api)
	}

	return r
}

func main() {
	r := setupRouter()

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
