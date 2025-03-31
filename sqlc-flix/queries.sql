-- name: CreateMovie :one
INSERT INTO movies (id, title, added_at, rating)
VALUES (?, ?, ?, ?)
RETURNING id;

-- name: GetOrCreatePerson :one
INSERT INTO people (name)
VALUES (?)
ON CONFLICT(name) DO UPDATE SET name = excluded.name
RETURNING id;

-- name: GetOrCreateCountry :one
INSERT INTO countries (name)
VALUES (?)
ON CONFLICT(name) DO UPDATE SET name = excluded.name
RETURNING id;

-- name: GetOrCreateGenre :one
INSERT INTO genres (name)
VALUES (?)
ON CONFLICT(name) DO UPDATE SET name = excluded.name
RETURNING id;

-- name: AddMovieDirector :exec
INSERT OR IGNORE INTO movie_directors (movie_id, person_id)
VALUES (?, ?);

-- name: AddMovieActor :exec
INSERT OR IGNORE INTO movie_actors (movie_id, person_id)
VALUES (?, ?);

-- name: AddMovieCountry :exec
INSERT OR IGNORE INTO movie_countries (movie_id, country_id)
VALUES (?, ?);

-- name: AddMovieGenre :exec
INSERT OR IGNORE INTO movie_genres (movie_id, genre_id)
VALUES (?, ?);

-- name: GetMovie :one
SELECT
    movies.id,
    movies.title,
    movies.added_at,
    movies.rating,
    CAST(IFNULL((
        SELECT GROUP_CONCAT(name)
        FROM (
            SELECT people.name
            FROM movie_directors
            JOIN people ON people.id = movie_directors.person_id
            WHERE movie_directors.movie_id = movies.id
            ORDER BY people.name ASC
        )
    ), '') AS TEXT) AS directors,
    CAST(IFNULL((
        SELECT GROUP_CONCAT(name)
        FROM (
            SELECT people.name
            FROM movie_actors
            JOIN people ON people.id = movie_actors.person_id
            WHERE movie_actors.movie_id = movies.id
            ORDER BY people.name ASC
        )
    ), '') AS TEXT) AS actors,
    CAST(IFNULL((
        SELECT GROUP_CONCAT(name)
        FROM (
            SELECT countries.name
            FROM movie_countries
            JOIN countries ON countries.id = movie_countries.country_id
            WHERE movie_countries.movie_id = movies.id
            ORDER BY countries.name ASC
        )
    ), '') AS TEXT) AS countries,
    CAST(IFNULL((
        SELECT GROUP_CONCAT(name)
        FROM (
            SELECT genres.name
            FROM movie_genres
            JOIN genres ON genres.id = movie_genres.genre_id
            WHERE movie_genres.movie_id = movies.id
            ORDER BY genres.name
        )
    ), '') AS TEXT) AS genres
FROM movies
WHERE movies.id = ?;

-- name: QueryMovies :many
SELECT
    movies.id,
    movies.title,
    movies.added_at,
    movies.rating,
    CAST(IFNULL((
        SELECT GROUP_CONCAT(name)
        FROM (
            SELECT people.name
            FROM movie_directors
            JOIN people ON people.id = movie_directors.person_id
            WHERE movie_directors.movie_id = movies.id
            ORDER BY people.name ASC
        )
    ), '') AS TEXT) AS directors,
    CAST(IFNULL((
        SELECT GROUP_CONCAT(name)
        FROM (
            SELECT people.name
            FROM movie_actors
            JOIN people ON people.id = movie_actors.person_id
            WHERE movie_actors.movie_id = movies.id
            ORDER BY people.name ASC
        )
    ), '') AS TEXT) AS actors,
    CAST(IFNULL((
        SELECT GROUP_CONCAT(name)
        FROM (
            SELECT countries.name
            FROM movie_countries
            JOIN countries ON countries.id = movie_countries.country_id
            WHERE movie_countries.movie_id = movies.id
            ORDER BY countries.name ASC
        )
    ), '') AS TEXT) AS countries,
    CAST(IFNULL((
        SELECT GROUP_CONCAT(name)
        FROM (
            SELECT genres.name
            FROM movie_genres
            JOIN genres ON genres.id = movie_genres.genre_id
            WHERE movie_genres.movie_id = movies.id
            ORDER BY genres.name
        )
    ), '') AS TEXT) AS genres
FROM movies
WHERE
    (:search = '' OR EXISTS (
        SELECT 1
        FROM movie_directors
        JOIN people ON people.id = movie_directors.person_id
        WHERE movie_directors.movie_id = movies.id
        AND INSTR(LOWER(people.name), LOWER(:search)) > 0
    )
    OR EXISTS (
        SELECT 1
        FROM movie_actors
        JOIN people ON people.id = movie_actors.person_id
        WHERE movie_actors.movie_id = movies.id
        AND INSTR(LOWER(people.name), LOWER(:search)) > 0
    ))
    AND (:genre = '' OR EXISTS (
        SELECT 1
        FROM movie_genres
        JOIN genres ON genres.id = movie_genres.genre_id
        WHERE movie_genres.movie_id = movies.id
        AND genres.name = :genre
    ))
    AND (:country = '' OR EXISTS (
        SELECT 1
        FROM movie_countries
        JOIN countries ON countries.id = movie_countries.country_id
        WHERE movie_countries.movie_id = movies.id
        AND countries.name = :country
    ))
    AND (:added_after IS NULL OR movies.added_at >= :added_after)
    AND (:added_before IS NULL OR movies.added_at <= :added_before)
    AND (:min_rating <= 0 OR movies.rating >= :min_rating)
    AND (:max_rating <= 0 OR movies.rating <= :max_rating)
ORDER BY movies.title ASC;