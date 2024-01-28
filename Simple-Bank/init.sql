-- Создание базы данных
CREATE DATABASE IF NOT EXISTS stats;

-- Создание таблицы для очереди данных из Kafka
CREATE TABLE stats.users_statistics_queue (
                                              iin UInt64,
                                              username varchar,
                                              action varchar,
                                              timestamp DateTime
) ENGINE = Kafka('kafka:9092', 'baeldung', 'stats-dev','JSONEachRow') SETTINGS kafka_thread_per_consumer = 1, kafka_num_consumers = 1, kafka_handle_error_mode = 'stream';

-- Создание таблицы для хранения статистики
CREATE TABLE IF NOT EXISTS stats.users_statistics (
                                                      iin UInt64,
                                                      username varchar,
                                                      action varchar,
                                                      timestamp DateTime
) ENGINE = MergeTree()
    PARTITION BY toYYYYMM(timestamp)
    ORDER BY (timestamp);

-- Создание материализованного представления для автоматической вставки данных в основную таблицу
CREATE MATERIALIZED VIEW IF NOT EXISTS stats.users_statistics_mv TO stats.users_statistics AS
SELECT * FROM stats.users_statistics_queue;
