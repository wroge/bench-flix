package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	benchflix "github.com/wroge/bench-flix"
	sqlxflix "github.com/wroge/bench-flix/sqlx-flix"
)

func main() {
	ctx := context.Background()
	r := sqlxflix.NewRepository("sqlite3", ":memory:?_fk=1")

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
		Limit: 10,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(movies)
}
