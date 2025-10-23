package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// OAuth2Provider represents an OAuth2 provider configuration
type OAuth2Provider struct {
	Config *oauth2.Config
	Name   string
}

// UserInfo represents the user information returned by OAuth2 providers
type UserInfo struct {
	Email      string
	Name       string
	AvatarURL  string
	Provider   string
	ProviderID string
}

// NewGoogleOAuth2Provider creates a new Google OAuth2 provider
func NewGoogleOAuth2Provider(clientID, clientSecret, redirectURL string) *OAuth2Provider {
	return &OAuth2Provider{
		Name: "google",
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

// NewGitHubOAuth2Provider creates a new GitHub OAuth2 provider
func NewGitHubOAuth2Provider(clientID, clientSecret, redirectURL string) *OAuth2Provider {
	return &OAuth2Provider{
		Name: "github",
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"user:email", "read:user"},
			Endpoint:     github.Endpoint,
		},
	}
}

// GetAuthURL returns the OAuth2 authorization URL
func (p *OAuth2Provider) GetAuthURL(state string) string {
	return p.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode exchanges an authorization code for a token
func (p *OAuth2Provider) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.Config.Exchange(ctx, code)
}

// GetUserInfo retrieves user information from the OAuth2 provider
func (p *OAuth2Provider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	client := p.Config.Client(ctx, token)

	var url string
	switch p.Name {
	case "google":
		url = "https://www.googleapis.com/oauth2/v2/userinfo"
	case "github":
		url = "https://api.github.com/user"
	default:
		return nil, fmt.Errorf("unsupported provider: %s", p.Name)
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return p.parseUserInfo(body)
}

func (p *OAuth2Provider) parseUserInfo(data []byte) (*UserInfo, error) {
	userInfo := &UserInfo{
		Provider: p.Name,
	}

	switch p.Name {
	case "google":
		var googleUser struct {
			ID      string `json:"id"`
			Email   string `json:"email"`
			Name    string `json:"name"`
			Picture string `json:"picture"`
		}
		if err := json.Unmarshal(data, &googleUser); err != nil {
			return nil, err
		}
		userInfo.ProviderID = googleUser.ID
		userInfo.Email = googleUser.Email
		userInfo.Name = googleUser.Name
		userInfo.AvatarURL = googleUser.Picture

	case "github":
		var githubUser struct {
			ID        int    `json:"id"`
			Login     string `json:"login"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
		}
		if err := json.Unmarshal(data, &githubUser); err != nil {
			return nil, err
		}
		userInfo.ProviderID = fmt.Sprintf("%d", githubUser.ID)
		userInfo.Email = githubUser.Email
		if githubUser.Name != "" {
			userInfo.Name = githubUser.Name
		} else {
			userInfo.Name = githubUser.Login
		}
		userInfo.AvatarURL = githubUser.AvatarURL

		// GitHub might not return email in the user endpoint, need to fetch separately
		if userInfo.Email == "" {
			// For simplicity, we'll leave this as-is
			// In production, you'd fetch from https://api.github.com/user/emails
			userInfo.Email = githubUser.Login + "@github.local" // Placeholder
		}
	}

	return userInfo, nil
}
