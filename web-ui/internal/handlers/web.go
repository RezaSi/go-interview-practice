package handlers

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"io/ioutil"
	"os"
	"path/filepath"
	"web-ui/internal/models"
	"web-ui/internal/services"
	"web-ui/internal/utils"
)

// WebHandler handles web page rendering
type WebHandler struct {
	content           embed.FS
	challengeService  *services.ChallengeService
	scoreboardService *services.ScoreboardService
	userService       *services.UserService
	packageService    *services.PackageService
}

// NewWebHandler creates a new web handler
func NewWebHandler(
	content embed.FS,
	challengeService *services.ChallengeService,
	scoreboardService *services.ScoreboardService,
	userService *services.UserService,
	packageService *services.PackageService,
) *WebHandler {
	return &WebHandler{
		content:           content,
		challengeService:  challengeService,
		scoreboardService: scoreboardService,
		userService:       userService,
		packageService:    packageService,
	}
}

// HomePage renders the home page with a list of challenges
func (h *WebHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.New("").Funcs(utils.GetTemplateFuncs()).ParseFS(h.content, "templates/base.html", "templates/home.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert map to slice for template
	var challengeList []*models.Challenge
	for _, challenge := range h.challengeService.GetChallenges() {
		challengeList = append(challengeList, challenge)
	}

	// Get packages for the Package Mastery tab
	packages := h.packageService.GetPackages()

	// Convert packages map to sorted slice by stars (descending)
	type PackageWithName struct {
		Name string
		*models.Package
	}

	var packagesList []*PackageWithName
	for name, pkg := range packages {
		packagesList = append(packagesList, &PackageWithName{
			Name:    name,
			Package: pkg,
		})
	}

	// Sort by stars descending (highest first)
	sort.Slice(packagesList, func(i, j int) bool {
		return packagesList[i].Stars > packagesList[j].Stars
	})

	// Get the username from cookie if available
	username := h.getUsernameFromCookie(r)

	// Get user attempts if username is set
	var userAttempt *models.UserAttemptedChallenges
	if username != "" {
		userAttempt = h.userService.GetUserAttempts(username, h.challengeService.GetChallenges())

		// For gin package, also check package challenge attempts
		if userAttempt == nil {
			userAttempt = &models.UserAttemptedChallenges{
				AttemptedIDs: make(map[int]bool),
				Scores:       make(map[int]int),
			}
		}

		// Add package challenge attempts to userAttempt for UI consistency
		for packageName, pkg := range packages {
			for i, challengeID := range pkg.LearningPath {
				attempted := h.hasUserAttemptedPackageChallenge(username, packageName, challengeID)
				if attempted {
					// Use negative IDs for package challenges to avoid conflicts with classic challenges
					// Create unique negative ID based on package and challenge index
					packageChallengeID := -(1000 + i*10 + len(packageName)) // Ensure unique negative IDs
					userAttempt.AttemptedIDs[packageChallengeID] = true
				}
			}
		}
	}

	data := struct {
		Challenges   []*models.Challenge
		Username     string
		UserAttempts *models.UserAttemptedChallenges
		Packages     map[string]*models.Package
		PackagesList []*PackageWithName
	}{
		Challenges:   challengeList,
		Username:     username,
		UserAttempts: userAttempt,
		Packages:     packages,
		PackagesList: packagesList,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		// Don't call http.Error here since headers may already be sent during template execution
	}
}

// ChallengePage renders a specific challenge page
func (h *WebHandler) ChallengePage(w http.ResponseWriter, r *http.Request) {
	// Extract challenge ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/challenge/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	challenge, exists := h.challengeService.GetChallenge(id)
	if !exists {
		http.NotFound(w, r)
		return
	}

	// Get username from cookie first
	username := h.getUsernameFromCookie(r)

	// If no username from cookie, try to get it from Git config
	if username == "" {
		gitInfo := utils.GetGitUsername()
		if gitInfo.Username != "" {
			username = gitInfo.Username
			// Set the cookie for future requests
			h.setUsernameCookie(w, username)
		}
	}

	existingSolution := ""
	hasAttempted := false

	if username != "" {
		existingSolution = h.userService.GetExistingSolution(username, id)
		// Check if user has attempted this challenge
		userAttempts := h.userService.GetUserAttempts(username, h.challengeService.GetChallenges())
		hasAttempted = userAttempts.AttemptedIDs[id]
	}

	tmpl, err := template.New("").Funcs(utils.GetTemplateFuncs()).ParseFS(h.content, "templates/base.html", "templates/challenge.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Challenge        *models.Challenge
		Username         string
		ExistingSolution string
		HasAttempted     bool
	}{
		Challenge:        challenge,
		Username:         username,
		ExistingSolution: existingSolution,
		HasAttempted:     hasAttempted,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		// Don't call http.Error here since headers may already be sent during template execution
	}
}

// ScoreboardPage renders the main scoreboard page
func (h *WebHandler) ScoreboardPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").Funcs(utils.GetTemplateFuncs()).ParseFS(h.content, "templates/base.html", "templates/scoreboard.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get all challenges for the scoreboard overview
	challenges := h.challengeService.GetChallenges()
	scoreboards := h.scoreboardService.GetAllScoreboards()

	data := struct {
		Challenges  models.ChallengeMap
		Scoreboards models.ScoreboardMap
	}{
		Challenges:  challenges,
		Scoreboards: scoreboards,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		// Don't call http.Error here since headers may already be sent during template execution
	}
}

// ScoreChallengeHandler renders scoreboard for a specific challenge
func (h *WebHandler) ScoreChallengeHandler(w http.ResponseWriter, r *http.Request) {
	// Extract challenge ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/scoreboard/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	challenge, exists := h.challengeService.GetChallenge(id)
	if !exists {
		http.NotFound(w, r)
		return
	}

	scoreboard, _ := h.scoreboardService.GetScoreboard(id)

	tmpl, err := template.New("").Funcs(utils.GetTemplateFuncs()).ParseFS(h.content, "templates/base.html", "templates/challenge_scoreboard.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Challenge *models.Challenge
		Entries   []models.ScoreboardEntry
	}{
		Challenge: challenge,
		Entries:   scoreboard,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		// Don't call http.Error here since headers may already be sent during template execution
	}
}

// PackageScoreboardPage renders the package-wide scoreboard page with same theme as main scoreboard
func (h *WebHandler) PackageScoreboardPage(w http.ResponseWriter, r *http.Request) {
	// URL format: /packages/{package}/scoreboard
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "packages" || parts[2] != "scoreboard" {
		http.NotFound(w, r)
		return
	}

	packageName := parts[1]

	// Get package and challenges
	pkg, err := h.packageService.GetPackage(packageName)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	challengesMap, err := h.packageService.GetPackageChallenges(packageName)
	if err != nil {
		challengesMap = make(map[string]*models.PackageChallenge)
	}

	// Build list of challenges in learning path order
	var challenges []*models.PackageChallenge
	for _, id := range pkg.LearningPath {
		if ch, ok := challengesMap[id]; ok {
			challenges = append(challenges, ch)
		}
	}

	// Create leaderboard
	leaderboard := h.createPackageLeaderboard(packageName, challenges)

	tmpl, err := template.New("").Funcs(utils.GetTemplateFuncs()).ParseFS(h.content, "templates/base.html", "templates/package_scoreboard.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Package         *models.Package
		Leaderboard     []models.PackageScoreboardEntry
		TotalChallenges int
	}{
		Package:         pkg,
		Leaderboard:     leaderboard,
		TotalChallenges: len(challenges),
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

// InterviewPage renders the interview simulator setup and runner
func (h *WebHandler) InterviewPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").Funcs(utils.GetTemplateFuncs()).ParseFS(h.content, "templates/base.html", "templates/interview.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert map to slice for template
	var challengeList []*models.Challenge
	for _, challenge := range h.challengeService.GetChallenges() {
		challengeList = append(challengeList, challenge)
	}

	// Get username from cookie if available
	username := h.getUsernameFromCookie(r)

	data := struct {
		Challenges []*models.Challenge
		Username   string
	}{
		Challenges: challengeList,
		Username:   username,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		// Don't call http.Error here since headers may already be sent during template execution
	}
}

// getUsernameFromCookie retrieves the username from cookie
func (h *WebHandler) getUsernameFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("username")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// setUsernameCookie sets the username cookie
func (h *WebHandler) setUsernameCookie(w http.ResponseWriter, username string) {
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

// PackageDetailPage renders the package detail page
func (h *WebHandler) PackageDetailPage(w http.ResponseWriter, r *http.Request) {
	// Extract package name from URL: /packages/gin
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 || parts[0] != "packages" {
		http.NotFound(w, r)
		return
	}

	packageName := parts[1]

	// Get package data
	pkg, err := h.packageService.GetPackage(packageName)
	if err != nil {
		log.Printf("Package not found: %v", err)
		http.NotFound(w, r)
		return
	}

	// Get challenges for this package
	challengesMap, err := h.packageService.GetPackageChallenges(packageName)
	if err != nil {
		log.Printf("Error getting package challenges: %v", err)
		challengesMap = make(map[string]*models.PackageChallenge)
	}

	// Convert map to sorted slice and add missing fields
	var challenges []*models.PackageChallenge
	challengeIDs := make([]string, 0, len(challengesMap))

	// Collect and sort challenge IDs
	for id := range challengesMap {
		challengeIDs = append(challengeIDs, id)
	}
	sort.Strings(challengeIDs)

	// Convert map to sorted slice using learning path order
	for _, challengeID := range pkg.LearningPath {
		if challenge, exists := challengesMap[challengeID]; exists {
			challenges = append(challenges, challenge)
		}
	}

	tmpl, err := template.New("").Funcs(utils.GetTemplateFuncs()).ParseFS(h.content, "templates/base.html", "templates/package_detail.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the username from cookie if available
	username := h.getUsernameFromCookie(r)

	// Check which package challenges the user has attempted
	packageAttempts := make(map[string]bool)
	completedCount := 0
	if username != "" {
		for _, challenge := range challenges {
			attempted := h.hasUserAttemptedPackageChallenge(username, packageName, challenge.ID)
			packageAttempts[challenge.ID] = attempted
			if attempted {
				completedCount++
			}
		}
	}

	// Calculate user progress
	userProgress := struct {
		CompletedCount     int
		ProgressPercentage float64
	}{
		CompletedCount:     completedCount,
		ProgressPercentage: float64(completedCount) / float64(len(challenges)) * 100,
	}

	// Create submission counts map for each challenge
	submissionCounts := make(map[string]int)
	for _, challenge := range challenges {
		submissionCounts[challenge.ID] = h.countPackageChallengeSubmissions(packageName, challenge.ID)
	}

	// Create actual leaderboard using submission data
	leaderboard := h.createPackageLeaderboard(packageName, challenges)

	data := struct {
		Package          *models.Package
		Challenges       []*models.PackageChallenge
		Username         string
		UserProgress     interface{}
		TotalChallenges  int
		Leaderboard      []models.PackageScoreboardEntry
		PackageAttempts  map[string]bool
		SubmissionCounts map[string]int
	}{
		Package:          pkg,
		Challenges:       challenges,
		Username:         username,
		UserProgress:     userProgress,
		TotalChallenges:  len(challenges),
		Leaderboard:      leaderboard,
		PackageAttempts:  packageAttempts,
		SubmissionCounts: submissionCounts,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
	}
}

// PackageChallengePage renders an individual package challenge page
func (h *WebHandler) PackageChallengePage(w http.ResponseWriter, r *http.Request) {
	// Extract package and challenge from URL: /packages/gin/challenge-1-basic-routing
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "packages" {
		http.NotFound(w, r)
		return
	}

	packageName := parts[1]
	challengeID := parts[2]

	// Get package data
	pkg, err := h.packageService.GetPackage(packageName)
	if err != nil {
		log.Printf("Package not found: %v", err)
		http.Error(w, "Package not found", http.StatusNotFound)
		return
	}

	// Get challenge data
	challenge, err := h.packageService.GetPackageChallenge(packageName, challengeID)
	if err != nil {
		log.Printf("Challenge not found: %v", err)
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	tmpl, err := template.New("").Funcs(utils.GetTemplateFuncs()).ParseFS(h.content, "templates/base.html", "templates/package_challenge.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get username from cookie first
	username := h.getUsernameFromCookie(r)

	// If no username from cookie, try to get it from Git config
	if username == "" {
		gitInfo := utils.GetGitUsername()
		if gitInfo.Username != "" {
			username = gitInfo.Username
			// Set the cookie for future requests
			h.setUsernameCookie(w, username)
		}
	}

	// Check if user has attempted this challenge
	hasAttempted := false
	existingSolution := ""
	if username != "" {
		hasAttempted = h.hasUserAttemptedPackageChallenge(username, packageName, challengeID)
		existingSolution = h.getUserPackageChallengeSolution(username, packageName, challengeID)
	}

	data := struct {
		Package          *models.Package
		Challenge        *models.PackageChallenge
		Username         string
		SubmissionCount  int
		HasAttempted     bool
		ExistingSolution string
	}{
		Package:          pkg,
		Challenge:        challenge,
		Username:         username,
		SubmissionCount:  0,
		HasAttempted:     hasAttempted,
		ExistingSolution: existingSolution,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		// Don't call http.Error here since headers may already be sent during template execution
	}
}

// hasUserAttemptedPackageChallenge checks if a user has attempted a package challenge
func (h *WebHandler) hasUserAttemptedPackageChallenge(username, packageName, challengeID string) bool {
	// Check if submission file exists in ../packages/{packageName}/{challengeID}/submissions/{username}/solution.go
	submissionPath := filepath.Join("..", "packages", packageName, challengeID, "submissions", username, "solution.go")
	if _, err := os.Stat(submissionPath); err == nil {
		return true
	}

	// Try alternative path in case of different file naming
	altSubmissionPath := filepath.Join("..", "packages", packageName, challengeID, "submissions", username, "solution-template.go")
	if _, err := os.Stat(altSubmissionPath); err == nil {
		return true
	}

	return false
}

// getUserPackageChallengeSolution retrieves a user's existing solution for a package challenge
func (h *WebHandler) getUserPackageChallengeSolution(username, packageName, challengeID string) string {
	if username == "" {
		return ""
	}

	// Try solution.go first
	submissionPath := filepath.Join("..", "packages", packageName, challengeID, "submissions", username, "solution.go")
	content, err := ioutil.ReadFile(submissionPath)
	if err == nil {
		return string(content)
	}

	// Try solution-template.go as fallback
	altSubmissionPath := filepath.Join("..", "packages", packageName, challengeID, "submissions", username, "solution-template.go")
	content, err = ioutil.ReadFile(altSubmissionPath)
	if err == nil {
		return string(content)
	}

	return ""
}

// countPackageChallengeSubmissions counts the number of submissions for a package challenge
func (h *WebHandler) countPackageChallengeSubmissions(packageName, challengeID string) int {
	submissionsDir := filepath.Join("..", "packages", packageName, challengeID, "submissions")

	// Check if submissions directory exists
	if _, err := os.Stat(submissionsDir); os.IsNotExist(err) {
		return 0
	}

	// Read the submissions directory
	entries, err := ioutil.ReadDir(submissionsDir)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			// Check if this user directory has a solution file
			userDir := filepath.Join(submissionsDir, entry.Name())
			solutionPath := filepath.Join(userDir, "solution.go")
			altSolutionPath := filepath.Join(userDir, "solution-template.go")

			if _, err := os.Stat(solutionPath); err == nil {
				count++
			} else if _, err := os.Stat(altSolutionPath); err == nil {
				count++
			}
		}
	}

	return count
}

// createPackageLeaderboard creates a leaderboard for package challenges similar to classic challenges
func (h *WebHandler) createPackageLeaderboard(packageName string, challenges []*models.PackageChallenge) []models.PackageScoreboardEntry {
	var leaderboard []models.PackageScoreboardEntry
	userStats := make(map[string]*userPackageStats)
	
	// Load sponsors for package leaderboard (reuse from API handler)
	// Create a temporary API handler instance to access LoadSponsors
	tempHandler := &APIHandler{}
	sponsors := tempHandler.LoadSponsors()

	// Collect submission data for each challenge
	for _, challenge := range challenges {
		submissionsDir := filepath.Join("..", "packages", packageName, challenge.ID, "submissions")

		// Check if submissions directory exists
		if _, err := os.Stat(submissionsDir); os.IsNotExist(err) {
			continue
		}

		// Read the submissions directory
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

				// Check if user has a solution file (either solution.go or solution-template.go)
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

				// Only count if not already counted
				if !userStats[username].challengesCompleted[challenge.ID] {
					userStats[username].completedCount++
					userStats[username].challengesCompleted[challenge.ID] = true

					// Update last submission time if this is more recent
					if modTime.After(userStats[username].lastSubmission) {
						userStats[username].lastSubmission = modTime
					}
				}
			}
		}
	}

	// Convert to leaderboard entries and sort
	for username, stats := range userStats {
		if stats.completedCount > 0 {
			entry := models.PackageScoreboardEntry{
				Username:    username,
				PackageName: packageName,
				ChallengeID: "", // Not specific to one challenge
				SubmittedAt: stats.lastSubmission,
				TestsPassed: stats.completedCount,
				TestsTotal:  len(challenges),
				IsSponsor:   sponsors[username],
			}
			leaderboard = append(leaderboard, entry)
		}
	}

	// Sort by completed count (descending), then by submission time (ascending for earliest)
	sort.Slice(leaderboard, func(i, j int) bool {
		if leaderboard[i].TestsPassed != leaderboard[j].TestsPassed {
			return leaderboard[i].TestsPassed > leaderboard[j].TestsPassed
		}
		return leaderboard[i].SubmittedAt.Before(leaderboard[j].SubmittedAt)
	})

	return leaderboard
}

// userPackageStats helper struct for collecting user statistics
type userPackageStats struct {
	username            string
	completedCount      int
	lastSubmission      time.Time
	challengesCompleted map[string]bool
}
