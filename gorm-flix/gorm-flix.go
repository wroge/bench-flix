package gormflix

import (
	"context"
	"time"

	benchflix "github.com/wroge/bench-flix"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Movie struct {
	ID        int64      `gorm:"primaryKey"`
	Title     string     `gorm:"not null"`
	AddedAt   time.Time  `gorm:"not null"`
	Rating    float64    `gorm:"not null"`
	Directors []*Person  `gorm:"many2many:movie_directors;constraint:OnDelete:CASCADE"`
	Actors    []*Person  `gorm:"many2many:movie_actors;constraint:OnDelete:CASCADE"`
	Countries []*Country `gorm:"many2many:movie_countries;constraint:OnDelete:CASCADE"`
	Genres    []*Genre   `gorm:"many2many:movie_genres;constraint:OnDelete:CASCADE"`
}

type Person struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

type Country struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

type Genre struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

type MovieDirector struct {
	MovieID  int64 `gorm:"primaryKey;not null"`
	PersonID int64 `gorm:"primaryKey;not null"`
}

type MovieActor struct {
	MovieID  int64 `gorm:"primaryKey;not null"`
	PersonID int64 `gorm:"primaryKey;not null"`
}

type MovieCountry struct {
	MovieID   int64 `gorm:"primaryKey;not null"`
	CountryID int64 `gorm:"primaryKey;not null"`
}

type MovieGenre struct {
	MovieID int64 `gorm:"primaryKey;not null"`
	GenreID int64 `gorm:"primaryKey;not null"`
}

func NewRepository(dsn string) benchflix.Repository {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(
		&Movie{}, &Person{}, &Country{}, &Genre{},
		&MovieDirector{}, &MovieActor{}, &MovieCountry{}, &MovieGenre{}); err != nil {
		panic(err)
	}

	return Repository{
		DB: db,
	}
}

type Repository struct {
	DB *gorm.DB
}

func (r Repository) Delete(ctx context.Context, id int64) error {
	return r.DB.Delete(Movie{ID: id}).Error
}

func (r Repository) Create(ctx context.Context, movie benchflix.Movie) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		create := Movie{
			ID:      movie.ID,
			Title:   movie.Title,
			AddedAt: movie.AddedAt,
			Rating:  movie.Rating,
		}

		if len(movie.Directors) > 0 {
			create.Directors = make([]*Person, len(movie.Directors))

			for i, name := range movie.Directors {
				create.Directors[i] = &Person{
					Name: name,
				}
			}

			if err := tx.Clauses(clause.OnConflict{
				DoUpdates: clause.Set{
					clause.Assignment{
						Column: clause.Column{Name: "name"},
						Value:  gorm.Expr("EXCLUDED.name"),
					},
				},
			}).Create(create.Directors).Error; err != nil {
				return err
			}
		}

		if len(movie.Actors) > 0 {
			create.Actors = make([]*Person, len(movie.Actors))

			for i, name := range movie.Actors {
				create.Actors[i] = &Person{
					Name: name,
				}
			}

			if err := tx.Clauses(clause.OnConflict{
				DoUpdates: clause.Set{
					clause.Assignment{
						Column: clause.Column{Name: "name"},
						Value:  gorm.Expr("EXCLUDED.name"),
					},
				},
			}).Create(create.Actors).Error; err != nil {
				return err
			}
		}

		if len(movie.Countries) > 0 {
			create.Countries = make([]*Country, len(movie.Countries))

			for i, name := range movie.Countries {
				create.Countries[i] = &Country{
					Name: name,
				}
			}

			if err := tx.Clauses(clause.OnConflict{
				DoUpdates: clause.Set{
					clause.Assignment{
						Column: clause.Column{Name: "name"},
						Value:  gorm.Expr("EXCLUDED.name"),
					},
				},
			}).Create(create.Countries).Error; err != nil {
				return err
			}
		}

		if len(movie.Genres) > 0 {
			create.Genres = make([]*Genre, len(movie.Genres))

			for i, name := range movie.Genres {
				create.Genres[i] = &Genre{
					Name: name,
				}
			}

			if err := tx.Clauses(clause.OnConflict{
				DoUpdates: clause.Set{
					clause.Assignment{
						Column: clause.Column{Name: "name"},
						Value:  gorm.Expr("EXCLUDED.name"),
					},
				},
			}).Create(create.Genres).Error; err != nil {
				return err
			}
		}

		return tx.Create(&create).Error
	})
}

func (r Repository) Query(ctx context.Context, query benchflix.Query) ([]benchflix.Movie, error) {
	var list []Movie

	db := r.DB.WithContext(ctx).
		Preload("Directors", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		}).
		Preload("Actors", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		}).
		Preload("Countries", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		}).
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		}).
		Distinct("movies.*").
		Order("movies.title ASC")

	if query.Limit > 0 {
		db = db.Limit(int(query.Limit))
	}

	if query.Search != "" {
		db = db.Joins("JOIN movie_directors md ON md.movie_id = movies.id").
			Joins("JOIN people d ON d.id = md.person_id").
			Joins("JOIN movie_actors ma ON ma.movie_id = movies.id").
			Joins("JOIN people a ON a.id = ma.person_id").
			Where("INSTR(d.name, ?) > 0 OR INSTR(a.name, ?) > 0", query.Search, query.Search)
	}

	if query.Genre != "" {
		db = db.Joins("JOIN movie_genres mg ON mg.movie_id = movies.id").
			Joins("JOIN genres g ON g.id = mg.genre_id").
			Where("g.name = ?", query.Genre)
	}

	if query.Country != "" {
		db = db.Joins("JOIN movie_countries mc ON mc.movie_id = movies.id").
			Joins("JOIN countries c ON c.id = mc.country_id").
			Where("c.name = ?", query.Country)
	}

	if !query.AddedBefore.IsZero() {
		db = db.Where("added_at < ?", query.AddedBefore)
	}

	if !query.AddedAfter.IsZero() {
		db = db.Where("added_at > ?", query.AddedAfter)
	}

	if query.MinRating > 0 {
		db = db.Where("rating >= ?", query.MinRating)
	}

	if query.MaxRating > 0 {
		db = db.Where("rating <= ?", query.MaxRating)
	}

	err := db.Find(&list).Error
	if err != nil {
		return nil, err
	}

	movies := make([]benchflix.Movie, len(list))

	for i, one := range list {
		movies[i] = ConvertMovie(one)
	}

	return movies, nil
}

func (r Repository) Read(ctx context.Context, id int64) (benchflix.Movie, error) {
	var one Movie

	err := r.DB.WithContext(ctx).
		Preload("Directors", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		}).
		Preload("Actors", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		}).
		Preload("Countries", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		}).
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		}).
		Order("movies.title ASC").
		First(&one, "id = ?", id).Error
	if err != nil {
		return benchflix.Movie{}, err
	}

	return ConvertMovie(one), nil
}

func ConvertMovie(m Movie) benchflix.Movie {
	movie := benchflix.Movie{
		ID:        m.ID,
		Title:     m.Title,
		AddedAt:   m.AddedAt,
		Rating:    m.Rating,
		Directors: make([]string, len(m.Directors)),
		Actors:    make([]string, len(m.Actors)),
		Countries: make([]string, len(m.Countries)),
		Genres:    make([]string, len(m.Genres)),
	}

	for i, d := range m.Directors {
		movie.Directors[i] = d.Name
	}

	for i, d := range m.Actors {
		movie.Actors[i] = d.Name
	}

	for i, d := range m.Countries {
		movie.Countries[i] = d.Name
	}

	for i, d := range m.Genres {
		movie.Genres[i] = d.Name
	}

	return movie
}
