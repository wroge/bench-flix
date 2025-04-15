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

func NewRepository(driverName, dataSourceName string) benchflix.Repository {
	db, err := sql.Open(driverName, dataSourceName)
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
			movie_id INTEGER REFERENCES movies (id) ON DELETE CASCADE,
			person_id INTEGER REFERENCES people (id) ON DELETE CASCADE,
			PRIMARY KEY (movie_id, person_id)
		);

		CREATE TABLE movie_actors (
			movie_id INTEGER REFERENCES movies (id) ON DELETE CASCADE,
			person_id INTEGER REFERENCES people (id) ON DELETE CASCADE,
			PRIMARY KEY (movie_id, person_id)
		);

		CREATE TABLE countries (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);

		CREATE TABLE movie_countries (
			movie_id INTEGER REFERENCES movies (id) ON DELETE CASCADE,
			country_id INTEGER REFERENCES countries (id) ON DELETE CASCADE,
			PRIMARY KEY (movie_id, country_id)
		);

		CREATE TABLE genres (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);

		CREATE TABLE movie_genres (
			movie_id INTEGER REFERENCES movies (id) ON DELETE CASCADE,
			genre_id INTEGER REFERENCES genres (id) ON DELETE CASCADE,
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

func (r Repository) Delete(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM movies WHERE id = ?", id)

	return err
}

func (r Repository) Create(ctx context.Context, movie benchflix.Movie) (err error) {
	var (
		actorsLen    = len(movie.Actors)
		directorsLen = len(movie.Directors)
		countriesLen = len(movie.Countries)
		genresLen    = len(movie.Genres)
	)

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

	if actorsLen+directorsLen > 0 {
		peopleArgs := make([]any, directorsLen+actorsLen)

		for i, p := range append(movie.Directors, movie.Actors...) {
			peopleArgs[i] = p
		}

		rows, err := tx.QueryContext(ctx,
			fmt.Sprintf(
				`INSERT INTO people (name) VALUES %s ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`,
				strings.Repeat(",(?)", len(peopleArgs))[1:],
			),
			peopleArgs...)
		if err != nil {
			return err
		}

		defer func() {
			err = errors.Join(err, rows.Err(), rows.Close())
		}()

		people := make([]int64, len(peopleArgs))
		index := 0

		for rows.Next() {
			if err = rows.Scan(&people[index]); err != nil {
				return err
			}

			index++
		}

		if directorsLen > 0 {
			movieDirectorArgs := make([]any, directorsLen*2)

			for i, id := range people[:directorsLen] {
				movieDirectorArgs[i*2] = movie.ID
				movieDirectorArgs[i*2+1] = id
			}

			_, err = tx.ExecContext(ctx,
				fmt.Sprintf(
					`INSERT INTO movie_directors (movie_id, person_id) VALUES %s;`,
					strings.Repeat(",(?, ?)", directorsLen)[1:],
				),
				movieDirectorArgs...,
			)
			if err != nil {
				return err
			}
		}

		if actorsLen > 0 {
			movieActorArgs := make([]any, actorsLen*2)

			for i, id := range people[directorsLen:] {
				movieActorArgs[i*2] = movie.ID
				movieActorArgs[i*2+1] = id
			}

			_, err = tx.ExecContext(ctx,
				fmt.Sprintf(
					`INSERT INTO movie_actors (movie_id, person_id) VALUES %s;`,
					strings.Repeat(",(?, ?)", actorsLen)[1:],
				),
				movieActorArgs...,
			)
			if err != nil {
				return err
			}
		}
	}

	if countriesLen > 0 {
		countryArgs := make([]any, countriesLen)

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

		movieCountryArgs := make([]any, countriesLen*2)

		for i := range countriesLen {
			movieCountryArgs[i*2] = movie.ID
			movieCountryArgs[i*2+1] = countryIDs[i]
		}

		_, err = tx.ExecContext(ctx,
			fmt.Sprintf(
				`INSERT INTO movie_countries (movie_id, country_id) VALUES %s;`,
				strings.Repeat(",(?, ?)", countriesLen)[1:],
			),
			movieCountryArgs...,
		)
		if err != nil {
			return err
		}
	}

	if genresLen > 0 {
		genreArgs := make([]any, genresLen)

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

		movieGenreArgs := make([]any, genresLen*2)

		for i := range genresLen {
			movieGenreArgs[i*2] = movie.ID
			movieGenreArgs[i*2+1] = genreIDs[i]
		}

		_, err = tx.ExecContext(ctx,
			fmt.Sprintf(
				`INSERT INTO movie_genres (movie_id, genre_id) VALUES %s;`,
				strings.Repeat(",(?, ?)", genresLen)[1:],
			),
			movieGenreArgs...,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r Repository) Query(ctx context.Context, query benchflix.Query) ([]benchflix.Movie, error) {
	builder := &strings.Builder{}
	args := []any{}

	if query.Search != "" {
		builder.WriteString(`AND (
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
		builder.WriteString(`AND EXISTS (
			SELECT 1 FROM movie_genres
			JOIN genres ON genres.id = movie_genres.genre_id
			WHERE movie_genres.movie_id = movies.id AND genres.name = ?
		)`)

		args = append(args, query.Genre)
	}

	if query.Country != "" {
		builder.WriteString(`AND EXISTS (
			SELECT 1 FROM movie_countries
			JOIN countries ON countries.id = movie_countries.country_id
			WHERE movie_countries.movie_id = movies.id AND countries.name = ?
		)`)

		args = append(args, query.Country)
	}

	if !query.AddedBefore.IsZero() {
		builder.WriteString(` AND added_at < ?`)

		args = append(args, query.AddedBefore)
	}

	if !query.AddedAfter.IsZero() {
		builder.WriteString(` AND added_at > ?`)

		args = append(args, query.AddedAfter)
	}

	if query.MinRating > 0 {
		builder.WriteString(` AND rating >= ?`)

		args = append(args, query.MinRating)
	}

	if query.MaxRating > 0 {
		builder.WriteString(` AND rating <= ?`)

		args = append(args, query.MaxRating)
	}

	builder.WriteString(" ORDER BY movies.title ASC")

	if query.Limit > 0 {
		builder.WriteString(" LIMIT ?")
		
		args = append(args, query.Limit)
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
			WHERE 1=1 %s;`,
			builder,
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

		movies = append(movies, ConvertMovie(movie, directors, actors, countries, genres))
	}

	return movies, nil
}

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
		directors, actors, countries, genres sql.NullString
	)

	if err := row.Scan(&movie.ID, &movie.Title, &movie.AddedAt, &movie.Rating, &directors, &actors, &countries, &genres); err != nil {
		return benchflix.Movie{}, err
	}

	return ConvertMovie(movie, directors, actors, countries, genres), nil
}

func ConvertMovie(movie benchflix.Movie, directors, actors, countries, genres sql.NullString) benchflix.Movie {
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

	return movie
}
