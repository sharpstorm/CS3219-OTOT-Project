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
