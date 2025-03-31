package sqltflix

import (
	"context"
	"database/sql"

	"github.com/Masterminds/sprig/v3"
	_ "github.com/mattn/go-sqlite3"
	benchflix "github.com/wroge/bench-flix"
	"github.com/wroge/sqlt"
)

var (
	config = sqlt.Config{
		Templates: []sqlt.Template{
			sqlt.Funcs(sprig.TxtFuncMap()),
		},
		Cache: &sqlt.Cache{},
	}
	schema = sqlt.Exec[any](
		config,
		sqlt.Parse(`
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
			);
		`),
	)

	create = sqlt.Transaction(nil,
		sqlt.Exec[benchflix.Movie](
			config,
			sqlt.Parse(`
				INSERT INTO movies (id, title, added_at, rating) VALUES (
					{{ .ID }}, {{ .Title }}, {{ .AddedAt }}, {{ .Rating }}
				);
			`),
		),
		sqlt.All[benchflix.Movie, int64](
			config,
			sqlt.Name("DirectorIDs"),
			sqlt.Parse(`
				{{ if .Directors }}
					INSERT INTO people (name) VALUES 
						{{ range $i, $p := .Directors }}
							{{ if $i }}, {{ end }}
							({{ $p }})
						{{ end }}
					ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id;
				{{ end }}
			`),
		),
		sqlt.All[benchflix.Movie, int64](
			config,
			sqlt.Name("ActorIDs"),
			sqlt.Parse(`
				{{ if .Actors }}
					INSERT INTO people (name) VALUES 
						{{ range $i, $p := .Actors }}
							{{ if $i }}, {{ end }}
							({{ $p }})
						{{ end }}
					ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id;
				{{ end }}
			`),
		),
		sqlt.Exec[benchflix.Movie](
			config,
			sqlt.Parse(`
				{{ if .Directors }}
					INSERT INTO movie_directors (movie_id, person_id) VALUES
						{{ range $i, $id := (Context "DirectorIDs") }}
							{{ if $i }}, {{ end }}
							({{ $.ID }}, {{ $id }})
						{{ end }}
				{{ end }}
			`),
		),
		sqlt.Exec[benchflix.Movie](
			config,
			sqlt.Parse(`
				{{ if .Actors }}
					INSERT INTO movie_actors (movie_id, person_id) VALUES
						{{ range $i, $id := (Context "ActorIDs") }}
							{{ if $i }}, {{ end }}
							({{ $.ID }}, {{ $id }})
						{{ end }}
				{{ end }}
			`),
		),
		sqlt.All[benchflix.Movie, int64](
			config,
			sqlt.Name("CountryIDs"),
			sqlt.Parse(`
				{{ if .Countries }}
					INSERT INTO countries (name) VALUES 
						{{ range $i, $p := .Countries }}
							{{ if $i }}, {{ end }}
							({{ $p }})
						{{ end }}
					ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id;
				{{ end }}
			`),
		),
		sqlt.Exec[benchflix.Movie](
			config,
			sqlt.Parse(`
				{{ if .Countries }}
					INSERT INTO movie_countries (movie_id, country_id) VALUES
						{{ range $i, $id := (Context "CountryIDs") }}
							{{ if $i }}, {{ end }}
							({{ $.ID }}, {{ $id }})
						{{ end }}
				{{ end }}
			`),
		),
		sqlt.All[benchflix.Movie, int64](
			config,
			sqlt.Name("GenreIDs"),
			sqlt.Parse(`
				{{ if .Genres }}
					INSERT INTO genres (name) VALUES 
						{{ range $i, $p := .Genres }}
							{{ if $i }}, {{ end }}
							({{ $p }})
						{{ end }}
					ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id;
				{{ end }}
			`),
		),
		sqlt.Exec[benchflix.Movie](
			config,
			sqlt.Parse(`
				{{ if .Genres }}
					INSERT INTO movie_genres (movie_id, genre_id) VALUES
						{{ range $i, $id := (Context "GenreIDs") }}
							{{ if $i }}, {{ end }}
							({{ $.ID }}, {{ $id }})
						{{ end }}
				{{ end }}
			`),
		),
	)

	queryByID = sqlt.First[int64, benchflix.Movie](
		config,
		sqlt.Parse(`
			SELECT
				movies.id,				{{ Scan "ID" }}
				movies.title,			{{ Scan "Title" }}
				movies.added_at,		{{ Scan "AddedAt" }}
				movies.rating,			{{ Scan "Rating" }}
				(
					SELECT GROUP_CONCAT(people.name ORDER BY people.name)
					FROM movie_directors
					JOIN people ON people.id = movie_directors.person_id
					WHERE movie_directors.movie_id = movies.id
				) AS directors,			{{ ScanSplit "Directors" "," }}
				(
					SELECT GROUP_CONCAT(people.name ORDER BY people.name)
					FROM movie_actors
					JOIN people ON people.id = movie_actors.person_id
					WHERE movie_actors.movie_id = movies.id
				) AS actors,			{{ ScanSplit "Actors" "," }}
				(
					SELECT GROUP_CONCAT(countries.name ORDER BY countries.name)
					FROM movie_countries
					JOIN countries ON countries.id = movie_countries.country_id
					WHERE movie_countries.movie_id = movies.id
				) AS countries,			{{ ScanSplit "Countries" "," }}
				(
					SELECT GROUP_CONCAT(genres.name ORDER BY genres.name)
					FROM movie_genres
					JOIN genres ON genres.id = movie_genres.genre_id
					WHERE movie_genres.movie_id = movies.id
				) AS genres				{{ ScanSplit "Genres" "," }}
			FROM movies
			WHERE movies.id = {{ . }}
			ORDER BY movies.title ASC;
		`),
	)

	queryAll = sqlt.All[benchflix.Query, benchflix.Movie](
		config,
		sqlt.Parse(`
			SELECT
				movies.id,				{{ Scan "ID" }}
				movies.title,			{{ Scan "Title" }}
				movies.added_at,		{{ Scan "AddedAt" }}
				movies.rating,			{{ Scan "Rating" }}
				(
					SELECT GROUP_CONCAT(people.name ORDER BY people.name)
					FROM movie_directors
					JOIN people ON people.id = movie_directors.person_id
					WHERE movie_directors.movie_id = movies.id
				) AS directors,			{{ ScanSplit "Directors" "," }}
				(
					SELECT GROUP_CONCAT(people.name ORDER BY people.name)
					FROM movie_actors
					JOIN people ON people.id = movie_actors.person_id
					WHERE movie_actors.movie_id = movies.id
				) AS actors,			{{ ScanSplit "Actors" "," }}
				(
					SELECT GROUP_CONCAT(countries.name ORDER BY countries.name)
					FROM movie_countries
					JOIN countries ON countries.id = movie_countries.country_id
					WHERE movie_countries.movie_id = movies.id
				) AS countries,			{{ ScanSplit "Countries" "," }}
				(
					SELECT GROUP_CONCAT(genres.name ORDER BY genres.name)
					FROM movie_genres
					JOIN genres ON genres.id = movie_genres.genre_id
					WHERE movie_genres.movie_id = movies.id
				) AS genres				{{ ScanSplit "Genres" "," }}
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
			ORDER BY movies.title ASC;
		`),
	)
)

func NewRepository() benchflix.Repository {
	db, err := sql.Open("sqlite3", ":memory:?_fk=1")
	if err != nil {
		panic(err)
	}

	_, err = schema.Exec(context.Background(), db, nil)
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
func (r Repository) Create(ctx context.Context, movie benchflix.Movie) error {
	_, err := create.Exec(ctx, r.DB, movie)

	return err
}

// Query implements benchflix.Repository.
func (r Repository) Query(ctx context.Context, query benchflix.Query) ([]benchflix.Movie, error) {
	return queryAll.Exec(ctx, r.DB, query)
}

// Read implements benchflix.Repository.
func (r Repository) Read(ctx context.Context, id int64) (benchflix.Movie, error) {
	return queryByID.Exec(ctx, r.DB, id)
}
