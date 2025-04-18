package sqlcflix

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	benchflix "github.com/wroge/bench-flix"
	"github.com/wroge/bench-flix/sqlc-flix/internal/db"
)

//go:embed schema.sql
var ddl string

func NewRepository(driverName, dataSourceName string) benchflix.Repository {
	sqldb, err := sql.Open(driverName, dataSourceName)
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

func (r Repository) Delete(ctx context.Context, id int64) error {
	return db.New(r.DB).DeleteMovie(ctx, id)
}

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

func (r Repository) Query(ctx context.Context, q benchflix.Query) ([]benchflix.Movie, error) {
	params := db.QueryMoviesParams{
		Search:      q.Search,
		Genre:       q.Genre,
		Country:     q.Country,
		AddedAfter:  sql.NullTime{Time: q.AddedAfter, Valid: !q.AddedAfter.IsZero()},
		AddedBefore: sql.NullTime{Time: q.AddedBefore, Valid: !q.AddedBefore.IsZero()},
		MinRating:   q.MinRating,
		MaxRating:   q.MaxRating,
		Limit:       q.Limit,
	}

	rows, err := db.New(r.DB).QueryMovies(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("HERE: %w", err)
	}

	movies := make([]benchflix.Movie, len(rows))
	for i, row := range rows {
		movies[i] = benchflix.Movie{
			ID:        row.ID,
			Title:     row.Title,
			AddedAt:   row.AddedAt,
			Rating:    row.Rating,
			Directors: splitCSV(row.Directors),
			Actors:    splitCSV(row.Actors),
			Countries: splitCSV(row.Countries),
			Genres:    splitCSV(row.Genres),
		}
	}

	return movies, nil
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}

	return strings.Split(s, ",")
}
