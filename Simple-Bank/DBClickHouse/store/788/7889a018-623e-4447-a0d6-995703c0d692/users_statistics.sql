ATTACH TABLE _ UUID 'b1832a4b-6e61-4d42-8f84-fe166a20cc01'
(
    `iin` UInt64,
    `username` String,
    `action` String,
    `timestamp` DateTime
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY timestamp
SETTINGS index_granularity = 8192
