package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	benchflix "github.com/wroge/bench-flix"
	entflix "github.com/wroge/bench-flix/ent-flix"
)

func main() {
	ctx := context.Background()
	r := entflix.NewRepository("sqlite3", ":memory:?_fk=1")

	file, err := os.Open("./movies.csv")
	if err != nil {
		panic(err)
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		panic(err)
	}

	for _, record := range records[1:] {
		movie, err := benchflix.NewMovie(record)
		if err != nil {
			panic(err)
		}

		if err = r.Create(ctx, movie); err != nil {
			panic(err)
		}
	}

	movies, err := r.Query(ctx, benchflix.Query{
		MinRating: 5,
		Limit:     1,
		// Search:      "Affleck",
		// Country:     "United Kingdom",
		// Genre:       "Drama",
		// AddedAfter:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		// AddedBefore: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		// MinRating:   4,
		// MaxRating:   8,
		// Limit:       1,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(movies)

	fmt.Println(r.Read(ctx, 1310741))
}
