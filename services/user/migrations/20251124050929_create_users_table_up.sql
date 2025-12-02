-- Migration up
CREATE TYPE user_role AS ENUM ('user', 'admin');

CREATE TABLE users (
    id BIGINT PRIMARY KEY,
--- local@domain => 64 (local) + 1 (@) + 253 (domain) = 318 characters
    email VARCHAR(318) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'user'
);
