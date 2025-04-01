package benchflix_test

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	benchflix "github.com/wroge/bench-flix"
	bunflix "github.com/wroge/bench-flix/bun-flix"
	entflix "github.com/wroge/bench-flix/ent-flix"
	gormflix "github.com/wroge/bench-flix/gorm-flix"
	sqlflix "github.com/wroge/bench-flix/sql-flix"
	sqlcflix "github.com/wroge/bench-flix/sqlc-flix"
	sqltflix "github.com/wroge/bench-flix/sqlt-flix"
	xormflix "github.com/wroge/bench-flix/xorm-flix"
)

var repositories = []func() benchflix.Repository{
	sqlflix.NewRepository,
	gormflix.NewRepository,
	sqltflix.NewRepository,
	entflix.NewRepository,
	sqlcflix.NewRepository,
	bunflix.NewRepository,
	xormflix.NewRepository,
}

type Case struct {
	Name   string
	Query  benchflix.Query
	Result string
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
				Search:  "Affleck",
				Country: "United Kingdom",
				Genre:   "Drama",
			},
			Result: `[{68734 Argo 2012-10-11 00:00:00 +0000 UTC [Ben Affleck] [Alan Arkin Ben Affleck Bryan Cranston John Goodman Victor Garber] [United Kingdom United States of America] 7.278 [Drama Thriller]} {157336 Interstellar 2014-11-05 00:00:00 +0000 UTC [Christopher Nolan] [Anne Hathaway Casey Affleck Jessica Chastain Matthew McConaughey Michael Caine] [United Kingdom United States of America] 8.5 [Adventure Drama Science Fiction]} {37414 The Killer Inside Me 2010-02-19 00:00:00 +0000 UTC [Michael Winterbottom] [Casey Affleck Jessica Alba Kate Hudson Ned Beatty Tom Bower] [Canada Sweden United Kingdom United States of America] 5.8 [Crime Drama Thriller]} {505225 The Last Thing He Wanted 2020-02-14 00:00:00 +0000 UTC [Dee Rees] [Anne Hathaway Ben Affleck Edi Gathegi Rosie Perez Willem Dafoe] [United Kingdom United States of America] 4.9 [Drama Thriller]} {23168 The Town 2010-09-15 00:00:00 +0000 UTC [Ben Affleck] [Ben Affleck Blake Lively Jeremy Renner Jon Hamm Rebecca Hall] [United Kingdom United States of America] 7.2 [Crime Drama Thriller]}]`,
		},
		{
			Name: "Mid",
			Query: benchflix.Query{
				Search:     "Affleck",
				AddedAfter: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			Result: `[{916728 Slingshot 2024-08-30 00:00:00 +0000 UTC [Mikael Håfström] [Casey Affleck David Morrissey Emily Beecham Laurence Fishburne Tomer Capone] [Hungary United States of America] 6.1 [Science Fiction Thriller]} {870028 The Accountant 2 2025-04-23 00:00:00 +0000 UTC [Gavin O'Connor] [Ben Affleck Cynthia Addai-Robinson Daniella Pineda J.K. Simmons Jon Bernthal] [United States of America] 0 [Action Crime Thriller]} {1059064 The Instigators 2024-08-02 00:00:00 +0000 UTC [Doug Liman] [Casey Affleck Jack Harlow Matt Damon Michael Stuhlbarg Ving Rhames] [United States of America] 6.2 [Action Comedy Crime]} {1217343 This Is Me…Now 2024-02-15 00:00:00 +0000 UTC [Dave Meyers] [Ben Affleck Fat Joe Idaliz Christian Jennifer Lopez Matthew Law] [United States of America] 5.3 [Drama Music]}]`,
		},
		{
			Name: "Simple",
			Query: benchflix.Query{
				MinRating: 9.5,
				MaxRating: 10,
			},
			Result: `[{723562 Actresses: Sex Audition 2020-07-03 00:00:00 +0000 UTC [Jeong Mi-na] [Han Seok-bong Jo Wan-jin Lee Chae-dam Lee Soo Tae Hee] [South Korea] 10 [Drama]} {1285792 An Unholy Affair: A Younger Man and a Busty Wife 2017-01-19 00:00:00 +0000 UTC [Sadao Sadaoka] [Saeko Matsushita Yuta Aoi] [Japan] 10 []} {974573 Another Simple Favor 2025-03-07 00:00:00 +0000 UTC [Paul Feig] [Andrew Rannells Anna Kendrick Bashir Salahuddin Blake Lively Henry Golding] [Canada United States of America] 9.5 [Comedy Crime Thriller]} {1164488 Balota 2024-08-02 00:00:00 +0000 UTC [Kip Oebanda] [Donna Cariaga Marian Rivera Nico Antonio Royce Cabrera Will Ashley] [Philippines] 10 [Drama Thriller]} {1440249 Butterfly 2025-02-28 00:00:00 +0000 UTC [Monish] [Jayanth Kamalesh Logeshwaran Sanjay] [India] 10 []} {1313003 Close To Me 2025-03-06 00:00:00 +0000 UTC [Stefano Sardo] [Giulio Beranek Maria Chiara Giannetta Mariela Garriga Paolo Pierobon Riccardo Scamarcio] [Italy] 10 [Romance Thriller]} {938086 Dangerous Younger Cousin 2021-11-23 00:00:00 +0000 UTC [Kim Tae-hoon] [Ji Min-iI Lee Chae-dam Min Do-yoon Yoo Ji-hyun Yoon Yool] [South Korea] 10 [Drama Family]} {1273049 Dragon 2025-02-21 00:00:00 +0000 UTC [Ashwath Marimuthu] [Anupama Parameswaran Gautham Vasudev Menon Kayadu Lohar Mysskin Pradeep Ranganathan] [India] 10 [Comedy Drama Romance]} {1262299 El Apocalipsis de san Juan 2024-10-07 00:00:00 +0000 UTC [Patricio Dondo Simón Delacre] [Carlos Secilio Héctor Hugo Larriera Miguel Angel Marchessi Ricardo Castro] [Argentina] 10 [Documentary History]} {1436457 F1 75 Live at The O2 2025-02-18 00:00:00 +0000 UTC [] [Alexander Albon Carlos Sainz Jr. Gabriel Bortoleto Nico Hülkenberg Yuki Tsunoda] [] 10 [Documentary]} {1037837 Family Matters 2022-12-25 00:00:00 +0000 UTC [Nuel C. Naval] [Agot Isidro Liza Lorena Mylene Dizon Noel Trinidad Nonie Buencamino] [Philippines] 10 [Drama]} {653127 Female Hostel 3 2019-06-11 00:00:00 +0000 UTC [Jo Tae-ho] [Choi Chul-min Kim In-ae Park Hyun-jung Yoon Da-hyeon] [South Korea] 9.8 [Romance]} {1412569 Hiram na Sandali 2025-01-16 00:00:00 +0000 UTC [G.B. Sampedro] [Aerol Carmelo Ara Altamira Denise Esteban Dyessa Garcia Vince Rillon] [Philippines] 10 [Drama Romance]} {546026 Inácio Garapa, Um Matuto Sonhador 2010-05-24 00:00:00 +0000 UTC [J. Gomes] [Maria Elivanete Wellington Marques] [] 10 [Comedy Drama Family]} {501007 It 2017-12-14 00:00:00 +0000 UTC [Anouk de Clercq Tom Callemin] [] [Belgium] 10 []} {1315905 Jailbreak Affair 2024-06-14 00:00:00 +0000 UTC [Ham Seo-yong] [Jin Seo-yool Lee Chae-dam Lee Ye-jin Min Do-yoon Son Ye-ha] [] 10 []} {1421982 Lore Of The Ring Light 2025-01-21 00:00:00 +0000 UTC [Allison Craig] [Eric Bauza Kimiko Glenn Michael Croner Nic Smal Vella Lovell] [] 10 [Adventure Animation Comedy Family Fantasy Music]} {1319933 Marco 2024-06-29 00:00:00 +0000 UTC [Arnau García Eneko Amasene Guifré Surinyac Martí Garcia Solé] [Alan García Arnau García Eneko Amasene Guifré Surinyac Martí Garcia Solé] [] 10 [Drama]} {1196942 Mere Husband Ki Biwi 2025-02-21 00:00:00 +0000 UTC [Mudassar Aziz] [Aditya Seal Arjun Kapoor Bhumi Pednekar Harsh Gujral Rakul Preet Singh] [India] 10 [Comedy Drama]} {696026 Nice Sister-In-Law 2019-06-05 00:00:00 +0000 UTC [Wang Ji-bang] [Ha Jin Jin Joo Jo Wan-jin] [South Korea] 10 [Drama Romance]} {1033023 Nobody Likes Me 2025-01-09 00:00:00 +0000 UTC [Petr Kazda Tomáš Weinreb] [Barbora Bobuľová Hana Vagnerová Leona Skleničková Mantas Zemleckas Rebeka Poláková] [Czech Republic France Slovakia] 9.5 [Drama Romance]} {484133 Nude 2017-10-29 00:00:00 +0000 UTC [Tony Sacco] [David Bellemere Jeannie Park Jessica Clements Rachel Cook Steve Shaw] [United States of America] 9.5 [Documentary]} {979480 Open Marriage: Aru Fuufu no Katachi 2018-11-02 00:00:00 +0000 UTC [Jirō Ishikawa] [Gaichi Masami Ichikawa Ren Fukusaki Ryozo Sousuke Yamamoto] [] 9.5 [Drama]} {1144932 Queen of the Ring 2025-03-07 00:00:00 +0000 UTC [Ash Avildsen] [Emily Bett Rickards Gavin Casalegno Josh Lucas Marie Avgeropoulos Walton Goggins] [United States of America] 10 [Drama]} {1211726 Red Silk 2025-02-20 00:00:00 +0000 UTC [Andrey Volgin] [Elena Podkaminskaya Gleb Kalyuzhny Gosha Kutsenko Miloš Biković Zheng Hanyi] [Russia] 10 [Action Adventure Mystery]} {1135869 Salome 2023-05-13 00:00:00 +0000 UTC [] [Ariele Bellus Bens Bima Prawira Radja Adipati Virly Virginia] [Indonesia] 10 [Drama Romance]} {688876 Secret Night Of Mother And Daughter 2020-03-27 00:00:00 +0000 UTC [Yoon Kyung-sik] [Jin Yi Lee Chae-dam Shin Yeon-woo Sin Seong-hoon] [South Korea] 10 [Drama Romance]} {694940 Sincheon Station Exit 3 2020-04-02 00:00:00 +0000 UTC [No Hyun-jin] [Ahn So-hee Hae Il Lee Sul-ah Min Do-yoon Si Woo] [South Korea] 10 [Drama Romance]} {481994 Some: An Erotic Tale 2017-04-13 00:00:00 +0000 UTC [Song Jeong-gyoo] [Ahn So-hee Han Ga-hee Hayashi Risa Lee Jae-seok Sang Woo] [South Korea] 9.5 [Comedy Drama]} {1130276 Succubus 2024-08-30 00:00:00 +0000 UTC [R.J. Daniel Hanna] [Brendan Bradley Derek Smith Emily Kincaid Olivia Grace Applegate Rachel Cook] [United States of America] 9.5 [Horror Thriller]} {875067 Swapping Guest House 2018-09-27 00:00:00 +0000 UTC [Park Eun-soo-I] [Park Joo-bin Sang Woo Seo Won Shin Joon-hyun] [South Korea] 10 [Romance]} {1149913 TOGEFILM - Mei Mei 2023-04-06 00:00:00 +0000 UTC [] [Virly Virginia] [] 10 [Drama Romance]} {1191815 The American Backyard 2025-03-06 00:00:00 +0000 UTC [Pupi Avati] [Armando De Ceccon Filippo Scotti Mildred Gustafsson Rita Tushingham Roberto De Francesco] [Italy United States of America] 10 [Drama Horror Mystery Thriller]} {1352874 The Crucifix 2025-01-09 00:00:00 +0000 UTC [Stephen Roach] [Darren Le Fevre Dean Kilbey Hannaj Bang Bendz Nicholas Anscombe Tom Carter] [United Kingdom] 10 [Horror]} {1192174 The Kite 2025-03-06 00:00:00 +0000 UTC [Alessandro Tonda] [Andrea Giannini Anna Ferzetti Claudio Santamaria Massimiliano Rossi Sonia Bergamasco] [Belgium Italy] 10 [Drama Thriller]} {445585 The Photographer 2017-02-10 00:00:00 +0000 UTC [Ji Hyun-sook] [Chae Won Jo Yong-bok Lee Soo Lee Soo-jung Oh Joo-ha] [South Korea] 10 []} {794253 The Shepherd 2020-10-20 00:00:00 +0000 UTC [Yiannis Stravolaimos] [Elena Thomopoulou Kostis Savvidakis Leyteris Tsatsis Thanasis Nakos Theo Theodoridis] [Greece] 10 [Drama]} {1317572 The Williams 2024-10-31 00:00:00 +0000 UTC [Raúl de la Fuente] [Iñaki Williams Nico Williams] [Spain] 10 [Documentary]} {1434243 The Wrong Obsession 2025-02-21 00:00:00 +0000 UTC [David DeCoteau] [Daniel Joo Gina Hiraizumi Matthew Pohlkamp Morgan Bradley Vivica A. Fox] [United States Minor Outlying Islands] 10 []} {721183 Three Sisters Swapping 2019-11-06 00:00:00 +0000 UTC [] [Han Seok-bong Ji Yeong Kang Min-woo Sae Bom Yoo Jung] [South Korea] 10 [Romance]} {1225093 Underpants Thief 2021-12-23 00:00:00 +0000 UTC [Somaratne Dissanayake] [Buddhi Randeniya Chinthaka Kulathunga Dilhani Ekanayake Pubudu Chathuranga] [Sri Lanka] 10 [Drama]} {552866 Youthful Mother-in-law 2018-06-15 00:00:00 +0000 UTC [Choi Eun-jung] [James Jin Joo Lee Sang-doo Si Ah 陈诗雅] [] 10 [Drama Romance]}]`,
		},
	}

	readCases = []IDCase{
		{
			ID:     10192,
			Result: `{10192 Shrek Forever After 2010-05-16 00:00:00 +0000 UTC [Mike Mitchell] [Antonio Banderas Cameron Diaz Eddie Murphy Mike Myers Walt Dohrn] [United States of America] 6.38 [Adventure Animation Comedy Family Fantasy]}`,
		},
	}
)

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
		for _, init := range repositories {
			r := init()

			t.Run(c.Name+"_"+strings.TrimSuffix(reflect.TypeOf(r).String(), "flix.Repository"), func(t *testing.T) {
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

				if fmt.Sprint(movies) != c.Result {
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

	for _, c := range readCases {
		for _, init := range repositories {
			r := init()

			t.Run(strings.TrimSuffix(reflect.TypeOf(r).String(), "flix.Repository"), func(t *testing.T) {
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

	for _, c := range readCases {
		for _, init := range repositories {
			r := init()

			for _, record := range records[1:] {
				movie, err := benchflix.NewMovie(record)
				if err != nil {
					b.Fatal(err)
				}

				if err = r.Create(ctx, movie); err != nil {
					b.Fatal(err)
				}
			}

			b.Run(strings.TrimSuffix(reflect.TypeOf(r).String(), "flix.Repository"), func(b *testing.B) {
				for b.Loop() {
					movie, err := r.Read(ctx, c.ID)
					if err != nil {
						b.Fatal(reflect.TypeOf(r), err)
					}

					if fmt.Sprint(movie) != c.Result {
						b.Fatal(reflect.TypeOf(r), movie)
					}
				}
			})
		}
	}
}
