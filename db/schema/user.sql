CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    first_name VARCHAR(300) NOT NULL,
    last_name VARCHAR(300) NOT NULL DEFAULT '',
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hsh TEXT NOT NULL,
    salt Text NOT NULL,
    avatar_url TEXT,
    bio TEXT,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ DEFAULT NULL
);
