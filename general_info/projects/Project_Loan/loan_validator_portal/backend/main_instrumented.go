package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	FirstName           string  `json:"first_name"`
	LastName            string  `json:"last_name"`
	DateOfBirth         string  `json:"date_of_birth"`
	LoanAmountRequested float64 `json:"loan_amount_requested"`
	LoanStatus          string  `json:"loan_status"`
}

var db *sql.DB
var tracer trace.Tracer

// Initialize OpenTelemetry
func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	otlpEndpoint := getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	
	log.Info().Str("endpoint", otlpEndpoint).Msg("Initializing OpenTelemetry tracer")

	// Create OTLP trace exporter
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(otlpEndpoint),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("loan-validator-portal"),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("environment", getEnv("ENV", "development")),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracer = tp.Tracer("loan-validator-portal")
	
	log.Info().Msg("OpenTelemetry tracer initialized successfully")
	return tp, nil
}

func main() {
	// Initialize structured logging with console writer for pretty output
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(consoleWriter)
	
	ctx := context.Background()

	// Initialize OpenTelemetry
	tp, err := initTracer(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize tracer")
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Error shutting down tracer provider")
		}
	}()

	// Database connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "loan_validator_db")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Create span for database connection
	ctx, span := tracer.Start(ctx, "database.connect")
	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.name", dbName),
		attribute.String("db.host", dbHost),
		attribute.String("db.port", dbPort),
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to open database connection")
		span.End()
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Test database connection with ping
	err = db.Ping()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to ping database")
		span.End()
		log.Fatal().Err(err).Str("host", dbHost).Str("port", dbPort).Msg("Failed to ping database - connection broken")
	}
	
	span.SetStatus(codes.Ok, "Database connected successfully")
	span.End()
	log.Info().Str("host", dbHost).Str("port", dbPort).Msg("Successfully connected to database")

	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// Add OpenTelemetry middleware
	r.Use(otelgin.Middleware("loan-validator-portal"))

	// Add trace ID to logs middleware
	r.Use(func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		if span.SpanContext().IsValid() {
			traceID := span.SpanContext().TraceID().String()
			c.Set("trace_id", traceID)
			log.Logger = log.With().Str("trace_id", traceID).Logger()
		}
		c.Next()
	})

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
	log.Info().Str("port", port).Msg("Server starting")
	
	if err := r.Run(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

func registerHandler(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "auth.register")
	defer span.End()

	firstName := c.PostForm("first_name")
	lastName := c.PostForm("last_name")
	username := c.PostForm("username")
	password := c.PostForm("password")

	span.SetAttributes(
		attribute.String("user.username", username),
		attribute.String("user.first_name", firstName),
		attribute.String("user.last_name", lastName),
	)

	log.Info().Str("username", username).Msg("User registration attempt")

	if firstName == "" || lastName == "" || username == "" || password == "" {
		span.SetStatus(codes.Error, "Missing required fields")
		log.Warn().Str("username", username).Msg("Registration failed: missing fields")
		c.Data(http.StatusBadRequest, "text/html", []byte(`<div class="error">All fields are required</div>`))
		return
	}

	// Check if username already exists
	var existingUser string
	err := db.QueryRowContext(ctx, "SELECT username FROM users WHERE username = $1", username).Scan(&existingUser)
	if err == nil {
		span.SetStatus(codes.Error, "Username already exists")
		log.Warn().Str("username", username).Msg("Registration failed: username exists")
		c.Data(http.StatusBadRequest, "text/html", []byte(`<div class="error">Username already exists</div>`))
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to hash password")
		log.Error().Err(err).Msg("Failed to hash password")
		c.Data(http.StatusInternalServerError, "text/html", []byte(`<div class="error">Failed to hash password</div>`))
		return
	}

	// Insert user into database
	_, err = db.ExecContext(ctx, `INSERT INTO users (first_name, last_name, username, password) VALUES ($1, $2, $3, $4)`,
		firstName, lastName, username, string(hashedPassword))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Database insert failed")
		log.Error().Err(err).Str("username", username).Msg("Failed to insert user")
		c.Data(http.StatusInternalServerError, "text/html", []byte(`<div class="error">Failed to create user</div>`))
		return
	}

	span.SetStatus(codes.Ok, "User registered successfully")
	log.Info().Str("username", username).Msg("User registered successfully")
	c.Data(http.StatusOK, "text/html", []byte(`<div class="success">Registration successful! Please login.</div>`))
}

func loginHandler(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "auth.login")
	defer span.End()

	username := c.PostForm("username")
	password := c.PostForm("password")

	span.SetAttributes(attribute.String("user.username", username))
	log.Info().Str("username", username).Msg("Login attempt")

	if username == "" || password == "" {
		span.SetStatus(codes.Error, "Missing credentials")
		log.Warn().Str("username", username).Msg("Login failed: missing credentials")
		c.Data(http.StatusBadRequest, "text/html", []byte(`<div class="error">Username and password are required</div>`))
		return
	}

	// Get user from database
	var user User
	err := db.QueryRowContext(ctx, `SELECT id, first_name, last_name, username, password FROM users WHERE username = $1`, username).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Password)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		log.Warn().Err(err).Str("username", username).Msg("Login failed: user not found")
		c.Data(http.StatusUnauthorized, "text/html", []byte(`<div class="error">Invalid username or password</div>`))
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		span.SetStatus(codes.Error, "Invalid password")
		log.Warn().Str("username", username).Msg("Login failed: invalid password")
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

	span.SetAttributes(attribute.Int("user.id", user.ID))
	span.SetStatus(codes.Ok, "Login successful")
	log.Info().Str("username", username).Int("user_id", user.ID).Msg("User logged in successfully")

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
	ctx, span := tracer.Start(c.Request.Context(), "api.validate_loan")
	defer span.End()

	session := sessions.Default(c)
	firstName := session.Get("first_name")
	lastName := session.Get("last_name")

	span.SetAttributes(
		attribute.String("user.first_name", fmt.Sprint(firstName)),
		attribute.String("user.last_name", fmt.Sprint(lastName)),
	)

	log.Info().
		Str("first_name", fmt.Sprint(firstName)).
		Str("last_name", fmt.Sprint(lastName)).
		Msg("Loan validation request")

	if firstName == nil || lastName == nil {
		span.SetStatus(codes.Error, "Session expired")
		log.Warn().Msg("Loan validation failed: session expired")
		c.Data(http.StatusUnauthorized, "text/html", []byte(`<div class="error">Session expired. Please login again.</div>`))
		return
	}

	// Call government loan bank API
	govBankURL := getEnv("GOV_BANK_URL", "http://localhost:8081")
	apiURL := fmt.Sprintf("%s/api/customer?first_name=%s&last_name=%s", govBankURL, firstName, lastName)

	// Create child span for external HTTP call
	ctx, httpSpan := tracer.Start(ctx, "http.client.government_bank")
	httpSpan.SetAttributes(
		attribute.String("http.url", apiURL),
		attribute.String("http.method", "GET"),
		attribute.String("peer.service", "government-loan-bank"),
	)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		httpSpan.RecordError(err)
		httpSpan.SetStatus(codes.Error, "Failed to create request")
		httpSpan.End()
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create request")
		log.Error().Err(err).Msg("Failed to create request to government portal")
		c.Data(http.StatusInternalServerError, "text/html", []byte(`<div class="error">Failed to create request</div>`))
		return
	}

	// Inject trace context into HTTP headers
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := client.Do(req)
	if err != nil {
		httpSpan.RecordError(err)
		httpSpan.SetStatus(codes.Error, "Connection failed")
		httpSpan.End()
		span.RecordError(err)
		span.SetStatus(codes.Error, "Cannot reach government portal")
		log.Error().Err(err).Str("url", govBankURL).Msg("Cannot reach government portal")
		errorMsg := `<div class="error">❌ Cannot reach the government portal to fetch details. Please try again later.</div>`
		c.Data(http.StatusServiceUnavailable, "text/html", []byte(errorMsg))
		return
	}
	defer resp.Body.Close()

	httpSpan.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
	httpSpan.End()

	if resp.StatusCode == http.StatusNotFound {
		span.SetStatus(codes.Ok, "User not found in government database")
		log.Warn().
			Str("first_name", fmt.Sprint(firstName)).
			Str("last_name", fmt.Sprint(lastName)).
			Msg("User not found in government portal")
		errorMsg := `<div class="error">⚠️ User details not found in the government portal database.</div>`
		c.Data(http.StatusOK, "text/html", []byte(errorMsg))
		return
	}

	if resp.StatusCode != http.StatusOK {
		span.SetStatus(codes.Error, fmt.Sprintf("Unexpected status code: %d", resp.StatusCode))
		log.Error().Int("status_code", resp.StatusCode).Msg("Error response from government portal")
		errorMsg := `<div class="error">❌ Error fetching data from government portal.</div>`
		c.Data(http.StatusServiceUnavailable, "text/html", []byte(errorMsg))
		return
	}

	// Parse response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to read response")
		log.Error().Err(err).Msg("Failed to read response from government portal")
		c.Data(http.StatusInternalServerError, "text/html", []byte(`<div class="error">Failed to read response</div>`))
		return
	}

	var loanInfo LoanInfo
	err = json.Unmarshal(body, &loanInfo)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to parse response")
		log.Error().Err(err).Msg("Failed to parse response from government portal")
		c.Data(http.StatusInternalServerError, "text/html", []byte(`<div class="error">Failed to parse response</div>`))
		return
	}

	span.SetAttributes(
		attribute.String("loan.status", loanInfo.LoanStatus),
		attribute.Float64("loan.amount", loanInfo.LoanAmountRequested),
	)
	span.SetStatus(codes.Ok, "Loan information retrieved successfully")
	
	log.Info().
		Str("first_name", loanInfo.FirstName).
		Str("last_name", loanInfo.LastName).
		Str("loan_status", loanInfo.LoanStatus).
		Float64("loan_amount", loanInfo.LoanAmountRequested).
		Msg("Loan validation successful")

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
