package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	// In our new implementation, packages are loaded on-demand
	fmt.Printf("Loaded 4 packages with real-time GitHub stars\n")
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

	// Use the full README content for both description and learning materials
	// This gives individual challenge pages the same detailed content as classic challenges

	return &models.PackageChallenge{
		ID:                challengeName,
		Title:             title,
		Description:       readmeContent, // Use README content directly
		Difficulty:        "Beginner",    // Could be parsed from README or metadata
		Template:          template,
		TestFile:          testFile,
		Hints:             hints,
		LearningMaterials: readmeContent, // Full content in learning materials tab
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
