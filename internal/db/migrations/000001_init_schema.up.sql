CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    username VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    bank BIGINT,
    awards INT
);

CREATE TABLE IF NOT EXISTS heros (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(255),
    description VARCHAR(255),
    rarity VARCHAR(255),
    damage_type VARCHAR(255),
    effect VARCHAR(255),
    hitpoint INT,
    damage INT,
    cost_elixir INT,
    damage_tower INT,
    speed INT,
    price INT
);

CREATE TABLE IF NOT EXISTS spells (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(255),
    description VARCHAR(255),
    area INT,
    damage_type VARCHAR(255),
    damage INT,
    duration BIGINT,
    effect VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS decks (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id INT REFERENCES users(id),
    name VARCHAR(255),
    description VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS deck_heros (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    deck_id INT REFERENCES decks(id),
    hero_id INT REFERENCES heros(id)
);

CREATE TABLE IF NOT EXISTS deck_spells (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    deck_id INT REFERENCES decks(id),
    spell_id INT REFERENCES spells(id)
);

CREATE TABLE IF NOT EXISTS user_heros (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id INT REFERENCES users(id),
    hero_id INT REFERENCES heros(id)
);

CREATE TABLE IF NOT EXISTS user_spells (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id INT REFERENCES users(id),
    spell_id INT REFERENCES spells(id)
);
