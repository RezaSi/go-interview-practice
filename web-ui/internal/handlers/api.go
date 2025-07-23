package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"web-ui/internal/models"
	"web-ui/internal/services"
	"web-ui/internal/utils"
)

// APIHandler handles all API endpoints
type APIHandler struct {
	challengeService  *services.ChallengeService
	scoreboardService *services.ScoreboardService
	userService       *services.UserService
	executionService  *services.ExecutionService
	packageService    *services.PackageService
	submissions       []models.Submission
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(
	challengeService *services.ChallengeService,
	scoreboardService *services.ScoreboardService,
	userService *services.UserService,
	executionService *services.ExecutionService,
	packageService *services.PackageService,
) *APIHandler {
	return &APIHandler{
		challengeService:  challengeService,
		scoreboardService: scoreboardService,
		userService:       userService,
		executionService:  executionService,
		packageService:    packageService,
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

	response := struct {
		Leaderboard []LeaderboardUser `json:"leaderboard"`
		Success     bool              `json:"success"`
	}{
		Leaderboard: leaderboard,
		Success:     true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// LeaderboardUser represents a user in the leaderboard
type LeaderboardUser struct {
	Username            string       `json:"username"`
	CompletedCount      int          `json:"completedCount"`
	CompletionRate      float64      `json:"completionRate"`
	CompletedChallenges map[int]bool `json:"completedChallenges"`
	Achievement         string       `json:"achievement"`
	Rank                int          `json:"rank"`
}

// calculateMainLeaderboard calculates the main leaderboard data
func (h *APIHandler) calculateMainLeaderboard() []LeaderboardUser {
	challenges := h.challengeService.GetChallenges()
	totalChallenges := len(challenges)
	userCompletions := make(map[string]map[int]bool)

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
