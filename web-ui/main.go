package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"web-ui/internal/server"
	"web-ui/internal/services"
)

//go:embed templates static
var content embed.FS

func main() {
	// Initialize services
	challengeService := services.NewChallengeService()
	scoreboardService := services.NewScoreboardService()
	userService := services.NewUserService()
	executionService := services.NewExecutionService()
	packageService := services.NewPackageService()

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
	)

	// Setup routes
	mux := srv.SetupRoutes()

	// Start server
	port := 8080
	log.Printf("Server starting on http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
