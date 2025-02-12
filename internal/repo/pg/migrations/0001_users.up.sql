CREATE TABLE users (
    username TEXT PRIMARY KEY,
    password_hash TEXT NOT NULL,
    balance INT DEFAULT 1000
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    sender_username TEXT REFERENCES users(username),
    receiver_username TEXT REFERENCES users(username),
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    username TEXT REFERENCES users(username),
    item TEXT NOT NULL,
    price INT NOT NULL,
    purchased_at TIMESTAMP DEFAULT NOW()
);
