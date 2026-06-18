package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"web-ui/internal/handlers"
	"web-ui/internal/services"
)

// Server represents the web server with all its dependencies
type Server struct {
	content           embed.FS
	challengeService  *services.ChallengeService
	scoreboardService *services.ScoreboardService
	userService       *services.UserService
	executionService  *services.ExecutionService
	packageService    *services.PackageService
	aiService         *services.AIService
}

// NewServer creates a new server instance
func NewServer(
	content embed.FS,
	challengeService *services.ChallengeService,
	scoreboardService *services.ScoreboardService,
	userService *services.UserService,
	executionService *services.ExecutionService,
	packageService *services.PackageService,
	aiService *services.AIService,
) *Server {
	return &Server{
		content:           content,
		challengeService:  challengeService,
		scoreboardService: scoreboardService,
		userService:       userService,
		executionService:  executionService,
		packageService:    packageService,
		aiService:         aiService,
	}
}

// SetupRoutes configures all HTTP routes
func (s *Server) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Setup static file handling
	s.setupStaticFiles(mux)

	// Initialize handlers
	apiHandler := handlers.NewAPIHandler(
		s.challengeService,
		s.scoreboardService,
		s.userService,
		s.executionService,
		s.packageService,
		s.aiService,
	)

	webHandler := handlers.NewWebHandler(
		s.content,
		s.challengeService,
		s.scoreboardService,
		s.userService,
		s.packageService,
	)

	// API routes
	mux.HandleFunc("/api/challenges", apiHandler.GetAllChallenges)
	mux.HandleFunc("/api/challenges/", apiHandler.GetChallengeByID)
	mux.HandleFunc("/api/submissions", apiHandler.HandleSubmissions)
	mux.HandleFunc("/api/scoreboard/", apiHandler.GetScoreboard)
	mux.HandleFunc("/api/run", apiHandler.RunCode)
	mux.HandleFunc("/api/save-to-filesystem", apiHandler.SaveSubmissionToFilesystem)
	mux.HandleFunc("/api/refresh-attempts", apiHandler.RefreshUserAttempts)
	mux.HandleFunc("/api/git-username", apiHandler.GetGitUsername)
	mux.HandleFunc("/api/main-scoreboard-rank", apiHandler.GetMainScoreboardRank)
	mux.HandleFunc("/api/main-leaderboard", apiHandler.GetMainLeaderboard)

	// Package challenge API routes
	mux.HandleFunc("/api/package-leaderboard", apiHandler.GetPackageLeaderboard)
	mux.HandleFunc("/api/packages/", apiHandler.HandlePackageChallenge)
	mux.HandleFunc("/api/packages-save-to-filesystem", apiHandler.SavePackageChallengeToFilesystem)

	// AI-powered API routes
	mux.HandleFunc("/api/ai/code-review", apiHandler.AICodeReview)
	mux.HandleFunc("/api/ai/interviewer-questions", apiHandler.AIInterviewerQuestions)
	mux.HandleFunc("/api/ai/code-hint", apiHandler.AICodeHint)
	mux.HandleFunc("/api/ai/debug", apiHandler.AIDebugResponse)

	// GitHub webhook route
	mux.HandleFunc("/webhook/github", apiHandler.GitHubWebhookHandler)

	// Health check endpoint for Railway
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "healthy",
			"service": "go-interview-practice",
			"version": "1.0.0",
		})
	})

	// Debug route for sponsors
	mux.HandleFunc("/api/debug/sponsors", apiHandler.GetSponsorsDebug)
	mux.HandleFunc("/api/ai/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		provider := os.Getenv("AI_PROVIDER")
		if provider == "" {
			provider = "gemini"
		}
		var apiKey string
		switch strings.ToLower(provider) {
		case "openai":
			apiKey = os.Getenv("OPENAI_API_KEY")
		case "claude":
			apiKey = os.Getenv("CLAUDE_API_KEY")
		default:
			apiKey = os.Getenv("GEMINI_API_KEY")
		}

		// Check if API key looks valid
		hasValidKey := apiKey != "" && !strings.Contains(apiKey, "Example") && len(apiKey) > 30

		response := map[string]interface{}{
			"provider":    provider,
			"status":      "ready",
			"message":     fmt.Sprintf("AI provider set to: %s", provider),
			"has_api_key": apiKey != "",
			"key_length":  len(apiKey),
			"key_preview": func() string {
				if len(apiKey) > 10 {
					return apiKey[:10] + "..."
				}
				return apiKey + "..."
			}(),
			"is_example_key": strings.Contains(apiKey, "Example"),
			"has_valid_key":  hasValidKey,
		}
		json.NewEncoder(w).Encode(response)
	})

	// Web routes
	mux.HandleFunc("/", webHandler.HomePage)
	mux.HandleFunc("/challenge/", webHandler.ChallengePage)
	mux.HandleFunc("/interview", webHandler.InterviewPage)
	mux.HandleFunc("/scoreboard", webHandler.ScoreboardPage)
	mux.HandleFunc("/scoreboard/", webHandler.ScoreChallengeHandler)
	mux.HandleFunc("/packages/", func(w http.ResponseWriter, r *http.Request) {
		// Route to appropriate handler based on URL structure
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) == 2 {
			// /packages/gin -> package detail page
			webHandler.PackageDetailPage(w, r)
		} else if len(parts) == 3 {
			if parts[2] == "scoreboard" {
				// /packages/gin/scoreboard -> package leaderboard page
				webHandler.PackageScoreboardPage(w, r)
			} else {
				// /packages/gin/challenge-1 -> package challenge page
				webHandler.PackageChallengePage(w, r)
			}
		} else {
			http.NotFound(w, r)
		}
	})

	return mux
}

// setupStaticFiles configures static file serving
func (s *Server) setupStaticFiles(mux *http.ServeMux) {
	fsys, err := fs.Sub(s.content, "static")
	if err != nil {
		log.Fatal(err)
	}

	staticHandler := http.FileServer(http.FS(fsys))
	mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set appropriate content type headers
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}
		staticHandler.ServeHTTP(w, r)
	})))
}
