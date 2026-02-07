-- name: CreateTicker :one
INSERT INTO inventory_symbol (symbol)
VALUES (?)
RETURNING *;

-- name: CreateOrGetTicker :one
INSERT INTO inventory_symbol (symbol)
VALUES (?)
ON CONFLICT(symbol) DO UPDATE SET
    symbol = excluded.symbol
RETURNING symbolid, symbol;

-- name: InsertDailyPrice :one
INSERT INTO prices (symbolid, trading_day, open, close)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: PctTopN :many
SELECT
    p.symbolid,
    s.symbol,
    p.trading_day,
    p.pct_delta
FROM prices p
JOIN inventory_symbol s
  ON p.symbolid = s.symbolid
WHERE p.trading_day = ?
ORDER BY p.pct_delta DESC
LIMIT ?
;


-- name: InsertWeeklyPrices :exec
INSERT INTO weekly_prices (symbolid, week_start, open, close)
SELECT DISTINCT
    symbolid,
    week_start,
    FIRST_VALUE(open) OVER w AS open,
    LAST_VALUE(close) OVER w AS close
FROM (
    SELECT
        symbolid,
        trading_day,
        open,
        close,
        date(trading_day, 'weekday 1', '-7 days') AS week_start
    FROM prices
    WHERE trading_day >= date(?, 'weekday 1', '-7 days')
      AND trading_day < date(?, 'weekday 1')
)
WINDOW w AS (
    PARTITION BY symbolid, week_start
    ORDER BY trading_day
    ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING
);

-- name: DeleteWeeklyPrices :exec
DELETE FROM weekly_prices
WHERE week_start = ? 
