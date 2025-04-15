package sqlxflix

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	benchflix "github.com/wroge/bench-flix"
)

func NewRepository(driverName, dataSourceName string) benchflix.Repository {
	sqldb, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	db := sqlx.NewDb(sqldb, driverName)

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
	DB *sqlx.DB
}

func (r Repository) Delete(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM movies WHERE id = ?", id)

	return err
}

func (r Repository) Create(ctx context.Context, movie benchflix.Movie) (err error) {
	tx, err := r.DB.BeginTxx(ctx, nil)
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

	if len(movie.Directors) > 0 {
		directorNames := make([]any, len(movie.Directors))

		for i, p := range movie.Directors {
			directorNames[i] = p
		}

		directorIDs := make([]int64, len(directorNames))

		err = tx.SelectContext(ctx,
			&directorIDs,
			fmt.Sprintf(
				`INSERT INTO people (name) VALUES %s ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`,
				strings.Repeat(",(?)", len(directorNames))[1:],
			),
			directorNames...)
		if err != nil {
			return err
		}

		movieDirectorArgs := make([]any, len(directorIDs)*2)

		for i, id := range directorIDs {
			movieDirectorArgs[i*2] = movie.ID
			movieDirectorArgs[i*2+1] = id
		}

		_, err = tx.ExecContext(ctx,
			fmt.Sprintf(
				`INSERT INTO movie_directors (movie_id, person_id) VALUES %s;`,
				strings.Repeat(",(?, ?)", len(directorIDs))[1:],
			),
			movieDirectorArgs...,
		)
		if err != nil {
			return err
		}
	}

	if len(movie.Actors) > 0 {
		actorNames := make([]any, len(movie.Actors))

		for i, p := range movie.Actors {
			actorNames[i] = p
		}

		actorIDs := make([]int64, len(actorNames))

		err = tx.SelectContext(ctx,
			&actorIDs,
			fmt.Sprintf(
				`INSERT INTO people (name) VALUES %s ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`,
				strings.Repeat(",(?)", len(actorNames))[1:],
			),
			actorNames...)
		if err != nil {
			return err
		}

		movieActorArgs := make([]any, len(actorIDs)*2)

		for i, id := range actorIDs {
			movieActorArgs[i*2] = movie.ID
			movieActorArgs[i*2+1] = id
		}

		_, err = tx.ExecContext(ctx,
			fmt.Sprintf(
				`INSERT INTO movie_actors (movie_id, person_id) VALUES %s;`,
				strings.Repeat(",(?, ?)", len(actorIDs))[1:],
			),
			movieActorArgs...,
		)
		if err != nil {
			return err
		}
	}

	if len(movie.Countries) > 0 {
		countryArgs := make([]any, len(movie.Countries))

		for i, c := range movie.Countries {
			countryArgs[i] = c
		}

		countryIDs := make([]int64, len(countryArgs))

		err = tx.SelectContext(ctx,
			&countryIDs,
			fmt.Sprintf(
				`INSERT INTO countries (name) VALUES %s ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id;`,
				strings.Repeat(",(?)", len(countryArgs))[1:],
			),
			countryArgs...,
		)
		if err != nil {
			return err
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

		genreIDs := make([]int64, len(genreArgs))

		err = tx.SelectContext(ctx,
			&genreIDs,
			fmt.Sprintf(
				`INSERT INTO genres (name) VALUES %s ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id;`,
				strings.Repeat(",(?)", len(genreArgs))[1:],
			),
			genreArgs...,
		)
		if err != nil {
			return err
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

	var movies []Movie

	err := r.DB.SelectContext(ctx,
		&movies,
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

	result := make([]benchflix.Movie, len(movies))

	for i, m := range movies {
		result[i] = ConvertMovie(m)
	}

	return result, nil
}

type Movie struct {
	benchflix.Movie
	Directors, Actors, Countries, Genres sql.NullString
}

func (r Repository) Read(ctx context.Context, id int64) (benchflix.Movie, error) {
	var movie Movie

	err := r.DB.GetContext(ctx,
		&movie,
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
	if err != nil {
		return benchflix.Movie{}, err
	}

	return ConvertMovie(movie), nil
}

func ConvertMovie(movie Movie) benchflix.Movie {
	if movie.Directors.Valid {
		movie.Movie.Directors = strings.Split(movie.Directors.String, ",")
	}

	if movie.Actors.Valid {
		movie.Movie.Actors = strings.Split(movie.Actors.String, ",")
	}

	if movie.Countries.Valid {
		movie.Movie.Countries = strings.Split(movie.Countries.String, ",")
	}

	if movie.Genres.Valid {
		movie.Movie.Genres = strings.Split(movie.Genres.String, ",")
	}

	return movie.Movie
}
