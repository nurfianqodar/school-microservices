CREATE TYPE user_role AS ENUM (
    'teacher',
    'staff',
    'student',
    'parent'
);

CREATE TABLE users (
    -- PK
    id uuid PRIMARY KEY,
    -- Main Data
    email varchar(255) NOT NULL,
    role user_role NOT NULL,
    password_hash varchar(255) NOT NULL,
    -- Timestamp
    created_at timestamptz DEFAULT current_timestamp,
    updated_at timestamptz DEFAULT current_timestamp,
    deleted_at timestamptz,
    UNIQUE NULLS NOT DISTINCT (email, deleted_at)
);

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_deleted_at ON users (deleted_at);
