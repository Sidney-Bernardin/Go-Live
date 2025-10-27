CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    username varchar(255) NOT NULL UNIQUE,
    email text NOT NULL UNIQUE,
    password_hash bytea NOT NULL,
    password_salt varchar(32) NOT NULL
)
