package xormflix

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	benchflix "github.com/wroge/bench-flix"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type Movie struct {
	ID      int64     `xorm:"id pk"`
	Title   string    `xorm:"notnull"`
	AddedAt time.Time `xorm:"added_at DATE notnull"`
	Rating  float64   `xorm:"notnull"`
}

func (Movie) TableName() string {
	return "movies"
}

type Person struct {
	ID   int64  `xorm:"id pk autoincr"`
	Name string `xorm:"not null unique"`
}

func (Person) TableName() string {
	return "people"
}

type Country struct {
	ID   int64  `xorm:"id pk autoincr"`
	Name string `xorm:"notnull unique"`
}

func (Country) TableName() string {
	return "countries"
}

type Genre struct {
	ID   int64  `xorm:"id pk autoincr"`
	Name string `xorm:"notnull unique"`
}

func (Genre) TableName() string {
	return "genres"
}

type MovieDirector struct {
	MovieID  int64 `xorm:"movie_id pk notnull"`
	PersonID int64 `xorm:"person_id pk notnull"`
}

func (MovieDirector) TableName() string {
	return "movie_directors"
}

type MovieActor struct {
	MovieID  int64 `xorm:"movie_id pk notnull"`
	PersonID int64 `xorm:"person_id pk notnull"`
}

func (MovieActor) TableName() string {
	return "movie_actors"
}

type MovieCountry struct {
	MovieID   int64 `xorm:"movie_id pk notnull"`
	CountryID int64 `xorm:"country_id pk notnull"`
}

func (MovieCountry) TableName() string {
	return "movie_countries"
}

type MovieGenre struct {
	MovieID int64 `xorm:"movie_id pk notnull"`
	GenreID int64 `xorm:"genre_id pk notnull"`
}

func (MovieGenre) TableName() string {
	return "movie_genres"
}

func NewRepository() benchflix.Repository {
	time.Local = time.UTC

	engine, err := xorm.NewEngine("sqlite3", ":memory:?_fk=1")
	if err != nil {
		panic(err)
	}

	if err = engine.CreateTables(
		new(Movie), new(Person), new(Country), new(Genre),
		new(MovieDirector), new(MovieActor), new(MovieCountry), new(MovieGenre)); err != nil {
		panic(err)
	}

	return Repository{
		Engine: engine,
	}
}

type Repository struct {
	Engine *xorm.Engine
}

func (r Repository) Create(ctx context.Context, movie benchflix.Movie) error {
	_, err := r.Engine.Transaction(func(s *xorm.Session) (interface{}, error) {
		if _, err := s.Insert(&Movie{
			ID:      movie.ID,
			Title:   movie.Title,
			AddedAt: movie.AddedAt,
			Rating:  movie.Rating,
		}); err != nil {
			return nil, err
		}

		personMap := make(map[string]int64)
		for _, name := range append(movie.Directors, movie.Actors...) {
			if _, ok := personMap[name]; ok {
				continue
			}

			var person Person
			exists, err := s.Where("name = ?", name).Get(&person)
			if err != nil {
				return nil, err
			}

			if !exists {
				person = Person{Name: name}
				if _, err := s.Insert(&person); err != nil {
					return nil, err
				}
			}

			personMap[name] = person.ID
		}

		for _, name := range movie.Directors {
			if _, err := s.Insert(&MovieDirector{
				MovieID:  movie.ID,
				PersonID: personMap[name],
			}); err != nil {
				return nil, err
			}
		}

		for _, name := range movie.Actors {
			if _, err := s.Insert(&MovieActor{
				MovieID:  movie.ID,
				PersonID: personMap[name],
			}); err != nil {
				return nil, err
			}
		}

		countryMap := make(map[string]int64)
		for _, name := range movie.Countries {
			if _, ok := countryMap[name]; ok {
				continue
			}

			var country Country
			exists, err := s.Where("name = ?", name).Get(&country)
			if err != nil {
				return nil, err
			}

			if !exists {
				country = Country{Name: name}
				if _, err := s.Insert(&country); err != nil {
					return nil, err
				}
			}

			countryMap[name] = country.ID

			if _, err := s.Insert(&MovieCountry{
				MovieID:   movie.ID,
				CountryID: country.ID,
			}); err != nil {
				return nil, err
			}
		}

		genreMap := make(map[string]int64)
		for _, name := range movie.Genres {
			if _, ok := genreMap[name]; ok {
				continue
			}

			var genre Genre
			exists, err := s.Where("name = ?", name).Get(&genre)
			if err != nil {
				return nil, err
			}

			if !exists {
				genre = Genre{Name: name}
				if _, err := s.Insert(&genre); err != nil {
					return nil, err
				}
			}

			genreMap[name] = genre.ID

			if _, err := s.Insert(&MovieGenre{
				MovieID: movie.ID,
				GenreID: genre.ID,
			}); err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return err
}

func (r Repository) Query(ctx context.Context, query benchflix.Query) ([]benchflix.Movie, error) {
	q := r.Engine.Context(ctx).Table("movies").Select(`
		movies.id,
		movies.title,
		movies.added_at,
		movies.rating,
		(
			SELECT JSON_GROUP_ARRAY(people.name ORDER BY people.name)
			FROM movie_directors
			JOIN people ON people.id = movie_directors.person_id
			WHERE movie_directors.movie_id = movies.id
		) AS directors,
		(
			SELECT JSON_GROUP_ARRAY(people.name ORDER BY people.name)
			FROM movie_actors
			JOIN people ON people.id = movie_actors.person_id
			WHERE movie_actors.movie_id = movies.id
		) AS actors,
		(
			SELECT JSON_GROUP_ARRAY(countries.name ORDER BY countries.name)
			FROM movie_countries
			JOIN countries ON countries.id = movie_countries.country_id
			WHERE movie_countries.movie_id = movies.id
		) AS countries,
		(
			SELECT JSON_GROUP_ARRAY(genres.name ORDER BY genres.name)
			FROM movie_genres
			JOIN genres ON genres.id = movie_genres.genre_id
			WHERE movie_genres.movie_id = movies.id
		) AS genres
	`).OrderBy("movies.title ASC")

	if query.Search != "" {
		q.Where(`(
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
		)`, query.Search, query.Search)
	}

	if query.Genre != "" {
		q.Where(`EXISTS (
			SELECT 1 FROM movie_genres
			JOIN genres ON genres.id = movie_genres.genre_id
			WHERE movie_genres.movie_id = movies.id AND genres.name = ?
		)`, query.Genre)
	}

	if query.Country != "" {
		q.Where(`EXISTS (
			SELECT 1 FROM movie_countries
			JOIN countries ON countries.id = movie_countries.country_id
			WHERE movie_countries.movie_id = movies.id AND countries.name = ?
		)`, query.Country)
	}

	if !query.AddedBefore.IsZero() {
		q.Where("added_at < ?", query.AddedBefore)
	}

	if !query.AddedAfter.IsZero() {
		q.Where("added_at > ?", query.AddedAfter)
	}

	if query.MinRating > 0 {
		q.Where("rating >= ?", query.MinRating)
	}

	if query.MaxRating > 0 {
		q.Where("rating <= ?", query.MaxRating)
	}

	var movies []benchflix.Movie

	if err := q.Find(&movies); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r Repository) Read(ctx context.Context, id int64) (benchflix.Movie, error) {
	var movie benchflix.Movie

	ok, err := r.Engine.Context(ctx).Table("movies").Select(`
		movies.id,
		movies.title,
		movies.added_at,
		movies.rating,
		(
			SELECT json_group_array(people.name ORDER BY people.name)
			FROM movie_directors
			JOIN people ON people.id = movie_directors.person_id
			WHERE movie_directors.movie_id = movies.id
		) AS directors,
		(
			SELECT json_group_array(people.name ORDER BY people.name)
			FROM movie_actors
			JOIN people ON people.id = movie_actors.person_id
			WHERE movie_actors.movie_id = movies.id
		) AS actors,
		(
			SELECT json_group_array(countries.name ORDER BY countries.name)
			FROM movie_countries
			JOIN countries ON countries.id = movie_countries.country_id
			WHERE movie_countries.movie_id = movies.id
		) AS countries,
		(
			SELECT json_group_array(genres.name ORDER BY genres.name)
			FROM movie_genres
			JOIN genres ON genres.id = movie_genres.genre_id
			WHERE movie_genres.movie_id = movies.id
		) AS genres
	`).Where(builder.Eq{"id": id}).OrderBy("movies.title ASC").Get(&movie)
	if err != nil {
		return benchflix.Movie{}, err
	}

	if !ok {
		return benchflix.Movie{}, sql.ErrNoRows
	}

	return movie, nil
}
