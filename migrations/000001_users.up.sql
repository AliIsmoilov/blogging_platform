CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    full_name VARCHAR,
    email VARCHAR,
    password VARCHAR,
    phone_number VARCHAR,
    balance decimal(18, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
