package models

import (
	"time"
)

// Package represents a Go package with associated challenges
type Package struct {
	Name             string                    `json:"name"`
	DisplayName      string                    `json:"display_name"`
	Description      string                    `json:"description"`
	Version          string                    `json:"version"`
	GitHubURL        string                    `json:"github_url"`
	DocumentationURL string                    `json:"documentation_url"`
	Stars            int                       `json:"stars"`
	Category         string                    `json:"category"`
	Difficulty       string                    `json:"difficulty"`
	Prerequisites    []string                  `json:"prerequisites"`
	LearningPath     []string                  `json:"learning_path"`
	Tags             []string                  `json:"tags"`
	EstimatedTime    string                    `json:"estimated_time"`
	RealWorldUsage   []string                  `json:"real_world_usage"`
	ChallengeDetails map[string]*ChallengeInfo `json:"challenge_details,omitempty"` // Dynamic challenge metadata
}

// ChallengeInfo contains metadata about each challenge in the learning path
type ChallengeInfo struct {
	ID                  string   `json:"id"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	Difficulty          string   `json:"difficulty"`
	EstimatedTime       string   `json:"estimated_time"`
	LearningObjectives  []string `json:"learning_objectives"`
	Prerequisites       []string `json:"prerequisites"`
	Tags                []string `json:"tags"`
	RealWorldConnection string   `json:"real_world_connection"`
	Icon                string   `json:"icon,omitempty"`   // Bootstrap icon class
	Status              string   `json:"status,omitempty"` // "available", "coming-soon", etc.
	Order               int      `json:"order"`            // Order in learning path
}

// ChallengeMetadata represents metadata that can be loaded from challenge directories
type ChallengeMetadata struct {
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	ShortDescription    string   `json:"short_description"` // Brief description for cards
	Difficulty          string   `json:"difficulty"`
	EstimatedTime       string   `json:"estimated_time"`
	LearningObjectives  []string `json:"learning_objectives"`
	Prerequisites       []string `json:"prerequisites"`
	Tags                []string `json:"tags"`
	RealWorldConnection string   `json:"real_world_connection"`
	Requirements        []string `json:"requirements"`
	BonusPoints         []string `json:"bonus_points"`
	Icon                string   `json:"icon,omitempty"`
	Order               int      `json:"order"`
}

// PackageChallenge represents a challenge specific to a package
type PackageChallenge struct {
	ID                  string   `json:"id"`           // e.g., "challenge-1-basic-routing"
	PackageName         string   `json:"package_name"` // e.g., "gin"
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	ShortDescription    string   `json:"short_description"` // Brief description for cards
	Difficulty          string   `json:"difficulty"`
	LearningObjectives  []string `json:"learning_objectives"`
	Template            string   `json:"template"`
	TestFile            string   `json:"testFile"`
	LearningMaterials   string   `json:"learningMaterials"`
	Hints               string   `json:"hints"`
	Requirements        []string `json:"requirements"`
	BonusPoints         []string `json:"bonus_points"`
	RealWorldConnection string   `json:"real_world_connection"`
	EstimatedTime       string   `json:"estimated_time"`
	Tags                []string `json:"tags"`
	Prerequisites       []string `json:"prerequisites"`
	Icon                string   `json:"icon,omitempty"`
	Order               int      `json:"order"`
	Status              string   `json:"status,omitempty"` // "available", "coming-soon", etc.
}

// PackageSubmission represents a user's submitted solution for a package challenge
type PackageSubmission struct {
	Username    string    `json:"username"`
	PackageName string    `json:"package_name"`
	ChallengeID string    `json:"challenge_id"`
	Code        string    `json:"code"`
	SubmittedAt time.Time `json:"submitted_at"`
	Passed      bool      `json:"passed"`
	TestOutput  string    `json:"test_output"`
	ExecutionMs int64     `json:"execution_ms"`
	TestsPassed int       `json:"tests_passed"`
	TestsTotal  int       `json:"tests_total"`
}

// PackageProgress tracks user progress in package learning paths
type PackageProgress struct {
	Username            string        `json:"username"`
	PackageName         string        `json:"package_name"`
	CompletedChallenges []string      `json:"completed_challenges"`
	InProgress          string        `json:"in_progress"`
	StartedAt           time.Time     `json:"started_at"`
	LastActivity        time.Time     `json:"last_activity"`
	TotalTime           time.Duration `json:"total_time"`
	Achievements        []string      `json:"achievements"`
	Score               int           `json:"score"`
}

// PackageScoreboardEntry represents an entry in the package scoreboard
type PackageScoreboardEntry struct {
	Username    string    `json:"username"`
	PackageName string    `json:"package_name"`
	ChallengeID string    `json:"challenge_id"`
	SubmittedAt time.Time `json:"submitted_at"`
	ExecutionMs int64     `json:"execution_ms"`
	TestsPassed int       `json:"tests_passed"`
	TestsTotal  int       `json:"tests_total"`
}

// Type aliases for collections
type PackageMap map[string]*Package
type PackageChallengeMap map[string]map[string]*PackageChallenge // package -> challenge_id -> challenge
type PackageSubmissionMap map[string][]PackageSubmission         // package -> submissions
type PackageProgressMap map[string]map[string]*PackageProgress   // username -> package -> progress
type PackageScoreboardMap map[string][]PackageScoreboardEntry    // package -> entries
