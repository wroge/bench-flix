package gormflix

import (
	"context"
	"time"

	benchflix "github.com/wroge/bench-flix"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Movie struct {
	ID        int64      `gorm:"primaryKey"`
	Title     string     `gorm:"not null"`
	AddedAt   time.Time  `gorm:"not null"`
	Rating    float64    `gorm:"not null"`
	Directors []*Person  `gorm:"many2many:movie_directors"`
	Actors    []*Person  `gorm:"many2many:movie_actors"`
	Countries []*Country `gorm:"many2many:movie_countries"`
	Genres    []*Genre   `gorm:"many2many:movie_genres"`
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

func NewRepository() benchflix.Repository {
	db, err := gorm.Open(sqlite.Open(":memory:?_fk=1"), &gorm.Config{
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

// Create implements benchflix.Repository.
func (r Repository) Create(ctx context.Context, movie benchflix.Movie) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		create := Movie{
			ID:      movie.ID,
			Title:   movie.Title,
			AddedAt: movie.AddedAt,
			Rating:  movie.Rating,
		}

		if err := tx.Create(&create).Error; err != nil {
			return err
		}

		for _, name := range movie.Directors {
			var person Person
			if err := tx.FirstOrCreate(&person, Person{Name: name}).Error; err != nil {
				return err
			}
			if err := tx.Model(&create).Association("Directors").Append(&person); err != nil {
				return err
			}
		}

		for _, name := range movie.Actors {
			var person Person
			if err := tx.FirstOrCreate(&person, Person{Name: name}).Error; err != nil {
				return err
			}
			if err := tx.Model(&create).Association("Actors").Append(&person); err != nil {
				return err
			}
		}

		for _, name := range movie.Countries {
			var country Country
			if err := tx.FirstOrCreate(&country, Country{Name: name}).Error; err != nil {
				return err
			}
			if err := tx.Model(&create).Association("Countries").Append(&country); err != nil {
				return err
			}
		}

		for _, name := range movie.Genres {
			var genre Genre
			if err := tx.FirstOrCreate(&genre, Genre{Name: name}).Error; err != nil {
				return err
			}
			if err := tx.Model(&create).Association("Genres").Append(&genre); err != nil {
				return err
			}
		}

		return nil
	})
}

// Query implements benchflix.Repository.
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

		movies[i] = movie
	}

	return movies, nil
}

// Read implements benchflix.Repository.
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
