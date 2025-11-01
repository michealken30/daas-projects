package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"shared"
)

var uploadDir = "./uploads"

func main() {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(err)
	}

	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Serve static files (uploaded images)
	r.Static("/uploads", uploadDir)

	// Upload routes
	r.POST("/upload", authMiddleware(), upload)
	r.GET("/profile", authMiddleware(), getProfile)

	r.Run(":8083") // Upload service on port 8083
}

func upload(c *gin.Context) {
	username := c.GetString("username")

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.UploadResponse{Error: "No file uploaded"})
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
		c.JSON(http.StatusInternalServerError, shared.UploadResponse{Error: "Could not save file"})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.UploadResponse{Error: "Could not save file"})
		return
	}

	imageURL := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, shared.UploadResponse{
		Message:  "File uploaded successfully",
		ImageURL: imageURL,
		Filename: filename,
	})
}

func getProfile(c *gin.Context) {
	username := c.GetString("username")
	c.JSON(http.StatusOK, shared.ProfileResponse{
		Username: username,
		Message:  fmt.Sprintf("Welcome %s!", username),
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
		if len(tokenString) > 7 && strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = tokenString[7:]
		}

		claims, err := shared.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}