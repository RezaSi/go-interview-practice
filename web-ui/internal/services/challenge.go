package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"web-ui/internal/models"
)

// ChallengeService handles challenge-related operations
type ChallengeService struct {
	challenges models.ChallengeMap
}

// NewChallengeService creates a new challenge service
func NewChallengeService() *ChallengeService {
	return &ChallengeService{
		challenges: make(models.ChallengeMap),
	}
}

// LoadChallenges loads all challenges from the filesystem
func (cs *ChallengeService) LoadChallenges() error {
	// Find challenge directories (challenge-1, challenge-2, etc.)
	challengeDirs, err := filepath.Glob("../challenge-*")
	if err != nil {
		return fmt.Errorf("failed to find challenge directories: %v", err)
	}

	for _, dir := range challengeDirs {
		// Extract challenge number
		re := regexp.MustCompile(`challenge-(\d+)`)
		match := re.FindStringSubmatch(dir)
		if len(match) < 2 {
			continue
		}

		id, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}

		challenge, err := cs.loadSingleChallenge(id, dir)
		if err != nil {
			log.Printf("Warning: Could not load challenge %d: %v", id, err)
			continue
		}

		cs.challenges[id] = challenge
	}

	log.Printf("Loaded %d challenges", len(cs.challenges))
	return nil
}

// loadSingleChallenge loads a single challenge from a directory
func (cs *ChallengeService) loadSingleChallenge(id int, dir string) (*models.Challenge, error) {
	// Read README.md for title and description
	readmePath := filepath.Join(dir, "README.md")
	readmeContent, err := ioutil.ReadFile(readmePath)
	if err != nil {
		return nil, fmt.Errorf("could not read README: %v", err)
	}

	// Extract title from README (first heading)
	title := cs.extractTitle(string(readmeContent), id)

	// Determine difficulty level
	difficulty := cs.determineDifficulty(id)

	// Read solution template
	templatePath := filepath.Join(dir, "solution-template.go")
	templateContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("could not read solution template: %v", err)
	}

	// Read test file
	testPath := filepath.Join(dir, "solution-template_test.go")
	testContent, err := ioutil.ReadFile(testPath)
	if err != nil {
		log.Printf("Warning: Could not read test file for challenge %d: %v", id, err)
	}

	// Read learning materials if available
	learningPath := filepath.Join(dir, "learning.md")
	learningContent := []byte("*No learning materials available for this challenge yet.*")
	if learningFileContent, err := ioutil.ReadFile(learningPath); err == nil {
		learningContent = learningFileContent
	}

	// Read hints if available
	hintsPath := filepath.Join(dir, "hints.md")
	hintsContent := []byte("*No hints available for this challenge yet.*")
	if hintsFileContent, err := ioutil.ReadFile(hintsPath); err == nil {
		hintsContent = hintsFileContent
	}

	// Create challenge
	challenge := &models.Challenge{
		ID:                id,
		Title:             title,
		Description:       cs.filterWebUIDescription(string(readmeContent)),
		Difficulty:        difficulty,
		Template:          string(templateContent),
		TestFile:          string(testContent),
		LearningMaterials: string(learningContent),
		Hints:             string(hintsContent),
	}

	return challenge, nil
}

// extractTitle extracts the title from README content
func (cs *ChallengeService) extractTitle(readmeContent string, id int) string {
	titleRe := regexp.MustCompile(`#\s+(.+)`)
	titleMatch := titleRe.FindStringSubmatch(readmeContent)

	if len(titleMatch) >= 2 {
		title := titleMatch[1]
		// Clean up the title - remove "Challenge X: " prefix if present
		cleanTitle := regexp.MustCompile(`^Challenge\s+\d+:\s+`).ReplaceAllString(title, "")
		return cleanTitle
	}

	return fmt.Sprintf("Challenge %d", id)
}

// determineDifficulty determines the difficulty level based on challenge ID
func (cs *ChallengeService) determineDifficulty(id int) string {
	switch {
	case id <= 3 || id == 6 || id == 18 || id == 21 || id == 22:
		return "Beginner"
	case id == 4 || id == 5 || id == 7 || id == 10 || id == 13 || id == 14 || id == 16 || id == 17 || id == 19 || id == 20 || id == 23 || id == 27 || id == 30:
		return "Intermediate"
	default:
		return "Advanced"
	}
}

// filterWebUIDescription removes manual instructions that are not relevant for web-ui users
func (cs *ChallengeService) filterWebUIDescription(content string) string {
	lines := strings.Split(content, "\n")
	var filteredLines []string
	skipSection := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Check if we're starting a section to skip
		if strings.HasPrefix(trimmedLine, "## Instructions") ||
			strings.HasPrefix(trimmedLine, "## Testing Your Solution Locally") ||
			(strings.HasPrefix(trimmedLine, "## ") && strings.Contains(strings.ToLower(trimmedLine), "testing") && strings.Contains(strings.ToLower(trimmedLine), "locally")) {
			skipSection = true
			continue
		}

		// Check if we're starting a new section (stop skipping)
		if skipSection && strings.HasPrefix(trimmedLine, "## ") {
			skipSection = false
		}

		// Add line if we're not skipping
		if !skipSection {
			filteredLines = append(filteredLines, line)
		}
	}

	// Clean up extra blank lines at the end
	for len(filteredLines) > 0 && strings.TrimSpace(filteredLines[len(filteredLines)-1]) == "" {
		filteredLines = filteredLines[:len(filteredLines)-1]
	}

	return strings.Join(filteredLines, "\n")
}

// GetChallenges returns all challenges
func (cs *ChallengeService) GetChallenges() models.ChallengeMap {
	return cs.challenges
}

// GetChallenge returns a specific challenge by ID
func (cs *ChallengeService) GetChallenge(id int) (*models.Challenge, bool) {
	challenge, exists := cs.challenges[id]
	return challenge, exists
}
