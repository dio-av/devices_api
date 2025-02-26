CREATE TABLE migrations (
	sequence           INTEGER PRIMARY KEY,
	filename           TEXT NOT NULL,
	revision           TEXT NOT NULL,
	revision_timestamp TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS devices (
    device_id    SERIAL PRIMARY KEY,
    device_name  TEXT NOT NULL,
    device_brand TEXT NOT NULL,
    device_state INTEGER NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT (NOW())
    updated_at   TIMESTAMPTZ NOT NULL
    deleted_at   TIMESTAMPTZ NOT NULL
);