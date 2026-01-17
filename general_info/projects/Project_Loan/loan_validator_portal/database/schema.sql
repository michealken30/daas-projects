-- Database: loan_validator_db
-- This database stores user authentication information for the loan validator portal

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on username for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Sample data (optional - uncomment to add test users)
-- Note: Password is 'password123' hashed with bcrypt
-- INSERT INTO users (first_name, last_name, username, password) 
-- VALUES 
--     ('John', 'Doe', 'johndoe', '$2a$10$YourBcryptHashHere'),
--     ('Jane', 'Smith', 'janesmith', '$2a$10$YourBcryptHashHere');
