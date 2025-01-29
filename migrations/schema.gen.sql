CREATE TABLE migrations (
	sequence           INTEGER PRIMARY KEY,
	filename           TEXT NOT NULL,
	revision           TEXT NOT NULL,
	revision_timestamp TIMESTAMP NOT NULL
);
CREATE TABLE devices(
    id                SERIAL PRIMARY KEY,
    d_name            TEXT NOT NULL,
    d_brand           TEXT NOT NULL,
    d_state           INTEGER NOT NULL,
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT (now())
 -- updated_at        TIMESTAMP NOT NULL
);