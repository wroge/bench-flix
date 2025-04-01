package sqlflix

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	benchflix "github.com/wroge/bench-flix"
)

func NewRepository() benchflix.Repository {
	db, err := sql.Open("sqlite3", ":memory:?_fk=1")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(
		`CREATE TABLE movies (
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
			movie_id INTEGER REFERENCES movies (id),
			person_id INTEGER REFERENCES people (id),
			PRIMARY KEY (movie_id, person_id)
		);

		CREATE TABLE movie_actors (
			movie_id INTEGER REFERENCES movies (id),
			person_id INTEGER REFERENCES people (id),
			PRIMARY KEY (movie_id, person_id)
		);

		CREATE TABLE countries (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);

		CREATE TABLE movie_countries (
			movie_id INTEGER REFERENCES movies (id),
			country_id INTEGER REFERENCES countries (id),
			PRIMARY KEY (movie_id, country_id)
		);

		CREATE TABLE genres (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);

		CREATE TABLE movie_genres (
			movie_id INTEGER REFERENCES movies (id),
			genre_id INTEGER REFERENCES genres (id),
			PRIMARY KEY (movie_id, genre_id)
		);`)
	if err != nil {
		panic(err)
	}

	return Repository{
		DB: db,
	}
}

type Repository struct {
	DB *sql.DB
}

// Create implements benchflix.Repository.
func (r Repository) Create(ctx context.Context, movie benchflix.Movie) (err error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO movies (id, title, added_at, rating) VALUES (?, ?, ?, ?);`,
		movie.ID, movie.Title, movie.AddedAt, movie.Rating,
	)
	if err != nil {
		return err
	}

	if len(movie.Directors)+len(movie.Actors) > 0 {
		personArgs := make([]any, len(movie.Directors)+len(movie.Actors))

		for i, p := range append(movie.Directors, movie.Actors...) {
			personArgs[i] = p
		}

		rows, err := tx.QueryContext(ctx,
			fmt.Sprintf(
				`INSERT INTO people (name) VALUES %s ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`,
				strings.Repeat(",(?)", len(personArgs))[1:],
			),
			personArgs...)
		if err != nil {
			return err
		}

		defer func() {
			err = errors.Join(err, rows.Err(), rows.Close())
		}()

		personIDs := make([]int64, len(personArgs))
		index := 0

		for rows.Next() {
			if err = rows.Scan(&personIDs[index]); err != nil {
				return err
			}

			index++
		}

		if len(movie.Directors) > 0 {
			movieDirectorArgs := make([]any, len(movie.Directors)*2)

			for i := range len(movie.Directors) {
				movieDirectorArgs[i*2] = movie.ID
				movieDirectorArgs[i*2+1] = personIDs[i]
			}

			_, err = tx.ExecContext(ctx,
				fmt.Sprintf(
					`INSERT INTO movie_directors (movie_id, person_id) VALUES %s;`,
					strings.Repeat(",(?, ?)", len(movie.Directors))[1:],
				),
				movieDirectorArgs...,
			)
			if err != nil {
				return err
			}
		}

		if len(movie.Actors) > 0 {
			movieActorArgs := make([]any, len(movie.Actors)*2)

			for i := range len(movie.Actors) {
				movieActorArgs[i*2] = movie.ID
				movieActorArgs[i*2+1] = personIDs[len(movie.Directors)+i]
			}

			_, err = tx.ExecContext(ctx,
				fmt.Sprintf(
					`INSERT INTO movie_actors (movie_id, person_id) VALUES %s;`,
					strings.Repeat(",(?, ?)", len(movie.Actors))[1:],
				),
				movieActorArgs...,
			)
			if err != nil {
				return err
			}
		}
	}

	if len(movie.Countries) > 0 {
		countryArgs := make([]any, len(movie.Countries))

		for i, c := range movie.Countries {
			countryArgs[i] = c
		}

		rows, err := tx.QueryContext(ctx,
			fmt.Sprintf(
				`INSERT INTO countries (name) VALUES %s ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id;`,
				strings.Repeat(",(?)", len(countryArgs))[1:],
			),
			countryArgs...,
		)
		if err != nil {
			return err
		}

		defer func() {
			err = errors.Join(err, rows.Err(), rows.Close())
		}()

		countryIDs := make([]int64, len(countryArgs))
		index := 0

		for rows.Next() {
			if err = rows.Scan(&countryIDs[index]); err != nil {
				return err
			}

			index++
		}

		movieCountryArgs := make([]any, len(movie.Countries)*2)

		for i := range len(movie.Countries) {
			movieCountryArgs[i*2] = movie.ID
			movieCountryArgs[i*2+1] = countryIDs[i]
		}

		_, err = tx.ExecContext(ctx,
			fmt.Sprintf(
				`INSERT INTO movie_countries (movie_id, country_id) VALUES %s;`,
				strings.Repeat(",(?, ?)", len(movie.Countries))[1:],
			),
			movieCountryArgs...,
		)
		if err != nil {
			return err
		}
	}

	if len(movie.Genres) > 0 {
		genreArgs := make([]any, len(movie.Genres))

		for i, c := range movie.Genres {
			genreArgs[i] = c
		}

		rows, err := tx.QueryContext(ctx,
			fmt.Sprintf(
				`INSERT INTO genres (name) VALUES %s ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id;`,
				strings.Repeat(",(?)", len(genreArgs))[1:],
			),
			genreArgs...,
		)
		if err != nil {
			return err
		}

		defer func() {
			err = errors.Join(err, rows.Err(), rows.Close())
		}()

		genreIDs := make([]int64, len(genreArgs))
		index := 0

		for rows.Next() {
			if err = rows.Scan(&genreIDs[index]); err != nil {
				return err
			}

			index++
		}

		movieGenreArgs := make([]any, len(movie.Genres)*2)

		for i := range len(movie.Genres) {
			movieGenreArgs[i*2] = movie.ID
			movieGenreArgs[i*2+1] = genreIDs[i]
		}

		_, err = tx.ExecContext(ctx,
			fmt.Sprintf(
				`INSERT INTO movie_genres (movie_id, genre_id) VALUES %s;`,
				strings.Repeat(",(?, ?)", len(movie.Genres))[1:],
			),
			movieGenreArgs...,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// Query implements benchflix.Repository.
func (r Repository) Query(ctx context.Context, query benchflix.Query) ([]benchflix.Movie, error) {
	filters := []string{"1=1"}
	args := []any{}

	if query.Search != "" {
		filters = append(filters, `(
			EXISTS (
				SELECT 1 FROM movie_directors
				JOIN people ON people.id = movie_directors.person_id
				WHERE movie_directors.movie_id = movies.id AND INSTR(people.name, ?) > 0
			)
			OR EXISTS (
				SELECT 1 FROM movie_actors
				JOIN people ON people.id = movie_actors.person_id 
				WHERE movie_actors.movie_id = movies.id AND INSTR(people.name, ?) > 0
			)
		)`)

		args = append(args, query.Search, query.Search)
	}

	if query.Genre != "" {
		filters = append(filters, `EXISTS (
			SELECT 1 FROM movie_genres
			JOIN genres ON genres.id = movie_genres.genre_id
			WHERE movie_genres.movie_id = movies.id AND genres.name = ?
		)`)

		args = append(args, query.Genre)
	}

	if query.Country != "" {
		filters = append(filters, `EXISTS (
			SELECT 1 FROM movie_countries
			JOIN countries ON countries.id = movie_countries.country_id
			WHERE movie_countries.movie_id = movies.id AND countries.name = ?
		)`)

		args = append(args, query.Country)
	}

	if !query.AddedBefore.IsZero() {
		filters = append(filters, `added_at < ?`)
		args = append(args, query.AddedBefore)
	}

	if !query.AddedAfter.IsZero() {
		filters = append(filters, `added_at > ?`)
		args = append(args, query.AddedAfter)
	}

	if query.MinRating > 0 {
		filters = append(filters, `rating >= ?`)
		args = append(args, query.MinRating)
	}

	if query.MaxRating > 0 {
		filters = append(filters, `rating <= ?`)
		args = append(args, query.MaxRating)
	}

	rows, err := r.DB.QueryContext(ctx,
		fmt.Sprintf(
			`SELECT
				movies.id,
				movies.title,
				movies.added_at,
				movies.rating,
				(
					SELECT GROUP_CONCAT(people.name ORDER BY people.name)
					FROM movie_directors
					JOIN people ON people.id = movie_directors.person_id
					WHERE movie_directors.movie_id = movies.id
				) AS directors,
				(
					SELECT GROUP_CONCAT(people.name ORDER BY people.name)
					FROM movie_actors
					JOIN people ON people.id = movie_actors.person_id
					WHERE movie_actors.movie_id = movies.id
				) AS actors,
				(
					SELECT GROUP_CONCAT(countries.name ORDER BY countries.name)
					FROM movie_countries
					JOIN countries ON countries.id = movie_countries.country_id
					WHERE movie_countries.movie_id = movies.id
				) AS countries,
				(
					SELECT GROUP_CONCAT(genres.name ORDER BY genres.name)
					FROM movie_genres
					JOIN genres ON genres.id = movie_genres.genre_id
					WHERE movie_genres.movie_id = movies.id
				) AS genres
			FROM movies
			WHERE %s
			ORDER BY movies.title ASC;`,
			strings.Join(filters, " AND "),
		),
		args...,
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = errors.Join(err, rows.Err(), rows.Close())
	}()

	var movies []benchflix.Movie

	for rows.Next() {
		var (
			movie                                benchflix.Movie
			directors, actors, countries, genres sql.NullString
		)

		if err := rows.Scan(&movie.ID, &movie.Title, &movie.AddedAt, &movie.Rating, &directors, &actors, &countries, &genres); err != nil {
			return nil, err
		}

		if directors.Valid {
			movie.Directors = strings.Split(directors.String, ",")
		}

		if actors.Valid {
			movie.Actors = strings.Split(actors.String, ",")
		}

		if countries.Valid {
			movie.Countries = strings.Split(countries.String, ",")
		}

		if genres.Valid {
			movie.Genres = strings.Split(genres.String, ",")
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

// Read implements benchflix.Repository.
func (r Repository) Read(ctx context.Context, id int64) (benchflix.Movie, error) {
	row := r.DB.QueryRowContext(ctx,
		`SELECT
			movies.id,
			movies.title,
			movies.added_at,
			movies.rating,
			(
				SELECT GROUP_CONCAT(people.name ORDER BY people.name)
				FROM movie_directors
				JOIN people ON people.id = movie_directors.person_id
				WHERE movie_directors.movie_id = movies.id
			) AS directors,
			(
				SELECT GROUP_CONCAT(people.name ORDER BY people.name)
				FROM movie_actors
				JOIN people ON people.id = movie_actors.person_id
				WHERE movie_actors.movie_id = movies.id
			) AS actors,
			(
				SELECT GROUP_CONCAT(countries.name ORDER BY countries.name)
				FROM movie_countries
				JOIN countries ON countries.id = movie_countries.country_id
				WHERE movie_countries.movie_id = movies.id
			) AS countries,
			(
				SELECT GROUP_CONCAT(genres.name ORDER BY genres.name)
				FROM movie_genres
				JOIN genres ON genres.id = movie_genres.genre_id
				WHERE movie_genres.movie_id = movies.id
			) AS genres
		FROM movies
		WHERE id = ?
		ORDER BY movies.title ASC;`,
		id,
	)

	var (
		movie                                benchflix.Movie
		directors, actors, countries, genres string
	)

	if err := row.Scan(&movie.ID, &movie.Title, &movie.AddedAt, &movie.Rating, &directors, &actors, &countries, &genres); err != nil {
		return benchflix.Movie{}, err
	}

	movie.Directors = strings.Split(directors, ",")
	movie.Actors = strings.Split(actors, ",")
	movie.Countries = strings.Split(countries, ",")
	movie.Genres = strings.Split(genres, ",")

	return movie, nil
}
