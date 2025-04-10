package sqltflix

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	benchflix "github.com/wroge/bench-flix"
	"github.com/wroge/sqlt"
)

var (
	schema = sqlt.Exec[any](sqlt.Parse(`
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
		);
	`))

	create = sqlt.Transaction(nil,
		sqlt.Exec[benchflix.Movie](sqlt.Parse(`
			INSERT INTO movies (id, title, added_at, rating) VALUES 
			({{ .ID }}, {{ .Title }}, {{ .AddedAt }}, {{ .Rating }});
		`)),
		sqlt.All[benchflix.Movie, int64](sqlt.Name("DirectorIDs"), sqlt.Parse(`
			{{ if .Directors }}
				INSERT INTO people (name) VALUES 
				{{ range $i, $p := .Directors }}
					{{ if $i }}, {{ end }}
					({{ $p }})
				{{ end }}
				ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id;
			{{ end }}
		`)),
		sqlt.All[benchflix.Movie, int64](sqlt.Name("ActorIDs"), sqlt.Parse(`
			{{ if .Actors }}
				INSERT INTO people (name) VALUES 
				{{ range $i, $p := .Actors }}
					{{ if $i }}, {{ end }}
					({{ $p }})
				{{ end }}
				ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id;
			{{ end }}
		`)),
		sqlt.Exec[benchflix.Movie](sqlt.Parse(`
			{{ if .Directors }}
				INSERT INTO movie_directors (movie_id, person_id) VALUES
				{{ range $i, $id := (Context "DirectorIDs") }}
					{{ if $i }}, {{ end }}
					({{ $.ID }}, {{ $id }})
				{{ end }}
			{{ end }}
		`)),
		sqlt.Exec[benchflix.Movie](sqlt.Parse(`
			{{ if .Actors }}
				INSERT INTO movie_actors (movie_id, person_id) VALUES
				{{ range $i, $id := (Context "ActorIDs") }}
					{{ if $i }}, {{ end }}
					({{ $.ID }}, {{ $id }})
				{{ end }}
			{{ end }}
		`)),
		sqlt.All[benchflix.Movie, int64](sqlt.Name("CountryIDs"), sqlt.Parse(`
			{{ if .Countries }}
				INSERT INTO countries (name) VALUES 
				{{ range $i, $p := .Countries }}
					{{ if $i }}, {{ end }}
					({{ $p }})
				{{ end }}
				ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id;
			{{ end }}
		`)),
		sqlt.Exec[benchflix.Movie](sqlt.Parse(`
			{{ if .Countries }}
				INSERT INTO movie_countries (movie_id, country_id) VALUES
				{{ range $i, $id := (Context "CountryIDs") }}
					{{ if $i }}, {{ end }}
					({{ $.ID }}, {{ $id }})
				{{ end }}
			{{ end }}
		`)),
		sqlt.All[benchflix.Movie, int64](sqlt.Name("GenreIDs"), sqlt.Parse(`
			{{ if .Genres }}
				INSERT INTO genres (name) VALUES 
				{{ range $i, $p := .Genres }}
					{{ if $i }}, {{ end }}
					({{ $p }})
				{{ end }}
				ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id;
			{{ end }}
		`)),
		sqlt.Exec[benchflix.Movie](sqlt.Parse(`
			{{ if .Genres }}
				INSERT INTO movie_genres (movie_id, genre_id) VALUES
				{{ range $i, $id := (Context "GenreIDs") }}
					{{ if $i }}, {{ end }}
					({{ $.ID }}, {{ $id }})
				{{ end }}
			{{ end }}
		`)),
	)

	first = sqlt.First[int64, benchflix.Movie](sqlt.Cache{}, sqlt.Parse(`
		SELECT
			movies.id,		{{ ScanInt "ID" }}
			movies.title,		{{ ScanString "Title" }}
			movies.added_at,	{{ ScanTime "AddedAt" }}
			movies.rating,		{{ ScanFloat "Rating" }}
			(
				SELECT GROUP_CONCAT(people.name ORDER BY people.name)
				FROM movie_directors
				JOIN people ON people.id = movie_directors.person_id
				WHERE movie_directors.movie_id = movies.id
			) AS directors,		{{ ScanStringSlice "Directors" "," }}
			(
				SELECT GROUP_CONCAT(people.name ORDER BY people.name)
				FROM movie_actors
				JOIN people ON people.id = movie_actors.person_id
				WHERE movie_actors.movie_id = movies.id
			) AS actors,		{{ ScanStringSlice "Actors" "," }}
			(
				SELECT GROUP_CONCAT(countries.name ORDER BY countries.name)
				FROM movie_countries
				JOIN countries ON countries.id = movie_countries.country_id
				WHERE movie_countries.movie_id = movies.id
			) AS countries,		{{ ScanStringSlice "Countries" "," }}
			(
				SELECT GROUP_CONCAT(genres.name ORDER BY genres.name)
				FROM movie_genres
				JOIN genres ON genres.id = movie_genres.genre_id
				WHERE movie_genres.movie_id = movies.id
			) AS genres 		{{ ScanStringSlice "Genres" "," }}
		FROM movies
		WHERE movies.id = {{ . }}
		ORDER BY movies.title ASC;
	`))

	all = sqlt.All[benchflix.Query, benchflix.Movie](sqlt.Cache{}, sqlt.Parse(`
		SELECT
			movies.id,		{{ ScanInt "ID" }}
			movies.title,		{{ ScanString "Title" }}
			movies.added_at,	{{ ScanTime "AddedAt" }}
			movies.rating,		{{ ScanFloat "Rating" }}
			(
				SELECT GROUP_CONCAT(people.name ORDER BY people.name)
				FROM movie_directors
				JOIN people ON people.id = movie_directors.person_id
				WHERE movie_directors.movie_id = movies.id
			) AS directors,		{{ ScanStringSlice "Directors" "," }}
			(
				SELECT GROUP_CONCAT(people.name ORDER BY people.name)
				FROM movie_actors
				JOIN people ON people.id = movie_actors.person_id
				WHERE movie_actors.movie_id = movies.id
			) AS actors,		{{ ScanStringSlice "Actors" "," }}
			(
				SELECT GROUP_CONCAT(countries.name ORDER BY countries.name)
				FROM movie_countries
				JOIN countries ON countries.id = movie_countries.country_id
				WHERE movie_countries.movie_id = movies.id
			) AS countries,		{{ ScanStringSlice "Countries" "," }}
			(
				SELECT GROUP_CONCAT(genres.name ORDER BY genres.name)
				FROM movie_genres
				JOIN genres ON genres.id = movie_genres.genre_id
				WHERE movie_genres.movie_id = movies.id
			) AS genres 		{{ ScanStringSlice "Genres" "," }}
		FROM movies
		WHERE 1=1
		{{ if .Search }}
			AND (
				EXISTS (
					SELECT 1 FROM movie_directors
					JOIN people ON people.id = movie_directors.person_id
					WHERE movie_directors.movie_id = movies.id
					AND INSTR(people.name, {{ .Search }}) > 0
				)
				OR EXISTS (
					SELECT 1 FROM movie_actors
					JOIN people ON people.id = movie_actors.person_id
					WHERE movie_actors.movie_id = movies.id
					AND INSTR(people.name, {{ .Search }}) > 0
				)
			)
		{{ end }}
		{{ if .Genre }}
			AND EXISTS (
				SELECT 1 FROM movie_genres
				JOIN genres ON genres.id = movie_genres.genre_id
				WHERE movie_genres.movie_id = movies.id
				AND genres.name = {{ .Genre }}
			)
		{{ end }}
		{{ if .Country }}
			AND EXISTS (
				SELECT 1 FROM movie_countries
				JOIN countries ON countries.id = movie_countries.country_id
				WHERE movie_countries.movie_id = movies.id
				AND countries.name = {{ .Country }}
			)
		{{ end }}
		{{ if not .AddedBefore.IsZero }}
			AND added_at < {{ .AddedBefore }}
		{{ end }}
		{{ if not .AddedAfter.IsZero }}
			AND added_at > {{ .AddedAfter }}
		{{ end }}
		{{ if .MinRating }}
			AND rating >= {{ .MinRating }}
		{{ end }}
		{{ if .MaxRating }}
			AND rating <= {{ .MaxRating }}
		{{ end }}
		ORDER BY movies.title ASC
		{{ if .Limit }}
			LIMIT {{ .Limit }}
		{{ end }};
	`))

	deleteMovie = sqlt.Exec[int64](sqlt.Cache{}, sqlt.Parse(`
		DELETE FROM movies WHERE id = {{ . }};
	`))
)

func NewRepository(driverName, dataSourceName string) benchflix.Repository {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	if _, err = schema.Exec(context.Background(), db, nil); err != nil {
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
	_, err := deleteMovie.Exec(ctx, r.DB, id)

	return err
}

func (r Repository) Create(ctx context.Context, movie benchflix.Movie) error {
	_, err := create.Exec(ctx, r.DB, movie)

	return err
}

func (r Repository) Query(ctx context.Context, query benchflix.Query) ([]benchflix.Movie, error) {
	return all.Exec(ctx, r.DB, query)
}

func (r Repository) Read(ctx context.Context, id int64) (benchflix.Movie, error) {
	return first.Exec(ctx, r.DB, id)
}
