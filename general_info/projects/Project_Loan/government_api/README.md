# Government Loan Bank API

A RESTful API service that provides customer loan information. This service acts as a data source for the Loan Validator Portal.

## Features

- RESTful API to retrieve customer loan information
- API endpoint to load/add customer data
- PostgreSQL database for storing customer records
- Query customers by first and last name
- Loan status tracking (Approved, Pending, Rejected, Under Review, Disbursed)

## Architecture

- **Backend**: Go (Golang) with Gin framework
- **Database**: PostgreSQL

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher

## Database Setup

1. Create the database:
```bash
createdb government_loan_db
```

2. Run the schema (includes sample data):
```bash
psql government_loan_db < database/schema.sql
```

Or connect to your PostgreSQL and create the database:
```sql
CREATE DATABASE government_loan_db;
\c government_loan_db
-- Then run the schema from database/schema.sql
```

The schema includes 5 sample customers for testing.

## Environment Variables

Create a `.env` file in the `backend/` directory or set these environment variables:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=government_loan_db

# Application Configuration
PORT=8081
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

The API will start on `http://localhost:8081`

## API Endpoints

### Get Customer Information

```bash
GET /api/customer?first_name=John&last_name=Doe
```

**Response** (200 OK):
```json
{
  "id": 1,
  "first_name": "John",
  "last_name": "Doe",
  "date_of_birth": "1985-05-15",
  "loan_amount_requested": 50000.00,
  "loan_status": "Approved",
  "created_at": "2024-12-04T10:00:00Z"
}
```

**Response** (404 Not Found):
```json
{
  "error": "Customer not found"
}
```

### Add Customer

```bash
POST /api/customer
Content-Type: application/json

{
  "first_name": "Alice",
  "last_name": "Johnson",
  "date_of_birth": "1995-03-20",
  "loan_amount_requested": 45000.00,
  "loan_status": "Pending"
}
```

**Response** (201 Created):
```json
{
  "message": "Customer added successfully",
  "customer_id": 6
}
```

**Valid Loan Statuses**: 
- `Approved`
- `Pending`
- `Rejected`
- `Under Review`
- `Disbursed`

### Get All Customers

```bash
GET /api/customers
```

**Response** (200 OK):
```json
{
  "count": 5,
  "customers": [
    {
      "id": 1,
      "first_name": "John",
      "last_name": "Doe",
      "date_of_birth": "1985-05-15",
      "loan_amount_requested": 50000.00,
      "loan_status": "Approved",
      "created_at": "2024-12-04T10:00:00Z"
    },
    ...
  ]
}
```

### Delete Customer

```bash
DELETE /api/customer/1
```

**Response** (200 OK):
```json
{
  "message": "Customer deleted successfully"
}
```

### Health Check

```bash
GET /health
```

**Response** (200 OK):
```json
{
  "status": "healthy",
  "service": "government_loan_bank"
}
```

## Testing the API

### Test with curl

1. **Get a customer**:
```bash
curl "http://localhost:8081/api/customer?first_name=John&last_name=Doe"
```

2. **Add a new customer**:
```bash
curl -X POST http://localhost:8081/api/customer \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Test",
    "last_name": "User",
    "date_of_birth": "1990-01-01",
    "loan_amount_requested": 25000.00,
    "loan_status": "Pending"
  }'
```

3. **Get all customers**:
```bash
curl http://localhost:8081/api/customers
```

4. **Health check**:
```bash
curl http://localhost:8081/health
```

## Sample Data

The database comes pre-loaded with 5 sample customers:

| First Name | Last Name | DOB | Loan Amount | Status |
|-----------|-----------|-----|-------------|--------|
| John | Doe | 1985-05-15 | $50,000 | Approved |
| Jane | Smith | 1990-08-22 | $75,000 | Pending |
| Michael | Johnson | 1978-03-10 | $100,000 | Disbursed |
| Sarah | Williams | 1992-11-30 | $35,000 | Under Review |
| David | Brown | 1988-07-18 | $60,000 | Rejected |

## Project Structure

```
Project_government_loan_bank/
├── backend/
│   ├── main.go           # Main application code
│   └── go.mod            # Go dependencies
└── database/
    └── schema.sql        # Database schema with sample data
```

## Integration with Loan Validator Portal

The Loan Validator Portal queries this API using:
- Endpoint: `GET /api/customer`
- Query params: `first_name` and `last_name` from the logged-in user

Make sure this service is running before starting the Loan Validator Portal.

## For Kubernetes Deployment

Once tested locally, this application can be containerized and deployed to Kubernetes. The following will be needed:

- Dockerfile for the Go backend
- Kubernetes manifests (Deployment, Service, ConfigMap, Secret)
- StatefulSet for PostgreSQL database
- Service discovery for inter-service communication

## Notes

- Customer names are case-insensitive for queries
- The API does not require authentication (can be added later for production)
- Duplicate customer names are not allowed
- All monetary values are stored with 2 decimal precision
