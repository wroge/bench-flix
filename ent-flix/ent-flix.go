package entflix

import (
	"context"
	"errors"

	"entgo.io/ent/dialect/sql"
	_ "github.com/mattn/go-sqlite3"
	benchflix "github.com/wroge/bench-flix"
	"github.com/wroge/bench-flix/ent-flix/ent"
	"github.com/wroge/bench-flix/ent-flix/ent/country"
	"github.com/wroge/bench-flix/ent-flix/ent/genre"
	"github.com/wroge/bench-flix/ent-flix/ent/movie"
	"github.com/wroge/bench-flix/ent-flix/ent/person"
)

func NewRepository() benchflix.Repository {
	client, err := ent.Open("sqlite3", ":memory:?_fk=1")
	if err != nil {
		panic(err)
	}

	if err = client.Schema.Create(context.Background()); err != nil {
		panic(err)
	}

	return Repository{
		Client: client,
	}
}

type Repository struct {
	Client *ent.Client
}

// Create implements benchflix.Repository.
func (r Repository) Create(ctx context.Context, movie benchflix.Movie) (err error) {
	tx, err := r.Client.BeginTx(ctx, nil)
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

	create := tx.Movie.Create().
		SetID(movie.ID).
		SetTitle(movie.Title).
		SetAddedAt(movie.AddedAt).
		SetRating(movie.Rating)

	if len(movie.Directors) > 0 {
		people := make([]int64, len(movie.Directors))

		for i, p := range movie.Directors {
			people[i], err = tx.Person.Create().SetName(p).OnConflict().UpdateName().ID(ctx)
			if err != nil {
				return err
			}
		}

		create.AddDirectorIDs(people...)
	}

	if len(movie.Actors) > 0 {
		people := make([]int64, len(movie.Actors))

		for i, p := range movie.Actors {
			people[i], err = tx.Person.Create().SetName(p).OnConflict().UpdateName().ID(ctx)
			if err != nil {
				return err
			}
		}

		create.AddActorIDs(people...)
	}

	if len(movie.Countries) > 0 {
		countries := make([]int64, len(movie.Countries))

		for i, p := range movie.Countries {
			countries[i], err = tx.Country.Create().SetName(p).OnConflict().UpdateName().ID(ctx)
			if err != nil {
				return err
			}
		}

		create.AddCountryIDs(countries...)
	}

	if len(movie.Genres) > 0 {
		genres := make([]int64, len(movie.Genres))

		for i, p := range movie.Genres {
			genres[i], err = tx.Genre.Create().SetName(p).OnConflict().UpdateName().ID(ctx)
			if err != nil {
				return err
			}
		}

		create.AddGenreIDs(genres...)
	}

	return create.Exec(ctx)
}

// Query implements benchflix.Repository.
func (r Repository) Query(ctx context.Context, query benchflix.Query) ([]benchflix.Movie, error) {
	q := r.Client.Movie.Query().
		WithDirectors(
			func(ptq *ent.PersonQuery) {
				ptq.Order(person.ByName(sql.OrderAsc()))
			},
		).WithActors(
		func(ptq *ent.PersonQuery) {
			ptq.Order(person.ByName(sql.OrderAsc()))
		}).
		WithCountries(
			func(ptq *ent.CountryQuery) {
				ptq.Order(country.ByName(sql.OrderAsc()))
			}).
		WithGenres(
			func(ptq *ent.GenreQuery) {
				ptq.Order(genre.ByName(sql.OrderAsc()))
			},
		).
		Order(movie.ByTitle(sql.OrderAsc()))

	if query.Search != "" {
		q.Where(movie.Or(
			movie.HasDirectorsWith(person.NameContains(query.Search)),
			movie.HasActorsWith(person.NameContains(query.Search)),
		))
	}

	if query.Genre != "" {
		q.Where(movie.HasGenresWith(genre.Name(query.Genre)))
	}

	if query.Country != "" {
		q.Where(movie.HasCountriesWith(country.Name(query.Country)))
	}

	if !query.AddedAfter.IsZero() {
		q.Where(movie.AddedAtGT(query.AddedAfter))
	}

	if !query.AddedBefore.IsZero() {
		q.Where(movie.AddedAtLT(query.AddedBefore))
	}

	if query.MinRating > 0 {
		q.Where(movie.RatingGTE(query.MinRating))
	}

	if query.MaxRating > 0 {
		q.Where(movie.RatingLTE(query.MaxRating))
	}

	result, err := q.All(ctx)
	if err != nil {
		return nil, err
	}

	movies := make([]benchflix.Movie, len(result))

	for i, each := range result {
		movies[i] = toMovie(each)
	}

	return movies, nil
}

// Read implements benchflix.Repository.
func (r Repository) Read(ctx context.Context, id int64) (benchflix.Movie, error) {
	result, err := r.Client.Movie.Query().Where(movie.ID(id)).
		WithDirectors(
			func(ptq *ent.PersonQuery) {
				ptq.Order(person.ByName(sql.OrderAsc()))
			},
		).WithActors(
		func(ptq *ent.PersonQuery) {
			ptq.Order(person.ByName(sql.OrderAsc()))
		}).
		WithCountries(
			func(ptq *ent.CountryQuery) {
				ptq.Order(country.ByName(sql.OrderAsc()))
			}).
		WithGenres(
			func(ptq *ent.GenreQuery) {
				ptq.Order(genre.ByName(sql.OrderAsc()))
			},
		).
		Order(movie.ByTitle(sql.OrderAsc())).
		Only(ctx)
	if err != nil {
		return benchflix.Movie{}, err
	}

	return toMovie(result), nil
}

func toMovie(result *ent.Movie) benchflix.Movie {
	movie := benchflix.Movie{
		ID:        result.ID,
		Title:     result.Title,
		AddedAt:   result.AddedAt,
		Rating:    result.Rating,
		Directors: make([]string, len(result.Edges.Directors)),
		Actors:    make([]string, len(result.Edges.Actors)),
		Countries: make([]string, len(result.Edges.Countries)),
		Genres:    make([]string, len(result.Edges.Genres)),
	}

	for j, v := range result.Edges.Directors {
		movie.Directors[j] = v.Name
	}

	for j, v := range result.Edges.Actors {
		movie.Actors[j] = v.Name
	}

	for j, v := range result.Edges.Countries {
		movie.Countries[j] = v.Name
	}

	for j, v := range result.Edges.Genres {
		movie.Genres[j] = v.Name
	}

	return movie
}
