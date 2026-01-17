# Loan Validator Portal

A web application that allows users to register, login, and validate their loan status by fetching data from the Government Loan Bank API.

## Features

- User registration and authentication
- Beautiful responsive UI using Tailwind CSS and htmx
- Loan validation by fetching data from Government Loan Bank
- Session-based authentication
- Error handling for network failures and missing data

## Architecture

- **Frontend**: HTML with htmx and Tailwind CSS (no JavaScript build tools required)
- **Backend**: Go (Golang) with Gin framework
- **Database**: PostgreSQL

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Running instance of Government Loan Bank API

## Database Setup

1. Create the database:
```bash
createdb loan_validator_db
```

2. Run the schema:
```bash
psql loan_validator_db < database/schema.sql
```

Or connect to your PostgreSQL and create the database:
```sql
CREATE DATABASE loan_validator_db;
\c loan_validator_db
-- Then run the schema from database/schema.sql
```

## Environment Variables

Create a `.env` file in the `backend/` directory or set these environment variables:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=loan_validator_db

# Application Configuration
PORT=8080

# Government Loan Bank API URL
GOV_BANK_URL=http://localhost:8081
```

## Running the Application

1. Navigate to the backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

The application will start on `http://localhost:8080`

## API Endpoints

### Authentication

- `POST /auth/register` - Register a new user
  - Body: `first_name`, `last_name`, `username`, `password`
  
- `POST /auth/login` - Login user
  - Body: `username`, `password`
  
- `GET /auth/check-session` - Check if user is logged in

### Loan Validation

- `GET /api/validate-loan` - Fetch loan information from Government Loan Bank (requires authentication)

## Usage Flow

1. **Register**: Create a new account with first name, last name, username, and password
2. **Login**: Sign in with your credentials
3. **Validate Loan**: Click the "Click Here to Validate Your Loan" button
4. The system will fetch your loan information from the Government Loan Bank using your first and last name
5. If successful, you'll see:
   - First Name
   - Last Name
   - Date of Birth
   - Loan Amount Requested
   - Loan Status

## Error Messages

- **Cannot reach the government portal**: The Government Loan Bank API is not accessible
- **User details not found**: Your name is not in the Government Loan Bank database
- **Session expired**: You need to login again

## Testing

### Test User Registration
1. Open `http://localhost:8080`
2. Click "Register"
3. Fill in the form with test data (e.g., John Doe)
4. Click "Register"

### Test Loan Validation
1. Login with a user whose name exists in the Government Loan Bank
2. Click "Click Here to Validate Your Loan"
3. Check that the loan information is displayed

## Project Structure

```
Project_loan_validator_portal/
├── backend/
│   ├── main.go           # Main application code
│   └── go.mod            # Go dependencies
├── frontend/
│   └── templates/
│       └── index.html    # Frontend UI
└── database/
    └── schema.sql        # Database schema
```

## For Kubernetes Deployment

Once tested locally, this application can be containerized and deployed to Kubernetes. The following will be needed:

- Dockerfile for the Go backend
- Kubernetes manifests (Deployment, Service, ConfigMap, Secret)
- StatefulSet for PostgreSQL database
- Ingress for external access

## Notes

- Make sure the Government Loan Bank API is running before starting this application
- The application uses session-based authentication with cookies
- All passwords are hashed using bcrypt before storing
