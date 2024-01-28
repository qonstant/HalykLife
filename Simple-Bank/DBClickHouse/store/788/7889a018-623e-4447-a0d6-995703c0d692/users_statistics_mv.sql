ATTACH MATERIALIZED VIEW _ UUID 'b919c2d9-086b-4cfd-ad57-5bb34f5d585c' TO stats.users_statistics
(
    `iin` UInt64,
    `username` String,
    `action` String,
    `timestamp` DateTime
) AS
SELECT *
FROM stats.users_statistics_queue
