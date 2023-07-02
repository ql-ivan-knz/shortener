CREATE TABLE IF NOT EXISTS links (
    original_url text UNIQUE NOT NULL,
    hash_url varchar(8) NOT NULL
);

CREATE INDEX hash_idx ON links (hash_url);