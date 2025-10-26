package interfaces

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gocleanarchitecture/entities"
	"gocleanarchitecture/frameworks/auth"

	"github.com/google/uuid"
)

// OAuth2Controller handles OAuth2 authentication
type OAuth2Controller struct {
	GoogleProvider *auth.OAuth2Provider
	GitHubProvider *auth.OAuth2Provider
	AuthUseCase    AuthUseCase
}

// NewOAuth2Controller creates a new OAuth2Controller
func NewOAuth2Controller(googleProvider, githubProvider *auth.OAuth2Provider, authUseCase AuthUseCase) *OAuth2Controller {
	return &OAuth2Controller{
		GoogleProvider: googleProvider,
		GitHubProvider: githubProvider,
		AuthUseCase:    authUseCase,
	}
}

// generateState generates a random state string for OAuth2
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// InitiateGoogleLogin redirects to Google OAuth2 consent screen
func (c *OAuth2Controller) InitiateGoogleLogin(w http.ResponseWriter, r *http.Request) {
	if c.GoogleProvider == nil {
		http.Error(w, "Google OAuth2 not configured", http.StatusServiceUnavailable)
		return
	}

	state, err := generateState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	// Store state in a cookie for verification (expires in 10 minutes)
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	authURL := c.GoogleProvider.GetAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the Google OAuth2 callback
func (c *OAuth2Controller) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if c.GoogleProvider == nil {
		http.Error(w, "Google OAuth2 not configured", http.StatusServiceUnavailable)
		return
	}

	// Verify state
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, "Missing state cookie", http.StatusBadRequest)
		return
	}

	stateParam := r.URL.Query().Get("state")
	if stateParam != stateCookie.Value {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Clear the state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})

	// Exchange code for token
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	token, err := c.GoogleProvider.ExchangeCode(context.Background(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange code: %v", err), http.StatusInternalServerError)
		return
	}

	// Get user info from Google
	userInfo, err := c.GoogleProvider.GetUserInfo(context.Background(), token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}

	// Process OAuth2 login (create or find existing user)
	user, jwtToken, err := c.processOAuth2Login(userInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process login: %v", err), http.StatusInternalServerError)
		return
	}

	// Return user and JWT token
	response := map[string]interface{}{
		"user":     user.Sanitize(),
		"token":    jwtToken,
		"provider": "google",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// InitiateGitHubLogin redirects to GitHub OAuth2 consent screen
func (c *OAuth2Controller) InitiateGitHubLogin(w http.ResponseWriter, r *http.Request) {
	if c.GitHubProvider == nil {
		http.Error(w, "GitHub OAuth2 not configured", http.StatusServiceUnavailable)
		return
	}

	state, err := generateState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	// Store state in a cookie for verification (expires in 10 minutes)
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	authURL := c.GitHubProvider.GetAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// GitHubCallback handles the GitHub OAuth2 callback
func (c *OAuth2Controller) GitHubCallback(w http.ResponseWriter, r *http.Request) {
	if c.GitHubProvider == nil {
		http.Error(w, "GitHub OAuth2 not configured", http.StatusServiceUnavailable)
		return
	}

	// Verify state
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, "Missing state cookie", http.StatusBadRequest)
		return
	}

	stateParam := r.URL.Query().Get("state")
	if stateParam != stateCookie.Value {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Clear the state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})

	// Exchange code for token
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	token, err := c.GitHubProvider.ExchangeCode(context.Background(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange code: %v", err), http.StatusInternalServerError)
		return
	}

	// Get user info from GitHub
	userInfo, err := c.GitHubProvider.GetUserInfo(context.Background(), token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}

	// Process OAuth2 login (create or find existing user)
	user, jwtToken, err := c.processOAuth2Login(userInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process login: %v", err), http.StatusInternalServerError)
		return
	}

	// Return user and JWT token
	response := map[string]interface{}{
		"user":     user.Sanitize(),
		"token":    jwtToken,
		"provider": "github",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// processOAuth2Login handles the user creation or retrieval after OAuth2 authentication
func (c *OAuth2Controller) processOAuth2Login(userInfo *auth.UserInfo) (*entities.User, string, error) {
	// Try to find existing user by email
	user, err := c.AuthUseCase.GetUserByEmail(userInfo.Email)

	var token string

	if err != nil || user == nil {
		// User doesn't exist, create a new one
		// Generate a username from email or name
		username := generateUsernameFromEmail(userInfo.Email)
		
		// Create random password (won't be used for OAuth logins)
		randomPassword := uuid.New().String()

		loginResp, err := c.AuthUseCase.Register(username, userInfo.Email, randomPassword, userInfo.Name)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create user: %w", err)
		}

		user = loginResp.User
		token = loginResp.Token

		// Update avatar URL if provided
		if userInfo.AvatarURL != "" {
			user, err = c.AuthUseCase.UpdateProfile(user.ID, user.FullName, user.Bio, userInfo.AvatarURL)
			if err != nil {
				// Log but don't fail - avatar update is optional
			}
		}
	} else {
		// User exists, update their avatar if provided and generate a new token
		if userInfo.AvatarURL != "" && userInfo.AvatarURL != user.AvatarURL {
			user, err = c.AuthUseCase.UpdateProfile(user.ID, user.FullName, user.Bio, userInfo.AvatarURL)
			if err != nil {
				// Log but don't fail - avatar update is optional
			}
		}

		// Generate a JWT token for the existing user
		token, err = c.AuthUseCase.GenerateTokenForUser(user.ID)
		if err != nil {
			return nil, "", fmt.Errorf("failed to generate token: %w", err)
		}
	}

	return user, token, nil
}

// generateUsernameFromEmail generates a username from an email address
func generateUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		username := strings.ToLower(parts[0])
		username = strings.ReplaceAll(username, ".", "_")
		username = strings.ReplaceAll(username, "+", "_")
		
		// Add a random suffix to avoid collisions
		suffix := uuid.New().String()[:8]
		return username + "_" + suffix
	}
	return "user_" + uuid.New().String()[:8]
}

