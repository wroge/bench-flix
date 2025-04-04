package benchflix

import (
	"context"
	"strconv"
	"strings"
	"time"
)

type Movie struct {
	ID        int64 `xorm:"id"`
	Title     string
	AddedAt   time.Time
	Directors []string
	Actors    []string
	Countries []string
	Rating    float64
	Genres    []string
}

type Query struct {
	Search                  string
	Genre                   string
	Country                 string
	AddedBefore, AddedAfter time.Time
	MinRating, MaxRating    float64
}

type Repository interface {
	Create(ctx context.Context, movie Movie) error
	Read(ctx context.Context, id int64) (Movie, error)
	Query(ctx context.Context, query Query) ([]Movie, error)
	Delete(ctx context.Context, id int64) error
}

func NewMovie(record []string) (Movie, error) {
	id, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return Movie{}, err
	}

	added, err := time.Parse(time.DateOnly, record[6])
	if err != nil {
		return Movie{}, err
	}

	rating, err := strconv.ParseFloat(record[8], 64)
	if err != nil {
		return Movie{}, err
	}

	return Movie{
		ID:        id,
		Title:     record[2],
		AddedAt:   added,
		Directors: Unique(strings.Split(record[3], ", ")),
		Actors:    Unique(strings.Split(record[4], ", ")),
		Countries: Unique(strings.Split(record[5], ", ")),
		Rating:    rating,
		Genres:    Unique(strings.Split(record[10], ", ")),
	}, nil
}

func Unique(list []string) []string {
	var (
		seen   = map[string]bool{}
		result []string
	)

	for _, s := range list {
		if s == "" || seen[s] {
			continue
		}

		seen[s] = true
		result = append(result, s)
	}

	return result
}
