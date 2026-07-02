package main

import (
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	handler "day16/internal/handler/http"
	repo "day16/internal/repository/postgres"
	"day16/internal/service"
	"day16/pkg/logger"
)

func main() {
	appLogger := logger.New()

	dsn := "host=127.0.0.1 user=postgres password=secret dbname=mydb port=5433 sslmode=disable"


	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		appLogger.Error("Failed to connect to database: %v", err)
		return
	}

	if err := db.AutoMigrate(&repo.UserDBModel{}); err != nil {
		appLogger.Error("Failed to run auto migration: %v", err)
		return
	}
	appLogger.Info("Database migration completed successfully.")

	userRepo := repo.NewPostgresUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, appLogger)

	http.HandleFunc("/users/", userHandler.HandleUserCRUD)

	appLogger.Info("Server running on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		appLogger.Error("Server crashed: %v", err)
	}
}
