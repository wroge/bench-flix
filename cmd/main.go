package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	benchflix "github.com/wroge/bench-flix"
	bunflix "github.com/wroge/bench-flix/bun-flix"
)

func main() {
	ctx := context.Background()
	r := bunflix.NewRepository()

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
		MinRating: 9.5,
		MaxRating: 10,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(movies)
}
