package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"web-ui/internal/models"
	"web-ui/internal/services"
	"web-ui/internal/utils"
)

// GitHubSponsorsResponse represents the GitHub GraphQL response for sponsors
type GitHubSponsorsResponse struct {
	Data struct {
		Viewer struct {
			SponsorshipsAsMaintainer struct {
				Nodes []struct {
					SponsorEntity struct {
						Login string `json:"login"`
					} `json:"sponsorEntity"`
				} `json:"nodes"`
			} `json:"sponsorshipsAsMaintainer"`
		} `json:"viewer"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// SponsorCache caches sponsor data to avoid hitting API limits
type SponsorCache struct {
	sponsors    map[string]bool
	lastUpdated time.Time
	mutex       sync.RWMutex
}

var sponsorCache = &SponsorCache{
	sponsors: make(map[string]bool),
}

// LoadSponsors loads sponsor list by scraping the public GitHub sponsors page
func (h *APIHandler) LoadSponsors() map[string]bool {
	sponsorCache.mutex.RLock()
	// Return cached data if it's less than 1 hour old
	if time.Since(sponsorCache.lastUpdated) < time.Hour && len(sponsorCache.sponsors) > 0 {
		defer sponsorCache.mutex.RUnlock()
		return sponsorCache.sponsors
	}
	sponsorCache.mutex.RUnlock()

	// Scrape sponsors from the public GitHub sponsors page
	sponsors := h.scrapeSponsorsFromGitHub()

	// Update cache
	sponsorCache.mutex.Lock()
	sponsorCache.sponsors = sponsors
	sponsorCache.lastUpdated = time.Now()
	sponsorCache.mutex.Unlock()
	return sponsors
}

// scrapeSponsorsFromGitHub scrapes the public GitHub sponsors page
func (h *APIHandler) scrapeSponsorsFromGitHub() map[string]bool {
	sponsorMap := make(map[string]bool)

	// Create HTTP client with timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Fetch the public sponsors page
	url := "https://github.com/sponsors/RezaSi"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return sponsorMap
	}

	// Set user agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; GoSponsorScraper/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching sponsors page: %v\n", err)
		return sponsorMap
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("GitHub sponsors page returned status %d\n", resp.StatusCode)
		return sponsorMap
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return sponsorMap
	}

	html := string(body)

	// Extract usernames from the HTML using regex - look for avatar images with alt="@username"
	avatarRegex := regexp.MustCompile(`alt="@([a-zA-Z0-9][a-zA-Z0-9\-]*)"`)
	matches := avatarRegex.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			username := match[1]
			// Filter out the repository owner from sponsors list
			if username != "RezaSi" {
				sponsorMap[username] = true
			}
		}
	}

	// Fallback: if no sponsors found with avatar method, try href patterns
	if len(sponsorMap) == 0 {
		// Look for href="/username" patterns that aren't common GitHub paths
		linkRegex := regexp.MustCompile(`href="/([a-zA-Z0-9][a-zA-Z0-9\-]+)"`)
		linkMatches := linkRegex.FindAllStringSubmatch(html, -1)

		for _, match := range linkMatches {
			if len(match) >= 2 {
				username := match[1]
				// Filter out common GitHub paths that aren't usernames
				if username != "sponsors" && username != "github" && username != "RezaSi" &&
					!strings.HasPrefix(username, "orgs/") &&
					!strings.Contains(username, "/") &&
					len(username) > 2 { // reasonable username length
					sponsorMap[username] = true
				}
			}
		}
	}

	return sponsorMap
}

// APIHandler handles all API endpoints
type APIHandler struct {
	challengeService  *services.ChallengeService
	scoreboardService *services.ScoreboardService
	userService       *services.UserService
	executionService  *services.ExecutionService
	packageService    *services.PackageService
	aiService         *services.AIService
	submissions       []models.Submission
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(
	challengeService *services.ChallengeService,
	scoreboardService *services.ScoreboardService,
	userService *services.UserService,
	executionService *services.ExecutionService,
	packageService *services.PackageService,
	aiService *services.AIService,
) *APIHandler {
	return &APIHandler{
		challengeService:  challengeService,
		scoreboardService: scoreboardService,
		userService:       userService,
		executionService:  executionService,
		packageService:    packageService,
		aiService:         aiService,
		submissions:       make([]models.Submission, 0),
	}
}

// GetAllChallenges returns all challenges
func (h *APIHandler) GetAllChallenges(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	challenges := h.challengeService.GetChallenges()

	// Convert map to slice for JSON response
	var challengeList []*models.Challenge
	for _, challenge := range challenges {
		challengeList = append(challengeList, challenge)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challengeList)
}

// GetChallengeByID returns a specific challenge by ID
func (h *APIHandler) GetChallengeByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract challenge ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/challenges/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	challenge, exists := h.challengeService.GetChallenge(id)
	if !exists {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenge)
}

// HandleSubmissions handles submission operations
func (h *APIHandler) HandleSubmissions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.createSubmission(w, r)
	case "GET":
		h.getSubmissions(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// createSubmission creates a new submission
func (h *APIHandler) createSubmission(w http.ResponseWriter, r *http.Request) {
	var submission models.Submission
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Set submission timestamp
	submission.SubmittedAt = time.Now()

	// Validate challenge exists
	challenge, exists := h.challengeService.GetChallenge(submission.ChallengeID)
	if !exists {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	// Run the code
	result := h.executionService.RunCode(submission.Code, challenge)
	submission.Passed = result.Passed
	submission.TestOutput = result.Output
	submission.ExecutionMs = result.ExecutionMs

	// Store submission
	h.submissions = append(h.submissions, submission)

	// Add to scoreboard if passed
	if submission.Passed {
		h.scoreboardService.AddSubmission(submission)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(submission)
}

// getSubmissions returns all submissions
func (h *APIHandler) getSubmissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.submissions)
}

// GetScoreboard returns the scoreboard for a challenge
func (h *APIHandler) GetScoreboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract challenge ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/scoreboard/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	scoreboard, exists := h.scoreboardService.GetScoreboard(id)
	if !exists {
		scoreboard = []models.ScoreboardEntry{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scoreboard)
}

// RunCode executes submitted code
func (h *APIHandler) RunCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ChallengeID int    `json:"challengeId"`
		Code        string `json:"code"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	challenge, exists := h.challengeService.GetChallenge(request.ChallengeID)
	if !exists {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	result := h.executionService.RunCode(request.Code, challenge)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// SaveSubmissionToFilesystem saves a submission to the filesystem
func (h *APIHandler) SaveSubmissionToFilesystem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request services.SaveSubmissionRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	if request.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Set username cookie
	h.setUsernameCookie(w, request.Username)

	// Validate challenge exists
	_, exists := h.challengeService.GetChallenge(request.ChallengeID)
	if !exists {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	response := h.executionService.SaveSubmissionToFilesystem(request)

	// Clear user attempts cache
	h.userService.RefreshUserAttempts(request.Username, h.challengeService.GetChallenges())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RefreshUserAttempts refreshes user's attempt cache
func (h *APIHandler) RefreshUserAttempts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Username string `json:"username"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	if request.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	attempts := h.userService.RefreshUserAttempts(request.Username, h.challengeService.GetChallenges())

	response := struct {
		Username     string       `json:"username"`
		AttemptedIDs map[int]bool `json:"attemptedIds"`
		Scores       map[int]int  `json:"scores"`
		Success      bool         `json:"success"`
	}{
		Username:     request.Username,
		AttemptedIDs: attempts.AttemptedIDs,
		Scores:       attempts.Scores,
		Success:      true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetGitUsername returns the username extracted from git configuration
func (h *APIHandler) GetGitUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gitInfo := utils.GetGitUsername()

	response := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Source   string `json:"source"`
		Success  bool   `json:"success"`
	}{
		Username: gitInfo.Username,
		Email:    gitInfo.Email,
		Source:   gitInfo.Source,
		Success:  gitInfo.Username != "",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// setUsernameCookie sets the username cookie
func (h *APIHandler) setUsernameCookie(w http.ResponseWriter, username string) {
	// Cookie expires in 30 days
	expiration := time.Now().Add(30 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     "username",
		Value:    username,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: false, // Allow JavaScript to access it
	}
	http.SetCookie(w, &cookie)
}

// GetMainScoreboardRank returns the user's rank in the main scoreboard
func (h *APIHandler) GetMainScoreboardRank(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username parameter required", http.StatusBadRequest)
		return
	}

	// Calculate user's rank in main scoreboard
	rank := h.calculateMainScoreboardRank(username)

	response := struct {
		Username string `json:"username"`
		Rank     int    `json:"rank"`
		Success  bool   `json:"success"`
	}{
		Username: username,
		Rank:     rank,
		Success:  true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// calculateMainScoreboardRank calculates the user's rank based on completed challenges
func (h *APIHandler) calculateMainScoreboardRank(username string) int {
	// Get all users and their completion counts (only count if ALL tests passed)
	challenges := h.challengeService.GetChallenges()
	userCompletions := make(map[string]int)

	// Process all challenge scoreboards to count actual completions
	for challengeID := range challenges {
		// Read scoreboard file directly to check test results
		scoreboardPath := fmt.Sprintf("../challenge-%d/SCOREBOARD.md", challengeID)
		content, err := ioutil.ReadFile(scoreboardPath)
		if err != nil {
			// Try alternative path
			scoreboardPath = fmt.Sprintf("challenge-%d/SCOREBOARD.md", challengeID)
			content, err = ioutil.ReadFile(scoreboardPath)
			if err != nil {
				continue
			}
		}

		// Parse scoreboard to find users who passed ALL tests
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			// Skip header and separator lines
			if !strings.Contains(line, "|") || strings.Contains(line, "Username") || strings.Contains(line, "---") {
				continue
			}

			parts := strings.Split(line, "|")
			if len(parts) < 4 {
				continue
			}

			username := strings.TrimSpace(parts[1])
			passedTestsStr := strings.TrimSpace(parts[2])
			totalTestsStr := strings.TrimSpace(parts[3])

			// Skip empty usernames or placeholders
			if username == "" || username == "------" {
				continue
			}

			// Parse test counts
			passedTests, err1 := strconv.Atoi(passedTestsStr)
			totalTests, err2 := strconv.Atoi(totalTestsStr)

			// Only count as completed if ALL tests passed
			if err1 == nil && err2 == nil && passedTests > 0 && passedTests == totalTests {
				userCompletions[username]++
			}
		}
	}

	// Get the target user's completion count
	targetCompletions := userCompletions[username]
	if targetCompletions == 0 {
		return 0 // User is unranked
	}

	// Count how many users have more completions (following Python script logic)
	rank := 1
	for user, completions := range userCompletions {
		if user != username && completions > targetCompletions {
			rank++
		}
	}

	return rank
}

// GetMainLeaderboard returns the main leaderboard data
func (h *APIHandler) GetMainLeaderboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Calculate leaderboard data
	leaderboard := h.calculateMainLeaderboard()

	// Include total number of classic challenges for dynamic UI rendering
	totalChallenges := len(h.challengeService.GetChallenges())

	response := struct {
		Leaderboard     []LeaderboardUser `json:"leaderboard"`
		Success         bool              `json:"success"`
		TotalChallenges int               `json:"totalChallenges"`
	}{
		Leaderboard:     leaderboard,
		Success:         true,
		TotalChallenges: totalChallenges,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPackageLeaderboard returns leaderboard data for a package learning path
func (h *APIHandler) GetPackageLeaderboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	packageName := r.URL.Query().Get("package")
	if packageName == "" {
		http.Error(w, "package parameter required", http.StatusBadRequest)
		return
	}

	// Load package and its challenges
	pkg, err := h.packageService.GetPackage(packageName)
	if err != nil {
		http.Error(w, "Package not found", http.StatusNotFound)
		return
	}

	// Build package challenge list in learning path order
	challengesMap, err := h.packageService.GetPackageChallenges(packageName)
	if err != nil {
		challengesMap = make(map[string]*models.PackageChallenge)
	}
	var challenges []*models.PackageChallenge
	for _, id := range pkg.LearningPath {
		if ch, ok := challengesMap[id]; ok {
			challenges = append(challenges, ch)
		}
	}

	// Reuse existing creator to gather leaderboard
	leaderboard := h.createPackageLeaderboard(packageName, challenges)

	response := map[string]interface{}{
		"success":         true,
		"leaderboard":     leaderboard,
		"totalChallenges": len(challenges),
		"package":         pkg.Name,
		"displayName":     pkg.DisplayName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// createPackageLeaderboard builds the package leaderboard based on filesystem submissions
func (h *APIHandler) createPackageLeaderboard(packageName string, challenges []*models.PackageChallenge) []models.PackageScoreboardEntry {
	var leaderboard []models.PackageScoreboardEntry
	type userPackageStats struct {
		username            string
		completedCount      int
		lastSubmission      time.Time
		challengesCompleted map[string]bool
	}
	userStats := make(map[string]*userPackageStats)

	// Load sponsors for package leaderboard
	sponsors := h.LoadSponsors()

	for _, challenge := range challenges {
		submissionsDir := filepath.Join("..", "packages", packageName, challenge.ID, "submissions")
		if _, err := os.Stat(submissionsDir); os.IsNotExist(err) {
			continue
		}
		entries, err := ioutil.ReadDir(submissionsDir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				username := entry.Name()
				userDir := filepath.Join(submissionsDir, username)
				solutionPath := filepath.Join(userDir, "solution.go")
				altSolutionPath := filepath.Join(userDir, "solution-template.go")

				var modTime time.Time
				if stat, err := os.Stat(solutionPath); err == nil {
					modTime = stat.ModTime()
				} else if stat, err := os.Stat(altSolutionPath); err == nil {
					modTime = stat.ModTime()
				} else {
					continue
				}

				if userStats[username] == nil {
					userStats[username] = &userPackageStats{
						username:            username,
						completedCount:      0,
						lastSubmission:      modTime,
						challengesCompleted: make(map[string]bool),
					}
				}
				if !userStats[username].challengesCompleted[challenge.ID] {
					userStats[username].completedCount++
					userStats[username].challengesCompleted[challenge.ID] = true
					if modTime.After(userStats[username].lastSubmission) {
						userStats[username].lastSubmission = modTime
					}
				}
			}
		}
	}

	for username, stats := range userStats {
		if stats.completedCount > 0 {
			leaderboard = append(leaderboard, models.PackageScoreboardEntry{
				Username:    username,
				PackageName: packageName,
				ChallengeID: "",
				SubmittedAt: stats.lastSubmission,
				TestsPassed: stats.completedCount,
				TestsTotal:  len(challenges),
				IsSponsor:   sponsors[username],
			})
		}
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		if leaderboard[i].TestsPassed != leaderboard[j].TestsPassed {
			return leaderboard[i].TestsPassed > leaderboard[j].TestsPassed
		}
		return leaderboard[i].SubmittedAt.Before(leaderboard[j].SubmittedAt)
	})

	return leaderboard
}

// LeaderboardUser represents a user in the leaderboard
type LeaderboardUser struct {
	Username            string       `json:"username"`
	CompletedCount      int          `json:"completedCount"`
	CompletionRate      float64      `json:"completionRate"`
	CompletedChallenges map[int]bool `json:"completedChallenges"`
	Achievement         string       `json:"achievement"`
	Rank                int          `json:"rank"`
	IsSponsor           bool         `json:"isSponsor"`
}

// calculateMainLeaderboard calculates the main leaderboard data
func (h *APIHandler) calculateMainLeaderboard() []LeaderboardUser {
	challenges := h.challengeService.GetChallenges()
	totalChallenges := len(challenges)
	userCompletions := make(map[string]map[int]bool)

	// Load sponsor information
	sponsors := h.LoadSponsors()

	// Process all challenge scoreboards to find completions
	for challengeID := range challenges {
		// Read scoreboard file directly to check test results
		scoreboardPath := fmt.Sprintf("../challenge-%d/SCOREBOARD.md", challengeID)
		content, err := ioutil.ReadFile(scoreboardPath)
		if err != nil {
			// Try alternative path
			scoreboardPath = fmt.Sprintf("challenge-%d/SCOREBOARD.md", challengeID)
			content, err = ioutil.ReadFile(scoreboardPath)
			if err != nil {
				continue
			}
		}

		// Parse scoreboard to find users who passed ALL tests
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			// Skip header and separator lines
			if !strings.Contains(line, "|") || strings.Contains(line, "Username") || strings.Contains(line, "---") {
				continue
			}

			parts := strings.Split(line, "|")
			if len(parts) < 4 {
				continue
			}

			username := strings.TrimSpace(parts[1])
			passedTestsStr := strings.TrimSpace(parts[2])
			totalTestsStr := strings.TrimSpace(parts[3])

			// Skip empty usernames or placeholders
			if username == "" || username == "------" {
				continue
			}

			// Parse test counts
			passedTests, err1 := strconv.Atoi(passedTestsStr)
			totalTests, err2 := strconv.Atoi(totalTestsStr)

			// Only count as completed if ALL tests passed
			if err1 == nil && err2 == nil && passedTests > 0 && passedTests == totalTests {
				if userCompletions[username] == nil {
					userCompletions[username] = make(map[int]bool)
				}
				userCompletions[username][challengeID] = true
			}
		}
	}

	// Convert to leaderboard format
	var leaderboard []LeaderboardUser
	for username, completions := range userCompletions {
		completedCount := len(completions)
		completionRate := float64(completedCount) / float64(totalChallenges) * 100

		// Determine achievement
		achievement := "🌱 Beginner"
		if completedCount >= 20 {
			achievement = "🔥 Master"
		} else if completedCount >= 15 {
			achievement = "⭐ Expert"
		} else if completedCount >= 10 {
			achievement = "💪 Advanced"
		} else if completedCount >= 5 {
			achievement = "🚀 Intermediate"
		}

		leaderboard = append(leaderboard, LeaderboardUser{
			Username:            username,
			CompletedCount:      completedCount,
			CompletionRate:      completionRate,
			CompletedChallenges: completions,
			Achievement:         achievement,
			IsSponsor:           sponsors[username],
		})
	}

	// Sort by completion count (descending), then by username
	for i := 0; i < len(leaderboard); i++ {
		for j := i + 1; j < len(leaderboard); j++ {
			if leaderboard[j].CompletedCount > leaderboard[i].CompletedCount ||
				(leaderboard[j].CompletedCount == leaderboard[i].CompletedCount && leaderboard[j].Username < leaderboard[i].Username) {
				leaderboard[i], leaderboard[j] = leaderboard[j], leaderboard[i]
			}
		}
	}

	// Assign ranks
	for i := range leaderboard {
		leaderboard[i].Rank = i + 1
	}

	return leaderboard
}

// HandlePackageChallenge handles package challenge test and submit requests
func (h *APIHandler) HandlePackageChallenge(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse URL path: /api/packages/{packageName}/{challengeId}/{action}
	path := strings.TrimPrefix(r.URL.Path, "/api/packages/")
	parts := strings.Split(path, "/")

	if len(parts) != 3 {
		http.Error(w, "Invalid URL format. Expected: /api/packages/{packageName}/{challengeId}/{action}", http.StatusBadRequest)
		return
	}

	packageName := parts[0]
	challengeId := parts[1]
	action := parts[2] // "test" or "submit"

	// Validate action
	if action != "test" && action != "submit" {
		http.Error(w, "Invalid action. Must be 'test' or 'submit'", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request struct {
		Code     string `json:"code"`
		Username string `json:"username"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if request.Code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	// Use the existing package service
	packageService := h.packageService

	challenge, err := packageService.GetPackageChallenge(packageName, challengeId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Challenge not found: %v", err), http.StatusNotFound)
		return
	}

	// Convert PackageChallenge to Challenge format for ExecutionService
	challengeForExecution := &models.Challenge{
		ID:       0, // Package challenges don't use numeric IDs
		Title:    challenge.Title,
		TestFile: challenge.TestFile,
	}

	// Run the actual tests using ExecutionService
	result := h.executionService.RunCode(request.Code, challengeForExecution)

	// Format response
	response := map[string]interface{}{
		"success":      result.Passed,
		"execution_ms": result.ExecutionMs,
		"output":       result.Output,
	}

	// Count passed tests from output for display
	testsPassed, testsTotal := h.parseTestResults(result.Output)
	response["tests_passed"] = testsPassed
	response["tests_total"] = testsTotal

	if action == "submit" && result.Passed {
		response["message"] = "Solution submitted successfully!"
		response["show_pr_instructions"] = true

		// Set username cookie if provided
		if request.Username != "" {
			h.setUsernameCookie(w, request.Username)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// parseTestResults parses Go test output to count passed and total tests
func (h *APIHandler) parseTestResults(output string) (passed int, total int) {
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		// Look for test result lines like "--- PASS: TestGetUsers" or "--- FAIL: TestCreateUser"
		if strings.Contains(line, "--- PASS:") {
			passed++
			total++
		} else if strings.Contains(line, "--- FAIL:") {
			total++
		}
	}

	// If no individual test results found, check for overall result
	if total == 0 {
		if strings.Contains(output, "PASS") && !strings.Contains(output, "FAIL") {
			// Assume basic success case
			passed = 1
			total = 1
		} else if strings.Contains(output, "FAIL") {
			// Assume basic failure case
			passed = 0
			total = 1
		}
	}

	return passed, total
}

// SavePackageChallengeToFilesystem saves a package challenge submission to the filesystem
func (h *APIHandler) SavePackageChallengeToFilesystem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Username    string `json:"username"`
		PackageName string `json:"packageName"`
		ChallengeID string `json:"challengeId"`
		Code        string `json:"code"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	if request.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Set username cookie
	h.setUsernameCookie(w, request.Username)

	// Validate challenge exists
	_, err = h.packageService.GetPackageChallenge(request.PackageName, request.ChallengeID)
	if err != nil {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	// Save to filesystem
	response := h.savePackageChallengeToFilesystem(request)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// savePackageChallengeToFilesystem handles the actual file saving for package challenges
func (h *APIHandler) savePackageChallengeToFilesystem(request struct {
	Username    string `json:"username"`
	PackageName string `json:"packageName"`
	ChallengeID string `json:"challengeId"`
	Code        string `json:"code"`
}) services.SaveSubmissionResponse {
	// Get working directory for correct relative paths
	workDir, _ := os.Getwd()

	// Try different path approaches to handle potential path issues
	var submissionDir string
	var fileSaved bool

	// Try multiple path options for package challenges
	pathOptions := []string{
		// Option 1: From web-ui directory (standard case)
		filepath.Join("..", "packages", request.PackageName, request.ChallengeID, "submissions", request.Username),
		// Option 2: From root workspace
		filepath.Join("packages", request.PackageName, request.ChallengeID, "submissions", request.Username),
		// Option 3: Absolute path from detected workspace root
		filepath.Join(workDir, "..", "packages", request.PackageName, request.ChallengeID, "submissions", request.Username),
	}

	for _, dirPath := range pathOptions {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			continue
		}

		solutionFile := filepath.Join(dirPath, "solution.go")
		err = ioutil.WriteFile(solutionFile, []byte(request.Code), 0644)
		if err != nil {
			continue
		}

		submissionDir = dirPath
		fileSaved = true
		break
	}

	if !fileSaved {
		return services.SaveSubmissionResponse{
			Success: false,
			Message: "Failed to save solution to any available path",
		}
	}

	// Return success response with git commands
	relativePath := filepath.Join("packages", request.PackageName, request.ChallengeID, "submissions", request.Username, "solution.go")
	return services.SaveSubmissionResponse{
		Success:  true,
		Message:  "Solution saved to filesystem",
		FilePath: filepath.Join(submissionDir, "solution.go"),
		GitCommands: []string{
			"cd " + filepath.Join(workDir, ".."),
			fmt.Sprintf("git add %s", relativePath),
			fmt.Sprintf("git commit -m \"Add solution for %s %s by %s\"", request.PackageName, request.ChallengeID, request.Username),
			"git push origin main",
		},
	}
}

// AICodeReview performs AI-powered code review
func (h *APIHandler) AICodeReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ChallengeID int    `json:"challengeId"`
		Code        string `json:"code"`
		Context     string `json:"context"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	challenge, exists := h.challengeService.GetChallenge(request.ChallengeID)
	if !exists {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	review, err := h.aiService.ReviewCode(request.Code, challenge, request.Context)
	if err != nil {
		http.Error(w, fmt.Sprintf("AI review failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

// AIInterviewerQuestions generates AI interviewer questions
func (h *APIHandler) AIInterviewerQuestions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ChallengeID  int    `json:"challengeId"`
		Code         string `json:"code"`
		UserProgress string `json:"userProgress"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	challenge, exists := h.challengeService.GetChallenge(request.ChallengeID)
	if !exists {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	questions, err := h.aiService.GetInterviewerQuestions(request.Code, challenge, request.UserProgress)
	if err != nil {
		http.Error(w, fmt.Sprintf("AI questions failed: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		Questions []string `json:"questions"`
		Success   bool     `json:"success"`
	}{
		Questions: questions,
		Success:   true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AICodeHint provides AI-powered code hints
func (h *APIHandler) AICodeHint(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ChallengeID int    `json:"challengeId"`
		Code        string `json:"code"`
		HintLevel   int    `json:"hintLevel"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	challenge, exists := h.challengeService.GetChallenge(request.ChallengeID)
	if !exists {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	// Validate hint level
	if request.HintLevel < 1 || request.HintLevel > 4 {
		request.HintLevel = 1
	}

	hint, err := h.aiService.GetCodeHint(request.Code, challenge, request.HintLevel)
	if err != nil {
		http.Error(w, fmt.Sprintf("AI hint failed: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		Hint      string `json:"hint"`
		HintLevel int    `json:"hintLevel"`
		Success   bool   `json:"success"`
	}{
		Hint:      hint,
		HintLevel: request.HintLevel,
		Success:   true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AIDebugResponse provides raw AI response for debugging
func (h *APIHandler) AIDebugResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ChallengeID int    `json:"challengeId"`
		Code        string `json:"code"`
		Context     string `json:"context"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	challenge, exists := h.challengeService.GetChallenge(request.ChallengeID)
	if !exists {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	// Get raw AI response for debugging
	prompt := h.aiService.BuildCodeReviewPrompt(request.Code, challenge, request.Context)
	rawResponse, err := h.aiService.CallLLMRaw(prompt)

	response := struct {
		RawResponse string `json:"raw_response"`
		Prompt      string `json:"prompt"`
		Success     bool   `json:"success"`
		Error       string `json:"error,omitempty"`
	}{
		RawResponse: rawResponse,
		Prompt:      prompt,
		Success:     err == nil,
	}

	if err != nil {
		response.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GitHubWebhookHandler handles GitHub sponsor webhooks
func (h *APIHandler) GitHubWebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Parse the webhook payload
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Check if this is a sponsorship event
	eventType := r.Header.Get("X-GitHub-Event")
	if eventType == "sponsorship" {
		// Clear the sponsor cache to force a refresh on next request
		sponsorCache.mutex.Lock()
		sponsorCache.sponsors = make(map[string]bool)
		sponsorCache.lastUpdated = time.Time{} // Reset to zero time to force refresh
		sponsorCache.mutex.Unlock()

		fmt.Printf("Sponsor cache cleared due to webhook event: %s\n", eventType)
	}

	// Respond with 200 OK to acknowledge receipt
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GetSponsorsDebug returns current sponsors for debugging
func (h *APIHandler) GetSponsorsDebug(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sponsors := h.LoadSponsors()

	response := struct {
		Sponsors map[string]bool `json:"sponsors"`
		Count    int             `json:"count"`
		Success  bool            `json:"success"`
	}{
		Sponsors: sponsors,
		Count:    len(sponsors),
		Success:  true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
