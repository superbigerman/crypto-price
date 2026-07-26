CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(20) REFERENCES currencies(symbol),
    price DECIMAL(20, 8),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_prices_symbol_created ON prices (symbol, created_at);