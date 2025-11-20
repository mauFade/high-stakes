package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	httphandler "github.com/mauFade/high-stakes/internal/adapter/handler/http"
	"github.com/mauFade/high-stakes/internal/adapter/repository/postgres"
	"github.com/mauFade/high-stakes/internal/core/service"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDBConnectionString() string {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}

	// Get values from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "high_stakes")
	dbSSLMode := getEnv("DB_SSL_MODE", "disable")

	// Build database connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	return connStr
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}

	// Get values from environment variables
	port := getEnv("PORT", "8080")

	connStr := getDBConnectionString()

	// Initialize database connection
	db, err := postgres.NewDB(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	ur := postgres.NewUserRepository(db)

	// Initialize services
	us := service.NewUserService(ur)

	// Initialize server
	s := httphandler.NewServer(us, port)
	s.SetupRoutes()

	if err := s.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
