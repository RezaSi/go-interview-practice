package models

import (
	"time"
)

// Challenge represents a coding challenge
type Challenge struct {
	ID                int    `json:"id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Difficulty        string `json:"difficulty"`
	Template          string `json:"template"`
	TestFile          string `json:"testFile"`
	LearningMaterials string `json:"learningMaterials"`
	Hints             string `json:"hints"`
}

// Submission represents a user's submitted solution
type Submission struct {
	Username    string    `json:"username"`
	ChallengeID int       `json:"challengeId"`
	Code        string    `json:"code"`
	SubmittedAt time.Time `json:"submittedAt"`
	Passed      bool      `json:"passed"`
	TestOutput  string    `json:"testOutput"`
	ExecutionMs int64     `json:"executionMs"`
}

// ScoreboardEntry represents an entry in the scoreboard
type ScoreboardEntry struct {
	Username    string    `json:"username"`
	ChallengeID int       `json:"challengeId"`
	SubmittedAt time.Time `json:"submittedAt"`
}

// UserAttemptedChallenges tracks attempted challenges by username
type UserAttemptedChallenges struct {
	Username     string       `json:"username"`
	AttemptedIDs map[int]bool `json:"attemptedIds"`
	Scores       map[int]int  `json:"scores"` // Scores (0-100) for each attempted challenge
}

// ChallengeMap is a type alias for the challenges map
type ChallengeMap map[int]*Challenge

// ScoreboardMap is a type alias for the scoreboards map
type ScoreboardMap map[int][]ScoreboardEntry

// UserAttemptsMap is a type alias for user attempts tracking
type UserAttemptsMap map[string]*UserAttemptedChallenges
