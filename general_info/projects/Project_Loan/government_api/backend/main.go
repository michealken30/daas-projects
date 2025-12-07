package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Customer struct {
	ID                  int       `json:"id"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	DateOfBirth         string    `json:"date_of_birth"`
	LoanAmountRequested float64   `json:"loan_amount_requested"`
	LoanStatus          string    `json:"loan_status"`
	CreatedAt           time.Time `json:"created_at"`
}

type CustomerInput struct {
	FirstName           string  `json:"first_name" binding:"required"`
	LastName            string  `json:"last_name" binding:"required"`
	DateOfBirth         string  `json:"date_of_birth" binding:"required"`
	LoanAmountRequested float64 `json:"loan_amount_requested" binding:"required"`
	LoanStatus          string  `json:"loan_status" binding:"required"`
}

var db *sql.DB

func main() {
	// Database connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "government_loan_db")

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

	// API routes
	api := r.Group("/api")
	{
		// Get customer by first and last name
		api.GET("/customer", getCustomerHandler)

		// Add new customer (for loading data)
		api.POST("/customer", addCustomerHandler)

		// Get all customers (for testing/admin)
		api.GET("/customers", getAllCustomersHandler)

		// Delete customer (for testing)
		api.DELETE("/customer/:id", deleteCustomerHandler)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "government_loan_bank",
		})
	})

	port := getEnv("PORT", "8081")
	log.Printf("Government Loan Bank API starting on port %s", port)
	r.Run(":" + port)
}

func getCustomerHandler(c *gin.Context) {
	firstName := c.Query("first_name")
	lastName := c.Query("last_name")

	if firstName == "" || lastName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "first_name and last_name query parameters are required",
		})
		return
	}

	var customer Customer
	query := `SELECT id, first_name, last_name, date_of_birth, loan_amount_requested, loan_status, created_at 
	          FROM customers 
	          WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2)`

	err := db.QueryRow(query, firstName, lastName).Scan(
		&customer.ID,
		&customer.FirstName,
		&customer.LastName,
		&customer.DateOfBirth,
		&customer.LoanAmountRequested,
		&customer.LoanStatus,
		&customer.CreatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Customer not found",
		})
		return
	}

	if err != nil {
		log.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve customer information",
		})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func addCustomerHandler(c *gin.Context) {
	var input CustomerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input: " + err.Error(),
		})
		return
	}

	// Validate loan status
	validStatuses := map[string]bool{
		"Approved":       true,
		"Pending":        true,
		"Rejected":       true,
		"Under Review":   true,
		"Disbursed":      true,
	}

	if !validStatuses[input.LoanStatus] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid loan status. Must be one of: Approved, Pending, Rejected, Under Review, Disbursed",
		})
		return
	}

	// Check if customer already exists
	var existingID int
	checkQuery := `SELECT id FROM customers WHERE LOWER(first_name) = LOWER($1) AND LOWER(last_name) = LOWER($2)`
	err := db.QueryRow(checkQuery, input.FirstName, input.LastName).Scan(&existingID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Customer already exists with this name",
			"existing_id": existingID,
		})
		return
	}

	// Insert customer
	var customerID int
	insertQuery := `INSERT INTO customers (first_name, last_name, date_of_birth, loan_amount_requested, loan_status) 
	                VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err = db.QueryRow(insertQuery,
		input.FirstName,
		input.LastName,
		input.DateOfBirth,
		input.LoanAmountRequested,
		input.LoanStatus,
	).Scan(&customerID)

	if err != nil {
		log.Println("Failed to insert customer:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add customer",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer added successfully",
		"customer_id": customerID,
	})
}

func getAllCustomersHandler(c *gin.Context) {
	rows, err := db.Query(`SELECT id, first_name, last_name, date_of_birth, loan_amount_requested, loan_status, created_at 
	                        FROM customers ORDER BY created_at DESC`)
	if err != nil {
		log.Println("Failed to query customers:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve customers",
		})
		return
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var customer Customer
		err := rows.Scan(
			&customer.ID,
			&customer.FirstName,
			&customer.LastName,
			&customer.DateOfBirth,
			&customer.LoanAmountRequested,
			&customer.LoanStatus,
			&customer.CreatedAt,
		)
		if err != nil {
			log.Println("Failed to scan customer:", err)
			continue
		}
		customers = append(customers, customer)
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(customers),
		"customers": customers,
	})
}

func deleteCustomerHandler(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM customers WHERE id = $1", id)
	if err != nil {
		log.Println("Failed to delete customer:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete customer",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Customer not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Customer deleted successfully",
	})
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
