package main

import (
	"log"
	"task-manager/internal/api"
	"task-manager/internal/api/handler"
	"task-manager/internal/config"
	"task-manager/internal/repository/sqlite"
	"task-manager/internal/service"
)

func main() {
	// initialize the config values
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// connect to sqlite database
	db, err := sqlite.NewDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	// initialize repositories
	taskRepo := sqlite.NewTaskRepository(db)

	// initialize services
	taskService := service.NewTaskService(taskRepo)

	// initialize handlers
	taskHandler := handler.NewTaskHandler(taskService)

	// setup router
	r := api.NewRouter(taskHandler)

	log.Printf("Starting server on port %s", cfg.Port)
	log.Printf("Database path: %s", cfg.DatabasePath)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
