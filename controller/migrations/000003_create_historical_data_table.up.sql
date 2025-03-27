CREATE TABLE IF NOT EXISTS historical_data (
    id BIGSERIAL PRIMARY KEY,
    ticker_id BIGINT,
    open BIGINT,
    high BIGINT,
    low BIGINT,
    close BIGINT,
    volume BIGINT,
    timestamp TIMESTAMP,
    UNIQUE (ticker_id, timestamp),
    FOREIGN KEY (ticker_id) REFERENCES tickers(id)
);

CREATE INDEX IF NOT EXISTS idx_ticker_timestamp ON historical_data (ticker_id, timestamp);