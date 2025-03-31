CREATE TABLE movies (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    added_at DATE NOT NULL,
    rating NUMERIC NOT NULL
);

CREATE TABLE people (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE movie_directors (
    movie_id INTEGER NOT NULL REFERENCES movies (id),
    person_id INTEGER NOT NULL REFERENCES people (id),
    PRIMARY KEY (movie_id, person_id)
);

CREATE TABLE movie_actors (
    movie_id INTEGER NOT NULL REFERENCES movies (id),
    person_id INTEGER NOT NULL REFERENCES people (id),
    PRIMARY KEY (movie_id, person_id)
);

CREATE TABLE countries (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE movie_countries (
    movie_id INTEGER NOT NULL REFERENCES movies (id),
    country_id INTEGER NOT NULL REFERENCES countries (id),
    PRIMARY KEY (movie_id, country_id)
);

CREATE TABLE genres (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE movie_genres (
    movie_id INTEGER NOT NULL REFERENCES movies (id),
    genre_id INTEGER NOT NULL REFERENCES genres (id),
    PRIMARY KEY (movie_id, genre_id)
);