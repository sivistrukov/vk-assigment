CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    is_admin bool DEFAULT false NOT NULL
);
CREATE TABLE IF NOT EXISTS actors (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    middle_name VARCHAR NULL,
    sex VARCHAR NOT NULL,
    birthday TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS films (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    release_date TIMESTAMP NOT NULL,
    rating INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS actors_and_films (
    id SERIAL PRIMARY KEY,
    actor_id INTEGER REFERENCES actors (id) ON DELETE CASCADE NOT NULL,
    film_id INTEGER REFERENCES films (id) ON DELETE CASCADE NOT NULL,
    CONSTRAINT uniq_actor_films UNIQUE (actor_id, film_id)
);