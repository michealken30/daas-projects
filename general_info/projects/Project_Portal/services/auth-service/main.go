package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"shared"
)

var users = make(map[string]string) // username -> hashed password

func main() {
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Auth routes
	r.POST("/register", register)
	r.POST("/login", login)
	r.POST("/register-form", registerForm)
	r.POST("/login-form", loginForm)

	r.Run(":8082") // Auth service on port 8082
}

func register(c *gin.Context) {
	var user shared.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, shared.AuthResponse{Error: err.Error()})
		return
	}

	// Check if user already exists
	if _, exists := users[user.Username]; exists {
		c.JSON(http.StatusConflict, shared.AuthResponse{Error: "User already exists"})
		return
	}

	// Hash password
	hashedPassword := shared.HashPassword(user.Password)
	users[user.Username] = hashedPassword

	c.JSON(http.StatusCreated, shared.AuthResponse{Message: "User registered successfully"})
}

func registerForm(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, shared.AuthResponse{Error: "Username and password are required"})
		return
	}

	// Check if user already exists
	if _, exists := users[username]; exists {
		c.JSON(http.StatusConflict, shared.AuthResponse{Error: "User already exists"})
		return
	}

	// Hash password
	hashedPassword := shared.HashPassword(password)
	users[username] = hashedPassword

	c.JSON(http.StatusCreated, shared.AuthResponse{Message: "User registered successfully"})
}

func login(c *gin.Context) {
	var user shared.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, shared.AuthResponse{Error: err.Error()})
		return
	}

	// Check if user exists and password is correct
	hashedPassword, exists := users[user.Username]
	if !exists || hashedPassword != shared.HashPassword(user.Password) {
		c.JSON(http.StatusUnauthorized, shared.AuthResponse{Error: "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := shared.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.AuthResponse{Error: "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, shared.AuthResponse{
		Token:    token,
		Username: user.Username,
		Message:  "Login successful",
	})
}

func loginForm(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, shared.AuthResponse{Error: "Username and password are required"})
		return
	}

	// Check if user exists and password is correct
	hashedPassword, exists := users[username]
	if !exists || hashedPassword != shared.HashPassword(password) {
		c.JSON(http.StatusUnauthorized, shared.AuthResponse{Error: "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := shared.GenerateJWT(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.AuthResponse{Error: "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, shared.AuthResponse{
		Token:    token,
		Username: username,
		Message:  "Login successful",
	})
}