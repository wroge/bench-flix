package sqlcflix

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	benchflix "github.com/wroge/bench-flix"
	"github.com/wroge/bench-flix/sqlc-flix/internal/db"
)

//go:embed schema.sql
var ddl string

func NewRepository() benchflix.Repository {
	sqldb, err := sql.Open("sqlite3", ":memory:?_fk=1")
	if err != nil {
		panic(err)
	}

	if _, err := sqldb.Exec(ddl); err != nil {
		panic(err)
	}

	return Repository{
		DB: sqldb,
	}
}

type Repository struct {
	DB *sql.DB
}

// Create inserts a movie and all its related data
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

	txdb := db.New(tx)

	if _, err := txdb.CreateMovie(ctx, db.CreateMovieParams{
		ID:      movie.ID,
		Title:   movie.Title,
		AddedAt: movie.AddedAt,
		Rating:  movie.Rating,
	}); err != nil {
		return err
	}

	for _, name := range movie.Directors {
		id, err := txdb.GetOrCreatePerson(ctx, name)
		if err != nil {
			return err
		}
		if err := txdb.AddMovieDirector(ctx, db.AddMovieDirectorParams{
			MovieID:  movie.ID,
			PersonID: id,
		}); err != nil {
			return err
		}
	}

	for _, name := range movie.Actors {
		id, err := txdb.GetOrCreatePerson(ctx, name)
		if err != nil {
			return err
		}
		if err := txdb.AddMovieActor(ctx, db.AddMovieActorParams{
			MovieID:  movie.ID,
			PersonID: id,
		}); err != nil {
			return err
		}
	}

	for _, name := range movie.Countries {
		id, err := txdb.GetOrCreateCountry(ctx, name)
		if err != nil {
			return err
		}
		if err := txdb.AddMovieCountry(ctx, db.AddMovieCountryParams{
			MovieID:   movie.ID,
			CountryID: id,
		}); err != nil {
			return err
		}
	}

	for _, name := range movie.Genres {
		id, err := txdb.GetOrCreateGenre(ctx, name)
		if err != nil {
			return err
		}
		if err := txdb.AddMovieGenre(ctx, db.AddMovieGenreParams{
			MovieID: movie.ID,
			GenreID: id,
		}); err != nil {
			return err
		}
	}

	return nil
}

// Read loads a movie and all its related string slices
func (r Repository) Read(ctx context.Context, id int64) (benchflix.Movie, error) {
	row, err := db.New(r.DB).GetMovie(ctx, id)
	if err != nil {
		return benchflix.Movie{}, err
	}

	return benchflix.Movie{
		ID:        row.ID,
		Title:     row.Title,
		AddedAt:   row.AddedAt,
		Rating:    row.Rating,
		Directors: splitCSV(row.Directors),
		Actors:    splitCSV(row.Actors),
		Countries: splitCSV(row.Countries),
		Genres:    splitCSV(row.Genres),
	}, nil
}

// Query performs a filtered search
func (r Repository) Query(ctx context.Context, q benchflix.Query) ([]benchflix.Movie, error) {
	params := db.QueryMoviesParams{
		Search:      q.Search,
		Genre:       q.Genre,
		Country:     q.Country,
		AddedAfter:  sql.NullTime{Time: q.AddedAfter, Valid: !q.AddedAfter.IsZero()},
		AddedBefore: sql.NullTime{Time: q.AddedBefore, Valid: !q.AddedBefore.IsZero()},
		MinRating:   q.MinRating,
		MaxRating:   q.MaxRating,
	}

	rows, err := db.New(r.DB).QueryMovies(ctx, params)
	if err != nil {
		return nil, err
	}

	var movies []benchflix.Movie
	for _, row := range rows {
		movies = append(movies, ConvertMovie(row))
	}

	return movies, nil
}

func ConvertMovie(m db.QueryMoviesRow) benchflix.Movie {
	return benchflix.Movie{
		ID:        m.ID,
		Title:     m.Title,
		AddedAt:   m.AddedAt,
		Rating:    m.Rating,
		Directors: splitCSV(m.Directors),
		Actors:    splitCSV(m.Actors),
		Countries: splitCSV(m.Countries),
		Genres:    splitCSV(m.Genres),
	}
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}

	return strings.Split(s, ",")
}
