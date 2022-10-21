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

INSERT INTO api_tokens (token, is_enabled, created_at) VALUES
    ('cs3219tokena', TRUE, NOW()),
    ('cs3219tokenb', TRUE, NOW()),
    ('cs3219tokenc', TRUE, NOW()),
    ('cs3219tokend', TRUE, NOW());

INSERT INTO cards (card_unique_id, card_pokemon, card_image) VALUES
    ('xy1-1', 'Venusaur-EX', 'https://images.pokemontcg.io/xy1/1_hires.png'),
    ('xy1-2', 'Mega Venusaur-EX', 'https://images.pokemontcg.io/xy1/2_hires.png'),
    ('xy1-3', 'Weedle', 'https://images.pokemontcg.io/xy1/3_hires.png'),
    ('xy1-15', 'Scatterbug', 'https://images.pokemontcg.io/xy1/15_hires.png'),
    ('xy1-16', 'Spewpa', 'https://images.pokemontcg.io/xy1/16_hires.png'),
    ('xy1-17', 'Vivillion', 'https://images.pokemontcg.io/xy1/17_hires.png'),
    ('xy1-18', 'Skiddo', 'https://images.pokemontcg.io/xy1/18_hires.png'),
    ('xy1-19', 'Gogoat', 'https://images.pokemontcg.io/xy1/19_hires.png'),
    ('xy1-20', 'Slugma', 'https://images.pokemontcg.io/xy1/20_hires.png');
