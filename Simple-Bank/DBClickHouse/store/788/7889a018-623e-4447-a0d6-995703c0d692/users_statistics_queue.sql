ATTACH TABLE _ UUID 'db87b75e-6644-41d1-84d3-947cc102bffb'
(
    `iin` UInt64,
    `username` String,
    `action` String,
    `timestamp` DateTime
)
ENGINE = Kafka('kafka:9092', 'baeldung', 'stats-dev', 'JSONEachRow')
SETTINGS kafka_thread_per_consumer = 1, kafka_num_consumers = 1, kafka_handle_error_mode = 'stream'
