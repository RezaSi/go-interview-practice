package handlers

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
}

// NewWebHandler creates a new web handler
func NewWebHandler(
	content embed.FS,
	challengeService *services.ChallengeService,
	scoreboardService *services.ScoreboardService,
	userService *services.UserService,
) *WebHandler {
	return &WebHandler{
		content:           content,
		challengeService:  challengeService,
		scoreboardService: scoreboardService,
		userService:       userService,
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

	// Get the username from cookie if available
	username := h.getUsernameFromCookie(r)

	// Get user attempts if username is set
	var userAttempt *models.UserAttemptedChallenges
	if username != "" {
		userAttempt = h.userService.GetUserAttempts(username, h.challengeService.GetChallenges())
	}

	data := struct {
		Challenges   []*models.Challenge
		Username     string
		UserAttempts *models.UserAttemptedChallenges
	}{
		Challenges:   challengeList,
		Username:     username,
		UserAttempts: userAttempt,
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
