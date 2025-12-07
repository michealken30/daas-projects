-- Database: government_loan_db
-- This database stores customer loan information for the government loan bank

-- Create customers table
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth VARCHAR(20) NOT NULL,
    loan_amount_requested DECIMAL(12, 2) NOT NULL,
    loan_status VARCHAR(50) NOT NULL CHECK (loan_status IN ('Approved', 'Pending', 'Rejected', 'Under Review', 'Disbursed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on names for faster lookups
CREATE INDEX IF NOT EXISTS idx_customers_names ON customers(LOWER(first_name), LOWER(last_name));

-- Sample data for testing
INSERT INTO customers (first_name, last_name, date_of_birth, loan_amount_requested, loan_status) 
VALUES 
    ('John', 'Doe', '1985-05-15', 50000.00, 'Approved'),
    ('Jane', 'Smith', '1990-08-22', 75000.00, 'Pending'),
    ('Michael', 'Johnson', '1978-03-10', 100000.00, 'Disbursed'),
    ('Sarah', 'Williams', '1992-11-30', 35000.00, 'Under Review'),
    ('David', 'Brown', '1988-07-18', 60000.00, 'Rejected')
ON CONFLICT DO NOTHING;
