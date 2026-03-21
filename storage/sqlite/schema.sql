CREATE TABLE IF NOT EXISTS analyses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    input_path TEXT NOT NULL,
    format TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    analysis_json TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS analyses_input_path_idx
    ON analyses (input_path);
