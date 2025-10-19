package main

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"web-ui/internal/server"
	"web-ui/internal/services"
)

//go:embed templates static
var content embed.FS

func main() {
	// Load environment variables from .env file
	loadEnvFile()

	// Initialize services
	challengeService := services.NewChallengeService()
	scoreboardService := services.NewScoreboardService()
	userService := services.NewUserService()
	executionService := services.NewExecutionService()
	packageService := services.NewPackageService()
	aiService := services.NewAIService()

	// Load data
	log.Println("Loading challenges...")
	if err := challengeService.LoadChallenges(); err != nil {
		log.Fatalf("Failed to load challenges: %v", err)
	}

	log.Println("Loading scoreboards...")
	if err := scoreboardService.LoadScoreboards(challengeService.GetChallenges()); err != nil {
		log.Fatalf("Failed to load scoreboards: %v", err)
	}

	log.Println("Loading packages...")
	if err := packageService.LoadPackages(); err != nil {
		log.Fatalf("Failed to load packages: %v", err)
	}

	// Initialize server
	srv := server.NewServer(
		content,
		challengeService,
		scoreboardService,
		userService,
		executionService,
		packageService,
		aiService,
	)

	// Setup routes
	mux := srv.SetupRoutes()

	// Start server
	port := 8080
	log.Printf("Server starting on http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

// loadEnvFile loads environment variables from a .env file
func loadEnvFile() {
	// Try to load .env from current directory and parent directories
	files := []string{".env", "../.env", "../../.env"}

	for _, file := range files {
		if loadEnvFromFile(file) {
			log.Printf("Loaded environment variables from %s", file)
			break
		}
	}
}

func loadEnvFromFile(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		// Set environment variable if not already set
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return true
}
