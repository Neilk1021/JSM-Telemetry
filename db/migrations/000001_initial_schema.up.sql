CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    player_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS map_attempts (
    id SERIAL PRIMARY KEY,
    player_id TEXT NOT NULL,
    map_name TEXT NOT NULL,
    passed BOOLEAN NOT NULL,
    attempts_to_complete INT NOT NULL,
    level_rounds INT NOT NULL,
    ram_used INT NOT NULL,
    cpu_used INT NOT NULL,
    duration BIGINT NOT NULL,
    first_time BOOLEAN NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS map_tower_usage (
    id SERIAL PRIMARY KEY,
    map_attempt_id INT NOT NULL REFERENCES map_attempts(id) ON DELETE CASCADE,
    tower_type TEXT NOT NULL,
    count INT NOT NULL,
    level INT NOT NULL
);
