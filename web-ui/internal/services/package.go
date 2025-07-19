package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"web-ui/internal/models"
)

type PackageService struct {
	httpClient   *http.Client
	packagesPath string
}

func NewPackageService() *PackageService {
	return &PackageService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		packagesPath: "../packages", // Relative to web-ui directory
	}
}

type PackageMetadata struct {
	Name             string   `json:"name"`
	DisplayName      string   `json:"display_name"`
	Description      string   `json:"description"`
	Version          string   `json:"version"`
	GitHubURL        string   `json:"github_url"`
	DocumentationURL string   `json:"documentation_url"`
	Stars            int      `json:"stars"`
	Category         string   `json:"category"`
	Difficulty       string   `json:"difficulty"`
	Prerequisites    []string `json:"prerequisites"`
	LearningPath     []string `json:"learning_path"`
	Tags             []string `json:"tags"`
	EstimatedTime    string   `json:"estimated_time"`
	RealWorldUsage   []string `json:"real_world_usage"`
}

func (s *PackageService) LoadPackages() error {
	// This method is called to ensure packages are loaded
	// Load packages and count them for logging
	packages := s.GetPackages()
	fmt.Printf("Loaded %d packages with real-time GitHub stars\n", len(packages))
	return nil
}

func (s *PackageService) GetPackages() map[string]*models.Package {
	packages := make(map[string]*models.Package)

	// Read packages directory
	entries, err := os.ReadDir(s.packagesPath)
	if err != nil {
		fmt.Printf("Error reading packages directory: %v\n", err)
		return packages
	}

	for _, entry := range entries {
		if entry.IsDir() {
			packagePath := filepath.Join(s.packagesPath, entry.Name())
			if pkg := s.loadPackage(packagePath, entry.Name()); pkg != nil {
				packages[pkg.Name] = pkg
			}
		}
	}

	return packages
}

func (s *PackageService) loadPackage(packagePath, packageName string) *models.Package {
	// Ensure httpClient is initialized
	if s.httpClient == nil {
		s.httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	// Load package.json
	metadataPath := filepath.Join(packagePath, "package.json")
	metadataBytes, err := os.ReadFile(metadataPath)
	if err != nil {
		fmt.Printf("Error reading package.json for %s: %v\n", packageName, err)
		return nil
	}

	var metadata PackageMetadata
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		fmt.Printf("Error parsing package.json for %s: %v\n", packageName, err)
		return nil
	}

	// Fetch real-time GitHub stars
	stars := s.fetchGitHubStars(metadata.GitHubURL)
	if stars > 0 {
		metadata.Stars = stars
	}

	// Load challenge details dynamically
	challengeDetails := s.loadChallengeDetails(packagePath, metadata.LearningPath)

	return &models.Package{
		Name:             packageName,
		DisplayName:      metadata.DisplayName,
		Description:      metadata.Description,
		Version:          metadata.Version,
		GitHubURL:        metadata.GitHubURL,
		DocumentationURL: metadata.DocumentationURL,
		Stars:            metadata.Stars,
		Category:         metadata.Category,
		Difficulty:       metadata.Difficulty,
		Prerequisites:    metadata.Prerequisites,
		LearningPath:     metadata.LearningPath,
		Tags:             metadata.Tags,
		EstimatedTime:    metadata.EstimatedTime,
		RealWorldUsage:   metadata.RealWorldUsage,
		ChallengeDetails: challengeDetails,
	}
}

// loadChallengeDetails dynamically loads metadata for each challenge in the learning path
func (s *PackageService) loadChallengeDetails(packagePath string, learningPath []string) map[string]*models.ChallengeInfo {
	challengeDetails := make(map[string]*models.ChallengeInfo)

	for i, challengeID := range learningPath {
		challengePath := filepath.Join(packagePath, challengeID)

		// Check if challenge directory exists
		if _, err := os.Stat(challengePath); os.IsNotExist(err) {
			// Challenge doesn't exist yet, mark as coming soon
			challengeDetails[challengeID] = &models.ChallengeInfo{
				ID:            challengeID,
				Title:         s.generateTitleFromID(challengeID),
				Description:   "Coming soon",
				Difficulty:    s.inferDifficultyFromOrder(i),
				EstimatedTime: s.inferEstimatedTime(challengeID),
				Status:        "coming-soon",
				Order:         i + 1,
				Icon:          s.inferIconFromID(challengeID),
			}
			continue
		}

		// Try to load challenge metadata
		metadata := s.loadChallengeMetadata(challengePath)
		if metadata != nil {
			challengeDetails[challengeID] = &models.ChallengeInfo{
				ID:                  challengeID,
				Title:               metadata.Title,
				Description:         metadata.ShortDescription,
				Difficulty:          metadata.Difficulty,
				EstimatedTime:       metadata.EstimatedTime,
				LearningObjectives:  metadata.LearningObjectives,
				Prerequisites:       metadata.Prerequisites,
				Tags:                metadata.Tags,
				RealWorldConnection: metadata.RealWorldConnection,
				Icon:                metadata.Icon,
				Status:              "available",
				Order:               i + 1,
			}
		} else {
			// Fallback to generated metadata
			challengeDetails[challengeID] = &models.ChallengeInfo{
				ID:            challengeID,
				Title:         s.generateTitleFromID(challengeID),
				Description:   s.generateDescriptionFromReadme(challengePath),
				Difficulty:    s.inferDifficultyFromOrder(i),
				EstimatedTime: s.inferEstimatedTime(challengeID),
				Status:        "available",
				Order:         i + 1,
				Icon:          s.inferIconFromID(challengeID),
			}
		}
	}

	return challengeDetails
}

// loadChallengeMetadata loads metadata from challenge directory
func (s *PackageService) loadChallengeMetadata(challengePath string) *models.ChallengeMetadata {
	metadataPath := filepath.Join(challengePath, "metadata.json")

	// Check if metadata.json exists
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return nil
	}

	metadataBytes, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil
	}

	var metadata models.ChallengeMetadata
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return nil
	}

	return &metadata
}

// Helper functions for generating metadata when not available

func (s *PackageService) generateTitleFromID(challengeID string) string {
	// Remove "challenge-N-" prefix and convert to title case
	parts := strings.Split(challengeID, "-")
	if len(parts) >= 3 {
		title := strings.Join(parts[2:], " ")
		return strings.Title(strings.ReplaceAll(title, "-", " "))
	}
	return strings.Title(strings.ReplaceAll(challengeID, "-", " "))
}

func (s *PackageService) generateDescriptionFromReadme(challengePath string) string {
	readmePath := filepath.Join(challengePath, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return "Challenge content available"
	}

	// Extract first paragraph as description
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "---") {
			if len(line) > 120 {
				return line[:120] + "..."
			}
			return line
		}
	}

	return "Learn practical implementation techniques"
}

func (s *PackageService) inferDifficultyFromOrder(order int) string {
	switch {
	case order < 2:
		return "Beginner"
	case order < 4:
		return "Intermediate"
	default:
		return "Advanced"
	}
}

func (s *PackageService) inferEstimatedTime(challengeID string) string {
	switch {
	case strings.Contains(challengeID, "basic"):
		return "30-45 min"
	case strings.Contains(challengeID, "middleware"):
		return "45-60 min"
	case strings.Contains(challengeID, "validation") || strings.Contains(challengeID, "error"):
		return "60-90 min"
	case strings.Contains(challengeID, "auth"):
		return "90-120 min"
	case strings.Contains(challengeID, "file") || strings.Contains(challengeID, "upload"):
		return "60-90 min"
	default:
		return "45-60 min"
	}
}

func (s *PackageService) inferIconFromID(challengeID string) string {
	switch {
	case strings.Contains(challengeID, "basic") || strings.Contains(challengeID, "routing"):
		return "bi-play-circle"
	case strings.Contains(challengeID, "middleware"):
		return "bi-layers"
	case strings.Contains(challengeID, "validation") || strings.Contains(challengeID, "error"):
		return "bi-shield-check"
	case strings.Contains(challengeID, "auth"):
		return "bi-person-lock"
	case strings.Contains(challengeID, "file") || strings.Contains(challengeID, "upload"):
		return "bi-cloud-upload"
	case strings.Contains(challengeID, "database") || strings.Contains(challengeID, "db"):
		return "bi-database"
	case strings.Contains(challengeID, "cli"):
		return "bi-terminal"
	default:
		return "bi-code-slash"
	}
}

func (s *PackageService) loadChallenges(packagePath string) []models.PackageChallenge {
	var challenges []models.PackageChallenge

	// Read challenge directories
	entries, err := os.ReadDir(packagePath)
	if err != nil {
		return challenges
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "challenge-") {
			challengePath := filepath.Join(packagePath, entry.Name())
			if challenge := s.loadChallenge(challengePath, entry.Name()); challenge != nil {
				challenges = append(challenges, *challenge)
			}
		}
	}

	return challenges
}

func (s *PackageService) loadChallenge(challengePath, challengeName string) *models.PackageChallenge {
	// Extract challenge title from directory name
	parts := strings.Split(challengeName, "-")
	if len(parts) < 3 {
		return nil
	}

	title := strings.Join(parts[2:], " ")
	title = strings.Title(strings.ReplaceAll(title, "-", " "))

	// Load README.md for full content
	readmeContent := s.readFileContent(filepath.Join(challengePath, "README.md"))
	if readmeContent == "" {
		readmeContent = "Challenge content not available"
	}

	// For individual challenge pages, use full content like classic challenges
	// For package listing, templates can extract brief descriptions as needed

	// Load solution template
	template := s.readFileContent(filepath.Join(challengePath, "solution-template.go"))
	if template == "" {
		template = "// Solution template not available"
	}

	// Load test file
	testFile := s.readFileContent(filepath.Join(challengePath, "solution-template_test.go"))
	if testFile == "" {
		testFile = "// Test file not available"
	}

	// Load hints
	hints := s.readFileContent(filepath.Join(challengePath, "hints.md"))
	if hints == "" {
		hints = "No hints available for this challenge."
	}

	// Load learning materials from learning.md (same as classic challenges)
	learningMaterials := s.readFileContent(filepath.Join(challengePath, "learning.md"))
	if learningMaterials == "" {
		learningMaterials = "*No learning materials available for this challenge yet.*"
	}

	// Determine difficulty - try to load from metadata first, then infer from challenge name
	difficulty := "Beginner" // default fallback

	// Try to load metadata.json for difficulty
	metadata := s.loadChallengeMetadata(challengePath)
	if metadata != nil && metadata.Difficulty != "" {
		difficulty = metadata.Difficulty
	} else {
		// Infer difficulty from challenge name/order
		// Extract challenge number from name (e.g., "challenge-1-basic-routing" -> 1)
		if len(parts) >= 2 {
			if challengeNum, err := strconv.Atoi(parts[1]); err == nil {
				difficulty = s.inferDifficultyFromOrder(challengeNum - 1) // Convert to 0-based index
			}
		}
	}

	return &models.PackageChallenge{
		ID:                challengeName,
		Title:             title,
		Description:       readmeContent, // Use README content for description
		Difficulty:        difficulty,
		Template:          template,
		TestFile:          testFile,
		Hints:             hints,
		LearningMaterials: learningMaterials, // Use learning.md for learning materials tab
	}
}

func (s *PackageService) readFileContent(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return string(content)
}

func (s *PackageService) fetchGitHubStars(githubURL string) int {
	// Add more robust nil checking
	if s == nil || s.httpClient == nil || githubURL == "" {
		return 0
	}

	// Extract owner/repo from GitHub URL
	parts := strings.Split(githubURL, "/")
	if len(parts) < 2 {
		return 0
	}

	repo := fmt.Sprintf("%s/%s", parts[len(parts)-2], parts[len(parts)-1])
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s", repo)

	resp, err := s.httpClient.Get(apiURL)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		fmt.Printf("GitHub API returned status 403 for %s\n", githubURL)
		return 0
	}

	if resp.StatusCode != 200 {
		return 0
	}

	var repoData struct {
		StargazersCount int `json:"stargazers_count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&repoData); err != nil {
		return 0
	}

	return repoData.StargazersCount
}

func (s *PackageService) GetPackage(packageID string) (*models.Package, error) {
	packages := s.GetPackages()
	if pkg, exists := packages[packageID]; exists {
		return pkg, nil
	}
	return nil, fmt.Errorf("package %s not found", packageID)
}

func (s *PackageService) GetChallenge(packageID, challengeID string) *models.PackageChallenge {
	// Load challenge directly from filesystem
	packagePath := filepath.Join(s.packagesPath, packageID)
	challengePath := filepath.Join(packagePath, challengeID)

	// Check if challenge directory exists
	if _, err := os.Stat(challengePath); os.IsNotExist(err) {
		return nil
	}

	return s.loadChallenge(challengePath, challengeID)
}

func (s *PackageService) GetPackageChallenges(packageID string) (map[string]*models.PackageChallenge, error) {
	packagePath := filepath.Join(s.packagesPath, packageID)

	// Check if package directory exists
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("package %s not found", packageID)
	}

	challengesList := s.loadChallenges(packagePath)
	challenges := make(map[string]*models.PackageChallenge)

	for _, challenge := range challengesList {
		// Create a copy to avoid the loop variable reference issue
		challengeCopy := challenge
		challenges[challenge.ID] = &challengeCopy
	}

	return challenges, nil
}

func (s *PackageService) GetPackageChallenge(packageID, challengeID string) (*models.PackageChallenge, error) {
	// Load challenge directly from filesystem
	packagePath := filepath.Join(s.packagesPath, packageID)
	challengePath := filepath.Join(packagePath, challengeID)

	// Check if challenge directory exists
	if _, err := os.Stat(challengePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("challenge %s not found in package %s", challengeID, packageID)
	}

	challenge := s.loadChallenge(challengePath, challengeID)
	if challenge == nil {
		return nil, fmt.Errorf("failed to load challenge %s from package %s", challengeID, packageID)
	}

	// Set the package name for the challenge
	challenge.PackageName = packageID

	return challenge, nil
}
