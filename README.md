# Bench-Flix

This benchmark imports a dataset of Netflix movies into a SQLite database and runs a range of queries to compare performance, memory usage, and allocation efficiency across different go packages.

⚠️ Results aren’t always perfectly comparable — for example, GORM uses preloading to handle many-to-many relationships.

I’m open to feedback and suggestions — I’m not an expert in every tool and aim to make this benchmark as fair and informative as possible.

👉 Want to add another SQL library? Just open a pull request!

- Dataset: [kaggle/netflix-movies](https://www.kaggle.com/datasets/bhargavchirumamilla/netflix-movies-and-tv-shows-till-2025)
- Sqlite Driver: [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
- sql: database/sql
- gorm: [gorm.io](https://gorm.io/)
- ent: [entgo.io](https://entgo.io/)
- sqlc: [sqlc.dev](https://sqlc.dev/)
- bun: [bun.uptrace.dev](https://bun.uptrace.dev/)
- sqlt: [wroge/sqlt](https://github.com/wroge/sqlt) (my own package)

## Benchmark Result

```bash
go test -bench . -run=xxx -benchmem ./bench-flix_test.go -benchtime=10s
goos: darwin
goarch: arm64
cpu: Apple M3 Pro
BenchmarkQuery/Complex_sql-12                572          20970552 ns/op           14156 B/op        271 allocs/op
BenchmarkQuery/Complex_gorm-12              1754           6980292 ns/op          124114 B/op       2439 allocs/op
BenchmarkQuery/Complex_sqlt-12               570          21409093 ns/op           13576 B/op        308 allocs/op
BenchmarkQuery/Complex_ent-12                567          21183904 ns/op           83963 B/op       1914 allocs/op
BenchmarkQuery/Complex_sqlc-12               397          30078865 ns/op           13122 B/op        250 allocs/op
BenchmarkQuery/Complex_bun-12                534          22275331 ns/op           56118 B/op       1126 allocs/op
BenchmarkQuery/Mid_sql-12                   3570           3374104 ns/op           10299 B/op        219 allocs/op
BenchmarkQuery/Mid_gorm-12                  2119           5719571 ns/op          104307 B/op       2022 allocs/op
BenchmarkQuery/Mid_sqlt-12                  3532           3363826 ns/op           10074 B/op        252 allocs/op
BenchmarkQuery/Mid_ent-12                    686          17507475 ns/op           67531 B/op       1540 allocs/op
BenchmarkQuery/Mid_sqlc-12                  2697           4419772 ns/op            9187 B/op        201 allocs/op
BenchmarkQuery/Mid_bun-12                   3414           3487511 ns/op           49158 B/op        898 allocs/op
BenchmarkQuery/Simple_sql-12               17971            667279 ns/op           79206 B/op       1677 allocs/op
BenchmarkQuery/Simple_gorm-12               9457           1299788 ns/op          604420 B/op      12200 allocs/op
BenchmarkQuery/Simple_sqlt-12              17614            681456 ns/op           85002 B/op       1863 allocs/op
BenchmarkQuery/Simple_ent-12               13209            907408 ns/op          313598 B/op       6698 allocs/op
BenchmarkQuery/Simple_sqlc-12              13891            863933 ns/op           89572 B/op       1513 allocs/op
BenchmarkQuery/Simple_bun-12               13538            885133 ns/op          200077 B/op       5928 allocs/op
BenchmarkRead/sql-12                      472965             25295 ns/op            2352 B/op         69 allocs/op
BenchmarkRead/gorm-12                     135396             88236 ns/op           60015 B/op       1004 allocs/op
BenchmarkRead/sqlt-12                     448791             26543 ns/op            3544 B/op         95 allocs/op
BenchmarkRead/ent-12                      203052             58833 ns/op           33617 B/op        848 allocs/op
BenchmarkRead/sqlc-12                     377048             31611 ns/op            2296 B/op         67 allocs/op
BenchmarkRead/bun-12                      250628             47776 ns/op           36537 B/op        414 allocs/op
PASS
ok      command-line-arguments  363.736s
```

## Benchmark Test

```go
var (
	queryCases = []Case{
		{
			Name: "Complex",
			Query: benchflix.Query{
				Search:  "Affleck",
				Country: "United Kingdom",
				Genre:   "Drama",
			},
			Result: `...`,
		},
		{
			Name: "Mid",
			Query: benchflix.Query{
				Search:     "Affleck",
				AddedAfter: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			Result: `...`,
		},
		{
			Name: "Simple",
			Query: benchflix.Query{
				MinRating: 9.5,
				MaxRating: 10,
			},
			Result: `...`,
		},
	}

	readCases = []IDCase{
		{
			ID:     10192,
			Result: `...`,
		},
	}
)

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

	for _, c := range queryCases {
		for _, init := range repositories {
			r := init()

			for _, record := range records[1:] {
				movie, err := benchflix.NewMovie(record)
				if err != nil {
					b.Fatal(reflect.TypeOf(r), err)
				}

				if err = r.Create(ctx, movie); err != nil {
					b.Fatal(reflect.TypeOf(r), err)
				}
			}

			b.Run(c.Name+"_"+strings.TrimSuffix(reflect.TypeOf(r).String(), "flix.Repository"), func(b *testing.B) {
				for b.Loop() {
					movies, err := r.Query(ctx, c.Query)
					if err != nil {
						b.Fatal(reflect.TypeOf(r), err)
					}

					if fmt.Sprint(movies) != c.Result {
						b.Fatal(reflect.TypeOf(r), movies)
					}
				}
			})
		}
	}
}
```