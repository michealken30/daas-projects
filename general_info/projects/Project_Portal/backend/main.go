package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type User struct { // I am just defining the data structure here for the users
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var (
	users     = make(map[string]string) // username -> hashed password
	jwtSecret = []byte("your-secret-key")
	uploadDir = "./uploads"
)

func main() { // Engine block of the entire flow
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(err)
	}

	r := gin.Default() // Start a gin router here 

	// Load HTML templates (check multiple possible paths)
	templatePaths := []string{
		"../frontend/templates/*",  // Local development
		"./templates/*",            // Docker container
		"templates/*",              // Alternative Docker path
	}
	
	var templatesLoaded bool
	for _, path := range templatePaths {
		if matches, _ := filepath.Glob(path); len(matches) > 0 {
			r.LoadHTMLGlob(path)
			templatesLoaded = true
			break
		}
	}
	
	if !templatesLoaded {
		panic("No HTML templates found")
	}

	// Configure CORS to allow the same origin
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8081"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Serve static files (uploaded images)
	r.Static("/uploads", uploadDir)

	// Serve the main HTML page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Authentication routes (form-based for HTMX)
	r.POST("/auth/login", loginForm)
	r.POST("/auth/register", registerForm)

	// API routes (JSON-based, keeping original for compatibility)
	r.POST("/api/register", register)
	r.POST("/api/login", login)
	r.POST("/api/upload", authMiddleware(), upload)
	r.GET("/api/profile", authMiddleware(), getProfile)

	r.Run(":8081")
}

func register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	if _, exists := users[user.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword := hashPassword(user.Password)
	users[user.Username] = hashedPassword

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func registerForm(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Check if user already exists
	if _, exists := users[username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword := hashPassword(password)
	users[username] = hashedPassword

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists and password is correct
	hashedPassword, exists := users[user.Username]
	if !exists || hashedPassword != hashPassword(user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"username": user.Username,
		"message":  "Login successful",
	})
}

func loginForm(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Check if user exists and password is correct
	hashedPassword, exists := users[username]
	if !exists || hashedPassword != hashPassword(password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := generateJWT(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"username": username,
		"message":  "Login successful",
	})
}

func upload(c *gin.Context) {
	username := c.GetString("username")

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Create unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d_%s", username, timestamp, header.Filename)
	filepath := filepath.Join(uploadDir, filename)

	// Save file
	out, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
		return
	}

	imageURL := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"imageUrl": imageURL,
		"filename": filename,
	})
}

func getProfile(c *gin.Context) {
	username := c.GetString("username")
	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"message":  fmt.Sprintf("Welcome %s!", username),
	})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(*Claims)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func generateJWT(username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}