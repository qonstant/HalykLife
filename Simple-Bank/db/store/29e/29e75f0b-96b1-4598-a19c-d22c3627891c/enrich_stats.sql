ATTACH TABLE _ UUID '1604dfd8-3ef5-42b7-bb2c-5dcb792aee55'
(
    `city_id` Int64,
    `country_iso_code` LowCardinality(FixedString(50)),
    `platform_id` Int32,
    `uniq` Int8,
    `listener_id` UUID,
    `podcast_id` Int64,
    `episode_id` Int64,
    `time_point` DateTime64(7, 'UTC'),
    `podcast_time_point` DateTime64(7, 'UTC'),
    `system` LowCardinality(FixedString(50)),
    `device` LowCardinality(FixedString(50)),
    `e_y` Int16,
    `e_m` Int8,
    `e_d` Int8,
    `e_h` Int8,
    `podcast_creation_time` DateTime64(7, 'UTC'),
    `podcast_rel_point` DateTime64(7, 'UTC'),
    `pr_y` Int16,
    `pr_m` Int8,
    `pr_d` Int8,
    `pr_h` Int8,
    `release_creation_time` DateTime64(7, 'UTC'),
    `release_rel_point` DateTime64(7, 'UTC'),
    `er_y` Int16,
    `er_m` Int8,
    `er_d` Int8,
    `er_h` Int8
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(time_point)
ORDER BY tuple()
SETTINGS index_granularity = 8192
