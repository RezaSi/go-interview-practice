package services

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"web-ui/internal/models"
)

// UserService handles user-related operations
type UserService struct {
	userAttempts models.UserAttemptsMap
	mutex        sync.RWMutex
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		userAttempts: make(models.UserAttemptsMap),
	}
}

// LoadUserAttempts checks the filesystem for submission directories
func (us *UserService) LoadUserAttempts(username string, challenges models.ChallengeMap) *models.UserAttemptedChallenges {
	// Check cache with read lock
	us.mutex.RLock()
	if attempts, ok := us.userAttempts[username]; ok {
		us.mutex.RUnlock()
		return attempts
	}
	us.mutex.RUnlock()

	// Create new tracking structure
	userAttempt := &models.UserAttemptedChallenges{
		Username:     username,
		AttemptedIDs: make(map[int]bool),
		Scores:       make(map[int]int),
	}

	// Scan all challenge directories for this user's submissions
	for id := range challenges {
		if us.hasUserSubmission(username, id) {
			userAttempt.AttemptedIDs[id] = true
			// Calculate score based on test results
			score := us.calculateScore(username, id)
			userAttempt.Scores[id] = score
		}
	}

	// Cache the results with write lock
	us.mutex.Lock()
	us.userAttempts[username] = userAttempt
	us.mutex.Unlock()
	return userAttempt
}

// hasUserSubmission checks if a user has a submission for a challenge
func (us *UserService) hasUserSubmission(username string, challengeID int) bool {
	// Try different path formats to handle potential path issues
	// Absolute path
	submissionDir := filepath.Join("..", fmt.Sprintf("challenge-%d", challengeID), "submissions", username)
	submissionFile := filepath.Join(submissionDir, "solution-template.go")

	// Check if the file exists
	if _, err := os.Stat(submissionFile); err == nil {
		return true
	}

	// Alternative path (direct from workspace root)
	altSubmissionFile := filepath.Join(fmt.Sprintf("challenge-%d", challengeID), "submissions", username, "solution-template.go")

	if _, err := os.Stat(altSubmissionFile); err == nil {
		return true
	}

	return false
}

// GetExistingSolution returns the content of an existing solution file if it exists
func (us *UserService) GetExistingSolution(username string, challengeID int) string {
	if username == "" {
		return ""
	}

	// Try different path formats
	// First try the relative path from web-ui
	submissionFile := filepath.Join("..", fmt.Sprintf("challenge-%d", challengeID), "submissions", username, "solution-template.go")
	content, err := ioutil.ReadFile(submissionFile)
	if err == nil {
		return string(content)
	}

	// Try alternative path from root directory
	altSubmissionFile := filepath.Join(fmt.Sprintf("challenge-%d", challengeID), "submissions", username, "solution-template.go")
	content, err = ioutil.ReadFile(altSubmissionFile)
	if err == nil {
		return string(content)
	}

	return ""
}

// RefreshUserAttempts clears the cache for a user and reloads their attempts
func (us *UserService) RefreshUserAttempts(username string, challenges models.ChallengeMap) *models.UserAttemptedChallenges {
	// Clear cache with write lock
	us.mutex.Lock()
	delete(us.userAttempts, username)
	us.mutex.Unlock()
	// Reload and return
	return us.LoadUserAttempts(username, challenges)
}

// GetUserAttempts returns the cached user attempts or loads them if not cached
func (us *UserService) GetUserAttempts(username string, challenges models.ChallengeMap) *models.UserAttemptedChallenges {
	us.mutex.RLock()
	if attempts, ok := us.userAttempts[username]; ok {
		us.mutex.RUnlock()
		return attempts
	}
	us.mutex.RUnlock()
	return us.LoadUserAttempts(username, challenges)
}

// calculateScore calculates the score for a user's submission for a challenge
func (us *UserService) calculateScore(username string, challengeID int) int {
	// Read the scoreboard file for this challenge
	scoreboardPath := filepath.Join("..", fmt.Sprintf("challenge-%d", challengeID), "SCOREBOARD.md")
	content, err := ioutil.ReadFile(scoreboardPath)
	if err != nil {
		// Try alternative path
		scoreboardPath = filepath.Join(fmt.Sprintf("challenge-%d", challengeID), "SCOREBOARD.md")
		content, err = ioutil.ReadFile(scoreboardPath)
		if err != nil {
			// No scoreboard file, return default score
			return 50
		}
	}

	scoreboardContent := string(content)

	// Parse the scoreboard to find this user's results
	lines := strings.Split(scoreboardContent, "\n")
	for _, line := range lines {
		// Look for lines that contain the username
		if strings.Contains(line, username) && strings.Contains(line, "|") {
			// Parse the table row: | Username | Passed Tests | Total Tests |
			parts := strings.Split(line, "|")
			if len(parts) >= 4 {
				passedStr := strings.TrimSpace(parts[2])
				totalStr := strings.TrimSpace(parts[3])

				passed, err1 := strconv.Atoi(passedStr)
				total, err2 := strconv.Atoi(totalStr)

				if err1 == nil && err2 == nil && total > 0 {
					// Calculate percentage score
					score := (passed * 100) / total
					return score
				}
			}
		}
	}

	// User not found in scoreboard, return 0
	return 0
}
