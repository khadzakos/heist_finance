CREATE TABLE IF NOT EXISTS tickers (
    id BIGSERIAL PRIMARY KEY,
    exchange VARCHAR(50),
    symbol VARCHAR(50),
    UNIQUE (exchange, symbol)
);