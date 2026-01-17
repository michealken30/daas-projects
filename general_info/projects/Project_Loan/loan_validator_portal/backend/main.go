package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type LoanInfo struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	DateOfBirth       string `json:"date_of_birth"`
	LoanAmountRequested float64 `json:"loan_amount_requested"`
	LoanStatus        string `json:"loan_status"`
}

var db *sql.DB

func main() {
	// Database connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "loan_validator_db")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Successfully connected to database")

	// Initialize Gin router
	r := gin.Default()

	// Session management
	store := cookie.NewStore([]byte("secret-key-change-this-in-production"))
	r.Use(sessions.Sessions("loan_validator_session", store))

	// Load HTML templates
	r.LoadHTMLGlob("../frontend/templates/*")

	// Serve frontend
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", registerHandler)
		auth.POST("/login", loginHandler)
		auth.GET("/check-session", checkSessionHandler)
	}

	// API routes (protected)
	api := r.Group("/api")
	api.Use(authRequired())
	{
		api.GET("/validate-loan", validateLoanHandler)
	}

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}

func registerHandler(c *gin.Context) {
	firstName := c.PostForm("first_name")
	lastName := c.PostForm("last_name")
	username := c.PostForm("username")
	password := c.PostForm("password")

	if firstName == "" || lastName == "" || username == "" || password == "" {
		c.Data(http.StatusBadRequest, "text/html", []byte(`<div class="error">All fields are required</div>`))
		return
	}

	// Check if username already exists
	var existingUser string
	err := db.QueryRow("SELECT username FROM users WHERE username = $1", username).Scan(&existingUser)
	if err == nil {
		c.Data(http.StatusBadRequest, "text/html", []byte(`<div class="error">Username already exists</div>`))
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html", []byte(`<div class="error">Failed to hash password</div>`))
		return
	}

	// Insert user into database
	_, err = db.Exec(`INSERT INTO users (first_name, last_name, username, password) VALUES ($1, $2, $3, $4)`,
		firstName, lastName, username, string(hashedPassword))
	if err != nil {
		log.Println("Failed to insert user:", err)
		c.Data(http.StatusInternalServerError, "text/html", []byte(`<div class="error">Failed to create user</div>`))
		return
	}

	c.Data(http.StatusOK, "text/html", []byte(`<div class="success">Registration successful! Please login.</div>`))
}

func loginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.Data(http.StatusBadRequest, "text/html", []byte(`<div class="error">Username and password are required</div>`))
		return
	}

	// Get user from database
	var user User
	err := db.QueryRow(`SELECT id, first_name, last_name, username, password FROM users WHERE username = $1`, username).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Password)
	if err != nil {
		c.Data(http.StatusUnauthorized, "text/html", []byte(`<div class="error">Invalid username or password</div>`))
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.Data(http.StatusUnauthorized, "text/html", []byte(`<div class="error">Invalid username or password</div>`))
		return
	}

	// Save session
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("username", user.Username)
	session.Set("first_name", user.FirstName)
	session.Set("last_name", user.LastName)
	session.Save()

	fullName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	response := fmt.Sprintf(`<div class="success" data-username="%s">Login successful! Welcome %s</div>`, fullName, fullName)
	c.Data(http.StatusOK, "text/html", []byte(response))
}

func checkSessionHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	firstName := session.Get("first_name")
	lastName := session.Get("last_name")

	if userID == nil {
		c.JSON(http.StatusOK, gin.H{"logged_in": false})
		return
	}

	fullName := fmt.Sprintf("%s %s", firstName, lastName)
	c.JSON(http.StatusOK, gin.H{
		"logged_in": true,
		"username":  fullName,
	})
}

func validateLoanHandler(c *gin.Context) {
	session := sessions.Default(c)
	firstName := session.Get("first_name")
	lastName := session.Get("last_name")

	if firstName == nil || lastName == nil {
		c.Data(http.StatusUnauthorized, "text/html", []byte(`<div class="error">Session expired. Please login again.</div>`))
		return
	}

	// Call government loan bank API
	govBankURL := getEnv("GOV_BANK_URL", "http://localhost:8081")
	apiURL := fmt.Sprintf("%s/api/customer?first_name=%s&last_name=%s", govBankURL, firstName, lastName)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		errorMsg := `<div class="error">❌ Cannot reach the government portal to fetch details. Please try again later.</div>`
		c.Data(http.StatusServiceUnavailable, "text/html", []byte(errorMsg))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		errorMsg := `<div class="error">⚠️ User details not found in the government portal database.</div>`
		c.Data(http.StatusOK, "text/html", []byte(errorMsg))
		return
	}

	if resp.StatusCode != http.StatusOK {
		errorMsg := `<div class="error">❌ Error fetching data from government portal.</div>`
		c.Data(http.StatusServiceUnavailable, "text/html", []byte(errorMsg))
		return
	}

	// Parse response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html", []byte(`<div class="error">Failed to read response</div>`))
		return
	}

	var loanInfo LoanInfo
	err = json.Unmarshal(body, &loanInfo)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html", []byte(`<div class="error">Failed to parse response</div>`))
		return
	}

	// Return formatted loan information
	htmlResponse := fmt.Sprintf(`
		<div class="space-y-3">
			<div class="info-row">
				<span class="info-label">First Name:</span>
				<span class="info-value">%s</span>
			</div>
			<div class="info-row">
				<span class="info-label">Last Name:</span>
				<span class="info-value">%s</span>
			</div>
			<div class="info-row">
				<span class="info-label">Date of Birth:</span>
				<span class="info-value">%s</span>
			</div>
			<div class="info-row">
				<span class="info-label">Loan Amount Requested:</span>
				<span class="info-value">$%.2f</span>
			</div>
			<div class="info-row">
				<span class="info-label">Loan Status:</span>
				<span class="info-value status-%s">%s</span>
			</div>
		</div>
	`, template.HTMLEscapeString(loanInfo.FirstName),
		template.HTMLEscapeString(loanInfo.LastName),
		template.HTMLEscapeString(loanInfo.DateOfBirth),
		loanInfo.LoanAmountRequested,
		template.HTMLEscapeString(loanInfo.LoanStatus),
		template.HTMLEscapeString(loanInfo.LoanStatus))

	c.Data(http.StatusOK, "text/html", []byte(htmlResponse))
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.Data(http.StatusUnauthorized, "text/html", []byte(`<div class="error">Please login to continue</div>`))
			c.Abort()
			return
		}
		c.Next()
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
