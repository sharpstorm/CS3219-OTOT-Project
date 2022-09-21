CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    user_name VARCHAR(255) UNIQUE,
    user_password VARCHAR(255)
);

CREATE TABLE wishlist (
    user_id INTEGER,
    card_id INTEGER,
    notes TEXT,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (card_id) REFERENCES cards(card_id),
    PRIMARY KEY (user_id, card_id)
);

CREATE TABLE api_tokens (
	token_id SERIAL PRIMARY KEY,
    token VARCHAR(255) UNIQUE,
    is_enabled BOOLEAN,
    created_at TIMESTAMP
);

CREATE TABLE cards (
    card_id SERIAL PRIMARY KEY,
    card_unique_id VARCHAR(255) UNIQUE,
    card_pokemon TEXT,
    card_image TEXT
);