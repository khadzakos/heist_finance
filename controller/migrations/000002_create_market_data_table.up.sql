CREATE TABLE IF NOT EXISTS market_data (
    id BIGSERIAL PRIMARY KEY,
    ticker_id BIGINT,
    price BIGINT,
    bid BIGINT,
    ask BIGINT,
    timestamp TIMESTAMP,
    FOREIGN KEY (ticker_id) REFERENCES tickers(id)
);

CREATE INDEX IF NOT EXISTS idx_ticker_timestamp ON market_data (ticker_id, timestamp);