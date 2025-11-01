package shared

import "github.com/golang-jwt/jwt/v5"

// User represents the user data structure for authentication
type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Claims represents JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthResponse represents the response from authentication operations
type AuthResponse struct {
	Token    string `json:"token,omitempty"`
	Username string `json:"username,omitempty"`
	Message  string `json:"message"`
	Error    string `json:"error,omitempty"`
}

// UploadResponse represents the response from upload operations
type UploadResponse struct {
	Message  string `json:"message"`
	ImageURL string `json:"imageUrl,omitempty"`
	Filename string `json:"filename,omitempty"`
	Error    string `json:"error,omitempty"`
}

// ProfileResponse represents the response from profile operations
type ProfileResponse struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Error    string `json:"error,omitempty"`
}