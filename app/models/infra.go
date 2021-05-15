package models

import (
	"encoding/json"
	"time"
)

// SystemSettings is the system-wide settings
type SystemSettings struct {
	Mode            string
	BuildTime       string
	Version         string
	Environment     string
	GoogleAnalytics string
	Compiler        string
	Domain          string
	HasLegal        bool
}

// Notification is the system generated notification entity
type Notification struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Link      string    `json:"link" db:"link"`
	Read      bool      `json:"read" db:"read"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// OAuthConfig is the configuration of a custom OAuth provider
type OAuthConfig struct {
	ID                int
	Provider          string
	DisplayName       string
	LogoBlobKey       string
	Status            int
	ClientID          string
	ClientSecret      string
	AuthorizeURL      string
	TokenURL          string
	ProfileURL        string
	Scope             string
	JSONUserIDPath    string
	JSONUserNamePath  string
	JSONUserEmailPath string
}

// MarshalJSON returns the JSON encoding of OAuthConfig
func (o OAuthConfig) MarshalJSON() ([]byte, error) {
	secret := "..."
	if len(o.ClientSecret) >= 10 {
		secret = o.ClientSecret[0:3] + "..." + o.ClientSecret[len(o.ClientSecret)-3:]
	}
	return json.Marshal(map[string]interface{}{
		"id":                o.ID,
		"provider":          o.Provider,
		"displayName":       o.DisplayName,
		"logoBlobKey":       o.LogoBlobKey,
		"status":            o.Status,
		"clientID":          o.ClientID,
		"clientSecret":      secret,
		"authorizeURL":      o.AuthorizeURL,
		"tokenURL":          o.TokenURL,
		"profileURL":        o.ProfileURL,
		"scope":             o.Scope,
		"jsonUserIDPath":    o.JSONUserIDPath,
		"jsonUserNamePath":  o.JSONUserNamePath,
		"jsonUserEmailPath": o.JSONUserEmailPath,
	})
}
