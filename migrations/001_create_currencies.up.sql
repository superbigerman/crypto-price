CREATE TABLE IF NOT EXISTS currencies (
    symbol VARCHAR(20) PRIMARY KEY,
    name VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);