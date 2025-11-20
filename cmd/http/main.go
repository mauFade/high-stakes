package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	httphandler "github.com/mauFade/high-stakes/internal/adapter/handler/http"
	"github.com/mauFade/high-stakes/internal/adapter/repository/postgres"
	"github.com/mauFade/high-stakes/internal/core/service"
)

func main() {
	// Parse command line flags
	port := flag.String("port", "8080", "Server port")
	dbHost := flag.String("db-host", "localhost", "Database host")
	dbPort := flag.String("db-port", "5432", "Database port")
	dbUser := flag.String("db-user", "postgres", "Database user")
	dbPassword := flag.String("db-password", "", "Database password")
	dbName := flag.String("db-name", "high_stakes", "Database name")
	dbSSLMode := flag.String("db-sslmode", "disable", "Database SSL mode")
	flag.Parse()

	// Build database connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		*dbHost, *dbPort, *dbUser, *dbPassword, *dbName, *dbSSLMode,
	)

	// Initialize database connection
	db, err := postgres.NewDB(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	userRepo := postgres.NewUserRepository(db)

	// Initialize service
	userService := service.NewUserService(userRepo)

	// Initialize router
	router := httphandler.NewRouter(userService)
	mux := router.SetupRoutes()

	// Start server
	addr := ":" + *port
	log.Printf("Server starting on port %s", *port)
	log.Printf("Health check: http://localhost%s/health", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
