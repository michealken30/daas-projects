package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func getAuthServiceURL() string {
	if url := os.Getenv("AUTH_SERVICE_URL"); url != "" {
		return url
	}
	return "http://localhost:8082" // fallback for local development
}

func getUploadServiceURL() string {
	if url := os.Getenv("UPLOAD_SERVICE_URL"); url != "" {
		return url
	}
	return "http://localhost:8083" // fallback for local development
}

func main() {
	r := gin.Default()

	// Load HTML templates (check multiple possible paths)
	templatePaths := []string{
		"../../frontend/templates/*",     // Local dev (from services/api-gateway - used by start script)
		"frontend/templates/*",           // Local development (from project root)
		"./templates/*",                  // Docker container
		"templates/*",                    // Alternative Docker path
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

	// Serve the main HTML page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Proxy routes to auth service
	r.POST("/auth/login", proxyToAuth("/login-form"))
	r.POST("/auth/register", proxyToAuth("/register-form"))
	r.POST("/api/register", proxyToAuth("/register"))
	r.POST("/api/login", proxyToAuth("/login"))

	// Proxy routes to upload service with file handling
	r.POST("/api/upload", proxyFileToUpload("/upload"))
	r.GET("/api/profile", proxyToUpload("/profile"))

	// Serve static files from upload service
	r.GET("/uploads/*filepath", func(c *gin.Context) {
		// Proxy to upload service for static files
		proxyStaticToUpload(c)
	})

	r.Run(":8081") // API Gateway on port 8081 (original port)
}

func proxyToAuth(endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := getAuthServiceURL() + endpoint

		// Handle form data or JSON
		var body io.Reader
		contentType := c.GetHeader("Content-Type")

		if contentType == "application/json" {
			// Forward JSON body
			jsonData, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
				return
			}
			body = bytes.NewReader(jsonData)
		} else {
			// Handle form data
			if err := c.Request.ParseForm(); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
				return
			}
			body = bytes.NewBufferString(c.Request.Form.Encode())
			contentType = "application/x-www-form-urlencoded"
		}

		// Create request to auth service
		req, err := http.NewRequest(c.Request.Method, url, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		req.Header.Set("Content-Type", contentType)

		// Forward the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth service unavailable"})
			return
		}
		defer resp.Body.Close()

		// Forward response
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		c.Data(resp.StatusCode, "application/json", respBody)
	}
}

func proxyToUpload(endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := getUploadServiceURL() + endpoint

		// Forward authorization header
		auth := c.GetHeader("Authorization")

		req, err := http.NewRequest(c.Request.Method, url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		if auth != "" {
			req.Header.Set("Authorization", auth)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload service unavailable"})
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		c.Data(resp.StatusCode, "application/json", respBody)
	}
}

func proxyFileToUpload(endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := getUploadServiceURL() + endpoint

		// Get the multipart form
		err := c.Request.ParseMultipartForm(32 << 20) // 32 MB
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
			return
		}

		// Create a new multipart form
		body := &bytes.Buffer{}
		writer := &bytes.Buffer{}

		// Forward the file
		file, header, err := c.Request.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}
		defer file.Close()

		// Read file content
		fileContent, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}

		// Create multipart form for forwarding
		boundary := "----WebKitFormBoundary7MA4YWxkTrZu0gW"
		contentType := fmt.Sprintf("multipart/form-data; boundary=%s", boundary)

		writer.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		writer.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"image\"; filename=\"%s\"\r\n", header.Filename))
		writer.WriteString(fmt.Sprintf("Content-Type: %s\r\n\r\n", header.Header.Get("Content-Type")))
		writer.Write(fileContent)
		writer.WriteString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

		body = writer

		// Forward authorization header
		auth := c.GetHeader("Authorization")

		req, err := http.NewRequest("POST", url, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		req.Header.Set("Content-Type", contentType)
		if auth != "" {
			req.Header.Set("Authorization", auth)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload service unavailable"})
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		c.Data(resp.StatusCode, "application/json", respBody)
	}
}

func proxyStaticToUpload(c *gin.Context) {
	filepath := c.Param("filepath")
	url := getUploadServiceURL() + "/uploads" + filepath

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer resp.Body.Close()

	// Forward the response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}