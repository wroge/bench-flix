# Bench-Flix

This benchmark imports a dataset of Netflix movies into a SQLite database and runs a range of queries to compare performance, memory usage, and allocation efficiency across different go packages.

⚠️ Results aren’t always perfectly comparable — for example, both GORM and Bun use preloading to resolve many-to-many relationships. 

I’m open to feedback and suggestions — I’m not an expert in every tool and aim to make this benchmark as fair and informative as possible.

👉 Want to add another SQL library? Just open a pull request!

- Dataset: [kaggle/netflix-movies](https://www.kaggle.com/datasets/bhargavchirumamilla/netflix-movies-and-tv-shows-till-2025)
- Sqlite Driver: [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
- sql: database/sql
- gorm: [gorm.io](https://gorm.io/)
- ent: [entgo.io](https://entgo.io/)
- sqlc: [sqlc.dev](https://sqlc.dev/)
- bun: [bun.uptrace.dev](https://bun.uptrace.dev/)
- xorm: [xorm.io](https://xorm.io/)
- sqlt: [wroge/sqlt](https://github.com/wroge/sqlt) (my own package)

## Benchmark

The “Complex” query in the ```gorm``` repository is significantly faster than in other implementations. This suggests that ```gorm```'s preloading strategy performs better for handling multiple many-to-many relationships compared to joining everything in a single query.
As expected, the implementation using standard SQL is the fastest overall.
```sqlt``` (my own library) is competitive with standard SQL, aiming for clean abstraction with minimal runtime overhead. ```sqlc``` is fast and efficient in simpler queries but struggles in complex multi-table lookups.

```bash
go test -bench . -run=xxx -benchmem -benchtime=10s
goos: darwin
goarch: arm64
pkg: github.com/wroge/bench-flix
cpu: Apple M3 Pro
BenchmarkQuery/Complex_sql-12                529          22509632 ns/op           14157 B/op        271 allocs/op
BenchmarkQuery/Complex_gorm-12              1623           7354243 ns/op          124140 B/op       2439 allocs/op
BenchmarkQuery/Complex_sqlt-12               541          22480862 ns/op           13580 B/op        308 allocs/op
BenchmarkQuery/Complex_ent-12                529          22549345 ns/op           83966 B/op       1914 allocs/op
BenchmarkQuery/Complex_sqlc-12               387          30614424 ns/op           13123 B/op        250 allocs/op
BenchmarkQuery/Complex_bun-12                524          23105291 ns/op           56119 B/op       1126 allocs/op
BenchmarkQuery/Complex_xorm-12               537          22714580 ns/op           37740 B/op        786 allocs/op
BenchmarkQuery/Mid_sql-12                   3499           3521276 ns/op           10300 B/op        219 allocs/op
BenchmarkQuery/Mid_gorm-12                  1959           6042494 ns/op          104302 B/op       2022 allocs/op
BenchmarkQuery/Mid_sqlt-12                  3504           3519769 ns/op           10074 B/op        252 allocs/op
BenchmarkQuery/Mid_ent-12                    657          18263914 ns/op           67532 B/op       1540 allocs/op
BenchmarkQuery/Mid_sqlc-12                  2692           4495674 ns/op            9187 B/op        201 allocs/op
BenchmarkQuery/Mid_bun-12                   3400           3548343 ns/op           49157 B/op        898 allocs/op
BenchmarkQuery/Mid_xorm-12                  3506           3424511 ns/op           27910 B/op        627 allocs/op
BenchmarkQuery/Simple_sql-12               17860            672030 ns/op           79193 B/op       1677 allocs/op
BenchmarkQuery/Simple_gorm-12               9163           1342026 ns/op          604424 B/op      12200 allocs/op
BenchmarkQuery/Simple_sqlt-12              16878            701401 ns/op           85016 B/op       1863 allocs/op
BenchmarkQuery/Simple_ent-12               12912            930565 ns/op          313551 B/op       6698 allocs/op
BenchmarkQuery/Simple_sqlc-12              13556            881140 ns/op           89547 B/op       1513 allocs/op
BenchmarkQuery/Simple_bun-12               13387            894804 ns/op          200051 B/op       5928 allocs/op
BenchmarkQuery/Simple_xorm-12              14360            835659 ns/op          205922 B/op       5035 allocs/op
BenchmarkRead/sql-12                      468254             25635 ns/op            2352 B/op         69 allocs/op
BenchmarkRead/gorm-12                     133970             89610 ns/op           60015 B/op       1004 allocs/op
BenchmarkRead/sqlt-12                     450860             27081 ns/op            3544 B/op         95 allocs/op
BenchmarkRead/ent-12                      205086             59128 ns/op           33617 B/op        848 allocs/op
BenchmarkRead/sqlc-12                     381436             32063 ns/op            2296 B/op         67 allocs/op
BenchmarkRead/bun-12                      249048             47880 ns/op           36537 B/op        414 allocs/op
BenchmarkRead/xorm-12                     335374             36718 ns/op           10920 B/op        261 allocs/op
PASS
ok      github.com/wroge/bench-flix     631.165s
```

## Benchstat

```bash
go install golang.org/x/perf/cmd/benchstat@latest
go test -bench . -run=xxx -benchmem -count=10 > bench.out
benchstat bench.out
goos: darwin
goarch: arm64
pkg: github.com/wroge/bench-flix
cpu: Apple M3 Pro
                      │  bench.out   │
                      │    sec/op    │
Query/Complex_sql-12    22.25m ±  1%
Query/Complex_gorm-12   7.356m ±  0%
Query/Complex_sqlt-12   22.44m ±  3%
Query/Complex_ent-12    22.88m ±  1%
Query/Complex_sqlc-12   31.22m ±  1%
Query/Complex_bun-12    23.23m ±  3%
Query/Complex_xorm-12   22.26m ±  2%
Query/Mid_sql-12        3.475m ±  4%
Query/Mid_gorm-12       6.051m ±  4%
Query/Mid_sqlt-12       3.451m ±  1%
Query/Mid_ent-12        18.31m ±  5%
Query/Mid_sqlc-12       4.542m ±  3%
Query/Mid_bun-12        3.525m ±  3%
Query/Mid_xorm-12       3.432m ±  1%
Query/Simple_sql-12     678.5µ ±  2%
Query/Simple_gorm-12    1.313m ± 12%
Query/Simple_sqlt-12    683.2µ ±  3%
Query/Simple_ent-12     914.6µ ±  1%
Query/Simple_sqlc-12    874.3µ ±  1%
Query/Simple_bun-12     895.6µ ±  1%
Query/Simple_xorm-12    819.6µ ±  2%
Read/sql-12             25.48µ ±  1%
Read/gorm-12            88.42µ ±  4%
Read/sqlt-12            26.54µ ±  3%
Read/ent-12             59.40µ ±  8%
Read/sqlc-12            31.33µ ±  1%
Read/bun-12             47.74µ ±  1%
Read/xorm-12            36.16µ ±  2%
geomean                 1.370m

                      │  bench.out   │
                      │     B/op     │
Query/Complex_sql-12    13.98Ki ± 1%
Query/Complex_gorm-12   121.3Ki ± 0%
Query/Complex_sqlt-12   13.42Ki ± 1%
Query/Complex_ent-12    82.12Ki ± 0%
Query/Complex_sqlc-12   13.02Ki ± 0%
Query/Complex_bun-12    54.96Ki ± 0%
Query/Complex_xorm-12   36.93Ki ± 0%
Query/Mid_sql-12        10.07Ki ± 0%
Query/Mid_gorm-12       101.9Ki ± 0%
Query/Mid_sqlt-12       9.856Ki ± 0%
Query/Mid_ent-12        66.06Ki ± 0%
Query/Mid_sqlc-12       8.999Ki ± 0%
Query/Mid_bun-12        48.02Ki ± 0%
Query/Mid_xorm-12       27.26Ki ± 0%
Query/Simple_sql-12     77.36Ki ± 0%
Query/Simple_gorm-12    590.3Ki ± 0%
Query/Simple_sqlt-12    83.02Ki ± 0%
Query/Simple_ent-12     306.3Ki ± 0%
Query/Simple_sqlc-12    87.47Ki ± 0%
Query/Simple_bun-12     195.4Ki ± 0%
Query/Simple_xorm-12    201.1Ki ± 0%
Read/sql-12             2.297Ki ± 0%
Read/gorm-12            58.61Ki ± 0%
Read/sqlt-12            3.461Ki ± 0%
Read/ent-12             32.83Ki ± 0%
Read/sqlc-12            2.242Ki ± 0%
Read/bun-12             35.68Ki ± 0%
Read/xorm-12            10.66Ki ± 0%
geomean                 35.21Ki

                      │  bench.out  │
                      │  allocs/op  │
Query/Complex_sql-12     271.0 ± 0%
Query/Complex_gorm-12   2.439k ± 0%
Query/Complex_sqlt-12    308.0 ± 0%
Query/Complex_ent-12    1.914k ± 0%
Query/Complex_sqlc-12    250.0 ± 0%
Query/Complex_bun-12    1.126k ± 0%
Query/Complex_xorm-12    786.0 ± 0%
Query/Mid_sql-12         219.0 ± 0%
Query/Mid_gorm-12       2.022k ± 0%
Query/Mid_sqlt-12        252.0 ± 0%
Query/Mid_ent-12        1.540k ± 0%
Query/Mid_sqlc-12        201.0 ± 0%
Query/Mid_bun-12         898.0 ± 0%
Query/Mid_xorm-12        627.0 ± 0%
Query/Simple_sql-12     1.677k ± 0%
Query/Simple_gorm-12    12.20k ± 0%
Query/Simple_sqlt-12    1.863k ± 0%
Query/Simple_ent-12     6.698k ± 0%
Query/Simple_sqlc-12    1.513k ± 0%
Query/Simple_bun-12     5.928k ± 0%
Query/Simple_xorm-12    5.035k ± 0%
Read/sql-12              69.00 ± 0%
Read/gorm-12            1.004k ± 0%
Read/sqlt-12             95.00 ± 0%
Read/ent-12              848.0 ± 0%
Read/sqlc-12             67.00 ± 0%
Read/bun-12              414.0 ± 0%
Read/xorm-12             261.0 ± 0%
geomean                  774.5
```