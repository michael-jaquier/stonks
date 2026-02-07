CREATE TABLE inventory_symbol (
    symbolid INTEGER PRIMARY KEY,
    symbol   TEXT NOT NULL UNIQUE
);

CREATE TABLE prices (
    symbolid INTEGER NOT NULL,
    trading_day DATE NOT NULL,

    open   REAL NOT NULL,
    close  REAL NOT NULL,

    delta  REAL GENERATED ALWAYS AS (close - open) STORED,
    pct_delta REAL GENERATED ALWAYS AS ((close - open) / open) STORED,

    PRIMARY KEY (symbolid, trading_day),

    FOREIGN KEY (symbolid)
        REFERENCES inventory_symbol(symbolid)
        ON DELETE CASCADE
);


CREATE TABLE weekly_prices (
    symbolid INTEGER NOT NULL,
    week_start DATE NOT NULL, -- Monday

    open   REAL NOT NULL,
    close  REAL NOT NULL,

    delta  REAL GENERATED ALWAYS AS (close - open) STORED,
    pct_delta REAL GENERATED ALWAYS AS ((close - open) / open) STORED,

    PRIMARY KEY (symbolid, week_start),

    FOREIGN KEY (symbolid)
        REFERENCES inventory_symbol(symbolid)
        ON DELETE CASCADE
);

CREATE INDEX idx_prices_pctdelta
ON prices(trading_day, pct_delta DESC);
CREATE INDEX idx_prices_delta
ON prices(trading_day, delta DESC);
CREATE INDEX idx_prices_delta_weekly
ON weekly_prices(week_start, pct_delta DESC);
