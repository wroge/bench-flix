package bunflix

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	benchflix "github.com/wroge/bench-flix"
)

type Movie struct {
	bun.BaseModel `bun:"table:movies"`

	ID      int64     `bun:",pk,autoincrement"`
	Title   string    `bun:",notnull"`
	AddedAt time.Time `bun:",notnull"`
	Rating  float64   `bun:",notnull"`

	Directors []*Person  `bun:"m2m:movie_directors"`
	Actors    []*Person  `bun:"m2m:movie_actors"`
	Countries []*Country `bun:"m2m:movie_countries"`
	Genres    []*Genre   `bun:"m2m:movie_genres"`
}

type Person struct {
	bun.BaseModel `bun:"table:people"`

	ID   int64  `bun:",pk,autoincrement"`
	Name string `bun:",unique,notnull"`
}

type Country struct {
	bun.BaseModel `bun:"table:countries"`

	ID   int64  `bun:",pk,autoincrement"`
	Name string `bun:",unique,notnull"`
}

type Genre struct {
	bun.BaseModel `bun:"table:genres"`

	ID   int64  `bun:",pk,autoincrement"`
	Name string `bun:",unique,notnull"`
}

type MovieDirector struct {
	bun.BaseModel `bun:"table:movie_directors"`

	MovieID  int64 `bun:",pk"`
	PersonID int64 `bun:",pk"`

	Movie  *Movie  `bun:"rel:belongs-to,join:movie_id=id"`
	Person *Person `bun:"rel:belongs-to,join:person_id=id"`
}

type MovieActor struct {
	bun.BaseModel `bun:"table:movie_actors"`

	MovieID  int64 `bun:",pk"`
	PersonID int64 `bun:",pk"`

	Movie  *Movie  `bun:"rel:belongs-to,join:movie_id=id"`
	Person *Person `bun:"rel:belongs-to,join:person_id=id"`
}

type MovieCountry struct {
	bun.BaseModel `bun:"table:movie_countries"`

	MovieID   int64 `bun:",pk"`
	CountryID int64 `bun:",pk"`

	Movie   *Movie   `bun:"rel:belongs-to,join:movie_id=id"`
	Country *Country `bun:"rel:belongs-to,join:country_id=id"`
}

type MovieGenre struct {
	bun.BaseModel `bun:"table:movie_genres"`

	MovieID int64 `bun:",pk"`
	GenreID int64 `bun:",pk"`

	Movie *Movie `bun:"rel:belongs-to,join:movie_id=id"`
	Genre *Genre `bun:"rel:belongs-to,join:genre_id=id"`
}

func NewRepository() benchflix.Repository {
	sqldb, err := sql.Open("sqlite3", ":memory:?_fk=1")
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())

	db.RegisterModel(
		(*MovieDirector)(nil),
		(*MovieActor)(nil),
		(*MovieCountry)(nil),
		(*MovieGenre)(nil),
	)

	if _, err = db.NewCreateTable().Model((*Movie)(nil)).Exec(context.Background()); err != nil {
		panic(err)
	}

	if _, err = db.NewCreateTable().Model((*Person)(nil)).Exec(context.Background()); err != nil {
		panic(err)
	}

	if _, err = db.NewCreateTable().Model((*Country)(nil)).Exec(context.Background()); err != nil {
		panic(err)
	}

	if _, err = db.NewCreateTable().Model((*Genre)(nil)).Exec(context.Background()); err != nil {
		panic(err)
	}

	if _, err = db.NewCreateTable().Model((*MovieDirector)(nil)).
		ForeignKey(`(movie_id) REFERENCES movies(id) ON DELETE CASCADE`).
		ForeignKey(`(person_id) REFERENCES people(id) ON DELETE CASCADE`).Exec(context.Background()); err != nil {
		panic(err)
	}

	if _, err = db.NewCreateTable().Model((*MovieActor)(nil)).
		ForeignKey(`(movie_id) REFERENCES movies(id) ON DELETE CASCADE`).
		ForeignKey(`(person_id) REFERENCES people(id) ON DELETE CASCADE`).Exec(context.Background()); err != nil {
		panic(err)
	}

	if _, err = db.NewCreateTable().Model((*MovieCountry)(nil)).
		ForeignKey(`(movie_id) REFERENCES movies(id) ON DELETE CASCADE`).
		ForeignKey(`(country_id) REFERENCES countries(id) ON DELETE CASCADE`).Exec(context.Background()); err != nil {
		panic(err)
	}

	if _, err = db.NewCreateTable().Model((*MovieGenre)(nil)).
		ForeignKey(`(movie_id) REFERENCES movies(id) ON DELETE CASCADE`).
		ForeignKey(`(genre_id) REFERENCES genres(id) ON DELETE CASCADE`).Exec(context.Background()); err != nil {
		panic(err)
	}

	return Repository{
		DB: db,
	}
}

type Repository struct {
	DB *bun.DB
}

func (r Repository) Delete(ctx context.Context, id int64) error {
	_, err := r.DB.NewDelete().Model(&Movie{}).Where("id = ?", id).Exec(ctx)

	return err
}

func (r Repository) Create(ctx context.Context, movie benchflix.Movie) error {
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

	if _, err := tx.NewInsert().Model(&Movie{
		ID:      movie.ID,
		Title:   movie.Title,
		AddedAt: movie.AddedAt,
		Rating:  movie.Rating,
	}).Exec(ctx); err != nil {
		return err
	}

	directorIDs := make([]int64, len(movie.Directors))

	for i, d := range movie.Directors {
		if err = tx.NewInsert().Model(&Person{
			Name: d,
		}).On("CONFLICT (name) DO UPDATE").Set("name = EXCLUDED.name").Returning("id").Scan(ctx, &directorIDs[i]); err != nil {
			return err
		}
	}

	if len(directorIDs) > 0 {
		movieDirectors := make([]MovieDirector, len(directorIDs))

		for i, d := range directorIDs {
			movieDirectors[i] = MovieDirector{
				MovieID:  movie.ID,
				PersonID: d,
			}
		}

		if _, err = tx.NewInsert().Model(&movieDirectors).Exec(ctx); err != nil {
			return err
		}
	}

	actorIDs := make([]int64, len(movie.Actors))

	for i, a := range movie.Actors {
		if err = tx.NewInsert().Model(&Person{
			Name: a,
		}).On("CONFLICT (name) DO UPDATE").Set("name = EXCLUDED.name").Returning("id").Scan(ctx, &actorIDs[i]); err != nil {
			return err
		}
	}

	if len(actorIDs) > 0 {
		movieActors := make([]MovieActor, len(actorIDs))

		for i, a := range actorIDs {
			movieActors[i] = MovieActor{
				MovieID:  movie.ID,
				PersonID: a,
			}
		}

		if _, err = tx.NewInsert().Model(&movieActors).Exec(ctx); err != nil {
			return err
		}
	}

	countryIDs := make([]int64, len(movie.Countries))

	for i, d := range movie.Countries {
		if err = tx.NewInsert().Model(&Country{
			Name: d,
		}).On("CONFLICT (name) DO UPDATE").Set("name = EXCLUDED.name").Returning("id").Scan(ctx, &countryIDs[i]); err != nil {
			return err
		}
	}

	if len(countryIDs) > 0 {
		movieCountries := make([]MovieCountry, len(countryIDs))

		for i, a := range countryIDs {
			movieCountries[i] = MovieCountry{
				MovieID:   movie.ID,
				CountryID: a,
			}
		}

		if _, err = tx.NewInsert().Model(&movieCountries).Exec(ctx); err != nil {
			return err
		}
	}

	genreIDs := make([]int64, len(movie.Genres))

	for i, d := range movie.Genres {
		if err = tx.NewInsert().Model(&Genre{
			Name: d,
		}).On("CONFLICT (name) DO UPDATE").Set("name = EXCLUDED.name").Returning("id").Scan(ctx, &genreIDs[i]); err != nil {
			return err
		}
	}

	if len(genreIDs) > 0 {
		movieGenres := make([]MovieGenre, len(genreIDs))

		for i, a := range genreIDs {
			movieGenres[i] = MovieGenre{
				MovieID: movie.ID,
				GenreID: a,
			}
		}

		if _, err = tx.NewInsert().Model(&movieGenres).Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (r Repository) Query(ctx context.Context, query benchflix.Query) ([]benchflix.Movie, error) {
	var movies []Movie

	q := r.DB.NewSelect().Model(&movies).
		Relation("Directors", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Order("name ASC")
		}).
		Relation("Actors", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Order("name ASC")
		}).
		Relation("Countries", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Order("name ASC")
		}).
		Relation("Genres", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Order("name ASC")
		}).
		Order("title ASC")

	if query.Search != "" {
		q = q.Where("(EXISTS (?) OR EXISTS (?))",
			r.DB.NewSelect().
				Table("movie_directors").
				ColumnExpr("1").
				Join("JOIN people ON people.id = movie_directors.person_id").
				Where("movie_directors.movie_id = movie.id").
				Where("INSTR(people.name, ?) > 0", query.Search),
			r.DB.NewSelect().
				Table("movie_actors").
				ColumnExpr("1").
				Join("JOIN people ON people.id = movie_actors.person_id").
				Where("movie_actors.movie_id = movie.id").
				Where("INSTR(people.name, ?) > 0", query.Search),
		)
	}

	if query.Genre != "" {
		q = q.Where("EXISTS (?)",
			r.DB.NewSelect().
				TableExpr("movie_genres").
				Join("JOIN genres ON genres.id = movie_genres.genre_id").
				Where("movie_genres.movie_id = movie.id").
				Where("genres.name = ?", query.Genre))
	}

	if query.Country != "" {
		q = q.Where("EXISTS (?)",
			r.DB.NewSelect().
				TableExpr("movie_countries").
				Join("JOIN countries ON countries.id = movie_countries.country_id").
				Where("movie_countries.movie_id = movie.id").
				Where("countries.name = ?", query.Country))
	}

	if !query.AddedBefore.IsZero() {
		q = q.Where("added_at < ?", query.AddedBefore)
	}

	if !query.AddedAfter.IsZero() {
		q = q.Where("added_at > ?", query.AddedAfter)
	}

	if query.MinRating > 0 {
		q = q.Where("rating >= ?", query.MinRating)
	}

	if query.MaxRating > 0 {
		q = q.Where("rating <= ?", query.MaxRating)
	}

	if err := q.Scan(ctx); err != nil {
		return nil, err
	}

	result := make([]benchflix.Movie, len(movies))

	for i, one := range movies {
		movie := benchflix.Movie{
			ID:        one.ID,
			Title:     one.Title,
			AddedAt:   one.AddedAt,
			Rating:    one.Rating,
			Directors: make([]string, len(one.Directors)),
			Actors:    make([]string, len(one.Actors)),
			Countries: make([]string, len(one.Countries)),
			Genres:    make([]string, len(one.Genres)),
		}

		for i, d := range one.Directors {
			movie.Directors[i] = d.Name
		}

		for i, d := range one.Actors {
			movie.Actors[i] = d.Name
		}

		for i, d := range one.Countries {
			movie.Countries[i] = d.Name
		}

		for i, d := range one.Genres {
			movie.Genres[i] = d.Name
		}

		result[i] = movie
	}

	return result, nil
}

func (r Repository) Read(ctx context.Context, id int64) (benchflix.Movie, error) {
	var one Movie

	if err := r.DB.NewSelect().Model(&one).
		Relation("Directors", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Order("name ASC")
		}).
		Relation("Actors", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Order("name ASC")
		}).
		Relation("Countries", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Order("name ASC")
		}).
		Relation("Genres", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Order("name ASC")
		}).
		Where("id = ?", id).Scan(ctx); err != nil {
		return benchflix.Movie{}, err
	}

	movie := benchflix.Movie{
		ID:        one.ID,
		Title:     one.Title,
		AddedAt:   one.AddedAt,
		Rating:    one.Rating,
		Directors: make([]string, len(one.Directors)),
		Actors:    make([]string, len(one.Actors)),
		Countries: make([]string, len(one.Countries)),
		Genres:    make([]string, len(one.Genres)),
	}

	for i, d := range one.Directors {
		movie.Directors[i] = d.Name
	}

	for i, d := range one.Actors {
		movie.Actors[i] = d.Name
	}

	for i, d := range one.Countries {
		movie.Countries[i] = d.Name
	}

	for i, d := range one.Genres {
		movie.Genres[i] = d.Name
	}

	return movie, nil
}
