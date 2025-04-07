package benchflix_test

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	benchflix "github.com/wroge/bench-flix"
	bunflix "github.com/wroge/bench-flix/bun-flix"
	entflix "github.com/wroge/bench-flix/ent-flix"
	gormflix "github.com/wroge/bench-flix/gorm-flix"
	sqlflix "github.com/wroge/bench-flix/sql-flix"
	sqlcflix "github.com/wroge/bench-flix/sqlc-flix"
	sqltflix "github.com/wroge/bench-flix/sqlt-flix"
	sqlxflix "github.com/wroge/bench-flix/sqlx-flix"
)

type Init struct {
	Name string
	New  func() benchflix.Repository
}

var inits = []Init{
	{
		"sql",
		func() benchflix.Repository {
			return sqlflix.NewRepository("sqlite3", ":memory:?_fk=1")
		},
	},
	{
		"gorm",
		func() benchflix.Repository {
			return gormflix.NewRepository(":memory:?_fk=1")
		},
	},
	{
		"sqlt",
		func() benchflix.Repository {
			return sqltflix.NewRepository("sqlite3", ":memory:?_fk=1")
		},
	},
	{
		"ent",
		func() benchflix.Repository {
			return entflix.NewRepository("sqlite3", ":memory:?_fk=1")
		},
	},
	{
		"sqlc",
		func() benchflix.Repository {
			return sqlcflix.NewRepository("sqlite3", ":memory:?_fk=1")
		},
	},
	{
		"bun",
		func() benchflix.Repository {
			return bunflix.NewRepository("sqlite3", ":memory:?_fk=1")
		},
	},
	{
		"sqlx",
		func() benchflix.Repository {
			return sqlxflix.NewRepository("sqlite3", ":memory:?_fk=1")
		},
	},
}

type Case struct {
	Name      string
	Query     benchflix.Query
	ResultLen int
	Result    string
}

type IDCase struct {
	ID     int64
	Result string
}

var (
	queryCases = []Case{
		{
			Name: "Complex",
			Query: benchflix.Query{
				Search:      "Affleck",
				Country:     "United Kingdom",
				Genre:       "Drama",
				AddedAfter:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				AddedBefore: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				MinRating:   4,
				MaxRating:   8,
				Limit:       1,
			},
			ResultLen: 1,
			Result:    `[{505225 The Last Thing He Wanted 2020-02-14 00:00:00 +0000 UTC [Dee Rees] [Anne Hathaway Ben Affleck Edi Gathegi Rosie Perez Willem Dafoe] [United Kingdom United States of America] 4.9 [Drama Thriller]}]`,
		},
		{
			Name: "Simple",
			Query: benchflix.Query{
				MinRating: 9,
				Limit:     1,
			},
			ResultLen: 1,
			Result:    `[{1310741 A Brother and 7 Siblings 2025-01-23 00:00:00 +0000 UTC [Yandy Laurens] [Ahmad Nadif Amanda Rawles Chicco Kurniawan Fatih Unru Freya Jayawardana] [Indonesia] 9.3 [Drama Family]}]`,
		},
		{
			Name: "10",
			Query: benchflix.Query{
				Limit: 10,
			},
			ResultLen: 10,
		},
		{
			Name: "100",
			Query: benchflix.Query{
				Limit: 100,
			},
			ResultLen: 100,
		},
		{
			Name: "1000",
			Query: benchflix.Query{
				Limit: 1000,
			},
			ResultLen: 1000,
		},
	}

	idCases = []IDCase{
		{
			ID:     10192,
			Result: `{10192 Shrek Forever After 2010-05-16 00:00:00 +0000 UTC [Mike Mitchell] [Antonio Banderas Cameron Diaz Eddie Murphy Mike Myers Walt Dohrn] [United States of America] 6.38 [Adventure Animation Comedy Family Fantasy]}`,
		},
	}
)

func BenchmarkSchemaAndCreate(b *testing.B) {
	ctx := context.Background()

	file, err := os.Open("./movies.csv")
	if err != nil {
		b.Fatal(err)
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		b.Fatal(err)
	}

	for _, init := range inits {
		for _, num := range []int{10, 100, 1000} {
			b.Run(fmt.Sprintf("%d_%s", num, init.Name), func(b *testing.B) {
				for b.Loop() {
					r := init.New()

					for _, record := range records[1:1000] {
						movie, err := benchflix.NewMovie(record)
						if err != nil {
							b.Fatal(reflect.TypeOf(r), err)
						}

						if err = r.Create(ctx, movie); err != nil {
							b.Fatal(reflect.TypeOf(r), err)
						}
					}
				}
			})
		}
	}
}

func BenchmarkCreateAndDelete(b *testing.B) {
	ctx := context.Background()

	file, err := os.Open("./movies.csv")
	if err != nil {
		b.Fatal(err)
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		b.Fatal(err)
	}

	do := func(r benchflix.Repository, num int) {
		ids := []int64{}

		for _, record := range records[1:num] {
			movie, err := benchflix.NewMovie(record)
			if err != nil {
				b.Fatal(reflect.TypeOf(r), err)
			}

			if err = r.Create(ctx, movie); err != nil {
				b.Fatal(reflect.TypeOf(r), err)
			}

			ids = append(ids, movie.ID)
		}

		for _, id := range ids {
			if err = r.Delete(ctx, id); err != nil {
				b.Fatal(reflect.TypeOf(r), err)
			}
		}
	}

	for _, init := range inits {
		for _, num := range []int{10, 100, 1000} {
			b.Run(fmt.Sprintf("%d_%s", num, init.Name), func(b *testing.B) {
				r := init.New()

				// Warmup
				do(r, num)

				for b.Loop() {
					do(r, num)
				}
			})
		}
	}
}

func Test_Query(t *testing.T) {
	ctx := context.Background()

	file, err := os.Open("./movies.csv")
	if err != nil {
		panic(err)
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		panic(err)
	}

	for _, c := range queryCases {
		for _, init := range inits {
			r := init.New()

			t.Run(c.Name+"_"+init.Name, func(t *testing.T) {
				for _, record := range records[1:] {
					movie, err := benchflix.NewMovie(record)
					if err != nil {
						t.Fatal(reflect.TypeOf(r), err)
					}

					if err = r.Create(ctx, movie); err != nil {
						t.Fatal(reflect.TypeOf(r), err)
					}
				}

				movies, err := r.Query(ctx, c.Query)
				if err != nil {
					t.Fatal(reflect.TypeOf(r), err)
				}

				if c.ResultLen != len(movies) {
					t.Fatalf("%s: %v: invalid number of movies: want %d got %d",
						reflect.TypeOf(r), c.Query, c.ResultLen, len(movies))
				}

				if c.Result != "" && fmt.Sprint(movies) != c.Result {
					t.Fatal(reflect.TypeOf(r), c.Query, movies)
				}
			})
		}
	}
}

func BenchmarkQuery(b *testing.B) {
	ctx := context.Background()

	file, err := os.Open("./movies.csv")
	if err != nil {
		b.Fatal(err)
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		b.Fatal(err)
	}

	do := func(r benchflix.Repository, c Case) {
		movies, err := r.Query(ctx, c.Query)
		if err != nil {
			b.Fatal(reflect.TypeOf(r), err)
		}

		if c.ResultLen != len(movies) {
			b.Fatalf("%s: %v: invalid number of movies: want %d got %d",
				reflect.TypeOf(r), c.Query, c.ResultLen, len(movies))
		}

		if c.Result != "" && fmt.Sprint(movies) != c.Result {
			b.Fatal(reflect.TypeOf(r), c.Query, movies)
		}
	}

	for _, c := range queryCases {
		for _, init := range inits {
			r := init.New()

			for _, record := range records[1:] {
				movie, err := benchflix.NewMovie(record)
				if err != nil {
					b.Fatal(reflect.TypeOf(r), err)
				}

				if err = r.Create(ctx, movie); err != nil {
					b.Fatal(reflect.TypeOf(r), err)
				}
			}

			// Warmup
			do(r, c)

			b.Run(c.Name+"_"+init.Name, func(b *testing.B) {
				for b.Loop() {
					do(r, c)
				}
			})
		}
	}
}

func Test_Read(t *testing.T) {
	ctx := context.Background()

	file, err := os.Open("./movies.csv")
	if err != nil {
		t.Fatal(err)
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range idCases {
		for _, init := range inits {
			r := init.New()

			t.Run(init.Name, func(t *testing.T) {
				for _, record := range records[1:] {
					movie, err := benchflix.NewMovie(record)
					if err != nil {
						t.Fatal(reflect.TypeOf(r), err)
					}

					if err = r.Create(ctx, movie); err != nil {
						t.Fatal(reflect.TypeOf(r), err)
					}
				}

				movie, err := r.Read(ctx, c.ID)
				if err != nil {
					t.Fatal(reflect.TypeOf(r), err)
				}

				if fmt.Sprint(movie) != c.Result {
					t.Fatal(reflect.TypeOf(r), movie)
				}
			})
		}
	}
}

func BenchmarkRead(b *testing.B) {
	ctx := context.Background()

	file, err := os.Open("./movies.csv")
	if err != nil {
		b.Fatal(err)
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		b.Fatal(err)
	}

	do := func(r benchflix.Repository, c IDCase) {
		movie, err := r.Read(ctx, c.ID)
		if err != nil {
			b.Fatal(reflect.TypeOf(r), err)
		}

		if fmt.Sprint(movie) != c.Result {
			b.Fatal(reflect.TypeOf(r), movie)
		}
	}

	for _, c := range idCases {
		for _, init := range inits {
			r := init.New()

			for _, record := range records[1:] {
				movie, err := benchflix.NewMovie(record)
				if err != nil {
					b.Fatal(err)
				}

				if err = r.Create(ctx, movie); err != nil {
					b.Fatal(err)
				}
			}

			// Warmup
			do(r, c)

			b.Run(init.Name, func(b *testing.B) {
				for b.Loop() {
					do(r, c)
				}
			})
		}
	}
}
