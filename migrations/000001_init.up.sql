CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE links (
    short VARCHAR(6) PRIMARY KEY,
    original TEXT NOT NULL,
    user_id UUID NOT NULL
);

