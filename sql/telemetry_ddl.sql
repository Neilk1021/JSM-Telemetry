CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    player_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL
);