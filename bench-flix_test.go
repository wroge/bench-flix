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
			Name: "Simple",
			Query: benchflix.Query{
				Search:     "Affleck",
				AddedAfter: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			Result: `[{916728 Slingshot 2024-08-30 00:00:00 +0000 UTC [Mikael Håfström] [Casey Affleck David Morrissey Emily Beecham Laurence Fishburne Tomer Capone] [Hungary United States of America] 6.1 [Science Fiction Thriller]} {870028 The Accountant 2 2025-04-23 00:00:00 +0000 UTC [Gavin O'Connor] [Ben Affleck Cynthia Addai-Robinson Daniella Pineda J.K. Simmons Jon Bernthal] [United States of America] 0 [Action Crime Thriller]} {1059064 The Instigators 2024-08-02 00:00:00 +0000 UTC [Doug Liman] [Casey Affleck Jack Harlow Matt Damon Michael Stuhlbarg Ving Rhames] [United States of America] 6.2 [Action Comedy Crime]} {1217343 This Is Me…Now 2024-02-15 00:00:00 +0000 UTC [Dave Meyers] [Ben Affleck Fat Joe Idaliz Christian Jennifer Lopez Matthew Law] [United States of America] 5.3 [Drama Music]}]`,
		},
		{
			Name: "10",
			Query: benchflix.Query{
				Limit: 10,
			},
			Result: `[{1013577 "Sr." 2022-09-02 00:00:00 +0000 UTC [Chris Smith] [Alan Arkin Chris Smith Robert Downey Jr. Robert Downey Sr. Sean Hayes] [United States of America] 6.9 [Documentary]} {41371 #1 Cheerleader Camp 2010-07-27 00:00:00 +0000 UTC [Mark Quod] [Charlene Tilton Erica Duke Harmony Blossom Hector David Jr. Jay Gillespie] [United States of America] 4.8 [Comedy]} {1090323 #AMFAD: All My Friends Are Dead 2024-08-29 00:00:00 +0000 UTC [Marcus Dunstan] [Ali Fumiko Whitney Jade Pettyjohn Jennifer Ens JoJo Siwa Michaella Russell] [United States of America] 6 [Horror Thriller]} {614696 #Alive 2020-06-24 00:00:00 +0000 UTC [Cho Il] [Jin So-yeon Kim Hak-seon Lee Hyun-wook Park Shin-hye Yoo Ah-in] [South Korea] 7.234 [Action Horror]} {714567 #Blue_Whale 2024-01-11 00:00:00 +0000 UTC [Anna Zaytseva] [Anna Potebnya Diana Shulmina Ekaterina Stulova Olga Pipchenko Timofey Eletskiy] [Russia] 6.3 [Drama Horror Thriller]} {301325 #Horror 2015-11-20 00:00:00 +0000 UTC [Tara Subkoff] [Balthazar Getty Chloë Sevigny Natasha Lyonne Taryn Manning Timothy Hutton] [United States of America] 3.591 [Crime Horror Thriller]} {605734 #Iamhere 2020-02-05 00:00:00 +0000 UTC [Eric Lartigau] [Alain Chabat Bae Doona Camille Rutherford Ilian Bergala Jules Sagot] [Belgium France] 5.6 [Comedy Romance]} {295050 #Stuck 2014-07-01 00:00:00 +0000 UTC [Stuart Acher] [Abraham Benrubi Jayson Blair Joanna Canton Joel David Moore Madeline Zima] [United States of America] 5.5 [Comedy Drama Romance]} {581802 #UNFIT: The Psychology of Donald Trump 2020-10-22 00:00:00 +0000 UTC [Dan Partland] [Donald Trump George Thomas Conway III John Gartner Kellyanne Conway Malcolm Nance] [United States of America] 6.608 [Documentary]} {455656 #realityhigh 2017-07-17 00:00:00 +0000 UTC [Fernando Lebrija] [Alicia Sanz Anne Winters Jake Borelli Keith Powers Nesta Cooper] [United States of America] 6.291 [Comedy]}]`,
		},
		{
			Name: "100",
			Query: benchflix.Query{
				Limit: 100,
			},
			Result: `[{1013577 "Sr." 2022-09-02 00:00:00 +0000 UTC [Chris Smith] [Alan Arkin Chris Smith Robert Downey Jr. Robert Downey Sr. Sean Hayes] [United States of America] 6.9 [Documentary]} {41371 #1 Cheerleader Camp 2010-07-27 00:00:00 +0000 UTC [Mark Quod] [Charlene Tilton Erica Duke Harmony Blossom Hector David Jr. Jay Gillespie] [United States of America] 4.8 [Comedy]} {1090323 #AMFAD: All My Friends Are Dead 2024-08-29 00:00:00 +0000 UTC [Marcus Dunstan] [Ali Fumiko Whitney Jade Pettyjohn Jennifer Ens JoJo Siwa Michaella Russell] [United States of America] 6 [Horror Thriller]} {614696 #Alive 2020-06-24 00:00:00 +0000 UTC [Cho Il] [Jin So-yeon Kim Hak-seon Lee Hyun-wook Park Shin-hye Yoo Ah-in] [South Korea] 7.234 [Action Horror]} {714567 #Blue_Whale 2024-01-11 00:00:00 +0000 UTC [Anna Zaytseva] [Anna Potebnya Diana Shulmina Ekaterina Stulova Olga Pipchenko Timofey Eletskiy] [Russia] 6.3 [Drama Horror Thriller]} {301325 #Horror 2015-11-20 00:00:00 +0000 UTC [Tara Subkoff] [Balthazar Getty Chloë Sevigny Natasha Lyonne Taryn Manning Timothy Hutton] [United States of America] 3.591 [Crime Horror Thriller]} {605734 #Iamhere 2020-02-05 00:00:00 +0000 UTC [Eric Lartigau] [Alain Chabat Bae Doona Camille Rutherford Ilian Bergala Jules Sagot] [Belgium France] 5.6 [Comedy Romance]} {295050 #Stuck 2014-07-01 00:00:00 +0000 UTC [Stuart Acher] [Abraham Benrubi Jayson Blair Joanna Canton Joel David Moore Madeline Zima] [United States of America] 5.5 [Comedy Drama Romance]} {581802 #UNFIT: The Psychology of Donald Trump 2020-10-22 00:00:00 +0000 UTC [Dan Partland] [Donald Trump George Thomas Conway III John Gartner Kellyanne Conway Malcolm Nance] [United States of America] 6.608 [Documentary]} {455656 #realityhigh 2017-07-17 00:00:00 +0000 UTC [Fernando Lebrija] [Alicia Sanz Anne Winters Jake Borelli Keith Powers Nesta Cooper] [United States of America] 6.291 [Comedy]} {1303498 $POSITIONS 2025-03-08 00:00:00 +0000 UTC [Brandon Daley] [Kaylyn Carter Michael Kunicki Reagan Fitzgerald Trevor Dawkins Vinny Kress] [United States of America] 0 [Drama]} {252178 '71 2014-10-10 00:00:00 +0000 UTC [Yann Demange] [Jack O'Connell Paul Anderson Sam Hazeldine Sam Reid Sean Harris] [United Kingdom] 6.803 [Action Drama Thriller War]} {1100642 'Twas the Text Before Christmas 2023-10-21 00:00:00 +0000 UTC [T.W. Peacocke] [Jayne Eastwood Marisa McIntyre Merritt Patterson Rob Stewart Trevor Donovan] [Canada United States of America] 6.4 [Comedy Romance TV Movie]} {470211 (Girl)Friend 2018-01-17 00:00:00 +0000 UTC [Victor Saint Macary] [Béatrice de Staël Camille Razat Jonathan Cohen Margot Bancilhon William Lebghil] [France] 5.4 [Comedy Romance]} {176068 +1 2013-09-20 00:00:00 +0000 UTC [Dennis Iliadis] [Ashley Grace Logan Miller Natalie Hall Rhys Wakefield Suzanne Dengel] [United States of America] 5.3 [Science Fiction Thriller]} {838197 ...Watch Out, We're Mad 2022-03-23 00:00:00 +0000 UTC [Antonio Usbergo Niccolò Celaia] [Alessandra Mastronardi Alessandro Roja Christian De Sica Edoardo Pesce Francesco Bruni] [Italy] 6.1 [Action Comedy]} {584586 0.0MHz 2019-05-29 00:00:00 +0000 UTC [You Sun-dong] [Choi Yoon-young Jung Eun-ji Jung Won-chang Lee Sung-yeol Shin Ju-hwan] [South Korea] 4.9 [Horror]} {1257926 0000 Kilometre 2024-10-31 00:00:00 +0000 UTC [Deniz Enyüksek] [Ahmet Haktan Zavlak Cavit Çetin Güner Derya Pınar Ak Gülin İyigün Ogün Kaptanoğlu] [Turkey] 5.5 [Drama Romance]} {1436465 0004ngel 2025-02-26 00:00:00 +0000 UTC [Eli Jean Tahchi] [] [Canada] 0 []} {127544 009 Re:Cyborg 2012-10-27 00:00:00 +0000 UTC [Kenji Kamiyama] [Chiwa Saito Daisuke Ono Mamoru Miyano Sakiko Tamagawa Toru Okawa] [Japan] 7 [Action Animation Science Fiction]} {259300 009-1: The End of the Beginning 2013-09-07 00:00:00 +0000 UTC [Koichi Sakamoto] [Aya Sugimoto Mao Ichimichi Mayuko Iwasa Minehiro Kinomoto Nao Nagasawa] [Japan] 5.8 [Action Science Fiction]} {217316 1 2013-09-30 00:00:00 +0000 UTC [Paul Crowder] [Jenson Button Lewis Hamilton Michael Fassbender Michael Schumacher Niki Lauda] [United States of America] 6.587 [Documentary]} {586032 1 2020-09-14 00:00:00 +0000 UTC [Andrzej Kozlowski] [Alice Amter Ethan Phillips Jeremy Craven Jude Ciccolella London Bridges] [] 6.705 [Drama]} {416691 1 Night 2016-10-14 00:00:00 +0000 UTC [Minhal Baig] [Anna Camp Isabelle Fuhrman Justin Chatwin Kelli Berglund Kyle Allen] [United States of America] 6.2 [Drama Romance]} {1080996 1+1+1 Life, Love, Chaos 2025-02-28 00:00:00 +0000 UTC [Yanie Dupont-Hébert] [Dany Lefebvre Irlande Côté Matai Stevens Noémie Yelle Victor Andres Trelles Turgeon] [Canada] 0 [Comedy Drama]} {333371 10 Cloverfield Lane 2016-03-10 00:00:00 +0000 UTC [Dan Trachtenberg] [Douglas M. Griffin John Gallagher Jr. John Goodman Mary Elizabeth Winstead Suzanne Cryer] [United States of America] 6.991 [Drama Horror Science Fiction Thriller]} {1074262 10 Days of a Bad Man 2023-08-18 00:00:00 +0000 UTC [Uluç Bayraktar] [Erdal Yıldız Ilayda Akdoğan Nejat İşler Nur Fettahoğlu Şenay Gürler] [Turkey] 6.5 [Crime Drama]} {1169361 10 Days of a Curious Man 2024-11-06 00:00:00 +0000 UTC [Uluç Bayraktar] [Ece İrtem Hazal Subaşı Ilayda Akdoğan Nejat İşler Şenay Gürler] [Turkey] 6.3 [Crime Drama Thriller]} {1073337 10 Days of a Good Man 2023-02-10 00:00:00 +0000 UTC [Uluç Bayraktar] [Ilayda Akdoğan Ilayda Alişan Nejat İşler Nur Fettahoğlu Şenay Gürler] [Turkey] 6.6 [Crime Drama]} {605735 10 Days with Dad 2020-02-19 00:00:00 +0000 UTC [Ludovic Bernard] [Alexis Michalik Alice David Aure Atika Franck Dubosc Héléna Noguerra] [France] 6 [Comedy Family]} {567811 10 Lives 2024-04-18 00:00:00 +0000 UTC [Christopher Jenkins] [Dylan Llewellyn Mo Gilligan Simone Ashley Sophie Okonedo Zayn Malik] [Canada United Kingdom United States of America] 7.8 [Animation Comedy Family Fantasy]} {552865 10 Minutes Gone 2019-09-30 00:00:00 +0000 UTC [Brian A. Miller] [Bruce Willis Kyle Schmid Meadow Williams Michael Chiklis Texas Battle] [United States of America] 5.3 [Action Crime Mystery Thriller]} {106942 10 Rules for Falling in Love 2012-03-16 00:00:00 +0000 UTC [Cristiano Bortone] [Enrica Pintore Giulio Berruti Guglielmo Scilla Pietro Masotti Vincenzo Salemme] [Italy] 5.102 [Comedy]} {216138 10 Rules for Sleeping Around 2013-08-23 00:00:00 +0000 UTC [Leslie Greif] [Bryan Callen Christopher Rodriguez Marquette Jesse Bradford Tammin Sursok Virginia Williams] [United States of America] 3.683 [Comedy Romance]} {534039 10 Things We Should Do Before We Break Up 2020-03-05 00:00:00 +0000 UTC [Galt Niederhoffer] [Brady Jenness Christina Ricci Hamish Linklater Lindsey Broad Mia Sinclair Jenness] [United States of America] 5.3 [Drama Romance]} {58547 10 Years 2012-09-14 00:00:00 +0000 UTC [Jamie Linden] [Channing Tatum Jenna Dewan Justin Long Max Minghella Oscar Isaac] [United States of America] 5.8 [Comedy Drama]} {253251 10,000 Km 2014-05-16 00:00:00 +0000 UTC [Carlos Marques-Marcet] [David Verdaguer Natalia Tena] [Spain] 6.1 [Drama Romance]} {253406 10,000 Saints 2015-08-14 00:00:00 +0000 UTC [Robert Pulcini Shari Springer Berman] [Asa Butterfield Emily Mortimer Ethan Hawke Hailee Steinfeld Julianne Nicholson] [United States of America] 6 [Comedy Drama Music]} {302666 10.0 Earthquake 2014-10-15 00:00:00 +0000 UTC [David Gidali] [Cameron Richardson Chasty Ballesteros Heather Sossaman Henry Ian Cusick Jeffrey Jones] [United States of America] 4.942 [Action Adventure Drama]} {126757 100 Bloody Acres 2012-08-04 00:00:00 +0000 UTC [Cameron Cairnes Colin Cairnes] [Angus Sampson Anna McGahan Damon Herriman John Jarratt Oliver Ackland] [Australia] 5.8 [Comedy Horror]} {1211728 100 Candles Game: The Last Possession 2023-11-09 00:00:00 +0000 UTC [Andrés Borghi Arie Socorro Carlos Goitia David Ferino Guillermo Lockhart Jerónimo Rocha Maximilian Niemann Ryan Graff] [Josefina Inés Fariña Justina Ceballos Magui Bravi Nacho Francavilla Zhon Li] [Argentina New Zealand] 6.25 [Horror]} {182228 100 Degrees Below Zero 2013-03-29 00:00:00 +0000 UTC [Richard Schenkman] [Andray Johnson Jeff Fahey John Rhys-Davies Judit Fekete Sara Malakul Lane] [United States of America] 4.3 [Action Science Fiction]} {402693 100 Meters 2016-11-04 00:00:00 +0000 UTC [Marcel Barrena] [Alexandra Jiménez Dani Rovira David Verdaguer Karra Elejalde Maria de Medeiros] [Spain] 7.3 [Comedy Drama]} {334532 100 Streets 2016-11-11 00:00:00 +0000 UTC [Jim O'Hanlon] [Gemma Arterton Idris Elba Kierston Wareing Ryan Gage Tom Cullen] [United Kingdom] 6.2 [Drama]} {546630 100 Things 2018-12-06 00:00:00 +0000 UTC [Florian David Fitz] [Florian David Fitz Hannelore Elsner Matthias Schweighöfer Miriam Stein Wolfgang Stumph] [Germany] 6.609 [Comedy Drama]} {1122824 100 Yards 2024-09-20 00:00:00 +0000 UTC [Xu Haofeng Xu Junfeng] [Andy On Chi-Kit Bea Hayden Kuo Jacky Heung Li Yuan Tang Shiyi] [China] 5.7 [Action Drama]} {314606 100 Yen Love 2014-11-15 00:00:00 +0000 UTC [Masaharu Take] [Hirofumi Arai Miyoko Inagawa Sakura Ando Saori Toshie Negishi] [Japan] 7.2 [Comedy Drama Romance]} {1327145 100 dni do matury 2025-02-28 00:00:00 +0000 UTC [Mikołaj Piszczan] [Bartosz Kubicki Bartosz Laskowski Kinga Banaś Patryk Baran Pola Sieczko] [Poland] 5 [Action Comedy]} {520946 100% Wolf 2020-06-26 00:00:00 +0000 UTC [Alexs Stadermann] [Ilai Swindells Jai Courtney Jane Lynch Rhys Darby Samara Weaving] [Australia Belgium] 6.1 [Adventure Animation Family Fantasy]} {852590 1000 Miles From Christmas 2021-12-24 00:00:00 +0000 UTC [Álvaro Fernández Armero] [Andrea Ros Fermí Reixach Peter Vives Tamar Novas Verónica Forqué] [Spain] 6.3 [Comedy Family Romance]} {268245 1000 to 1 2014-03-04 00:00:00 +0000 UTC [Michael Levine] [Cassi Thomson David Henrie Hannah Marks Luke Kleintank Myk Watford] [United States of America] 5.905 [Drama]} {337302 10000 Years Later 2015-03-27 00:00:00 +0000 UTC [Li Yi] [Chong Wang Joma Yalayam] [China] 6.7 [Action Animation Fantasy]} {505177 10x10 2018-03-16 00:00:00 +0000 UTC [Suzi Ewing] [Jason Maza Kelly Reilly Luke Evans Noel Clarke Olivia Chenery] [United Kingdom] 5.2 [Drama Mystery Thriller]} {51248 11-11-11 2011-11-11 00:00:00 +0000 UTC [Darren Lynn Bousman] [Brendan Price Lluís Soler Michael Landes Timothy Gibbs Wendy Glenn] [United States of America] 4.6 [Horror Thriller]} {182246 11.6 2013-04-03 00:00:00 +0000 UTC [Philippe Godeau] [Bouli Lanners Corinne Masiero François Cluzet Johan Libéreau Juana Acosta] [France] 5.565 [Drama Thriller]} {79078 11/11/11 2011-11-01 00:00:00 +0000 UTC [Keith Allan] [Erin Coker Hayden Byerly Jon Briddell Rebecca Light Scott McKinley] [United States of America] 3.9 [Horror Thriller]} {81393 12 Dates of Christmas 2011-12-11 00:00:00 +0000 UTC [James Hayman] [Amy Smart Benjamin Ayres Mark-Paul Gosselaar Mary Long Peter MacNeill] [Canada Spain United States of America] 6.125 [Comedy Fantasy Romance TV Movie]} {459928 12 Feet Deep 2018-11-08 00:00:00 +0000 UTC [Matt Eskandari] [Alexandra Park Diane Farr Dogen Eyeler Nora-Jane Noone Tobin Bell] [United States of America] 6 [Thriller]} {363483 12 Gifts of Christmas 2015-11-26 00:00:00 +0000 UTC [Peter Sullivan] [Aaron O'Connell Alesandra Durham Donna Mills Katrina Law Melanie Nelson] [United States of America] 6.063 [Comedy Family Romance TV Movie]} {667141 12 Hour Shift 2020-10-02 00:00:00 +0000 UTC [Brea Grant] [Angela Bettis Chloe Farnworth David Arquette Kit Williamson Mick Foley] [United States of America] 5.5 [Comedy Crime Horror]} {625169 12 Mighty Orphans 2021-06-18 00:00:00 +0000 UTC [Ty Roberts] [Jake Austin Walker Luke Wilson Martin Sheen Vinessa Shaw Wayne Knight] [United States of America] 7.2 [Action Drama History]} {195269 12 Rounds 2: Reloaded 2013-06-01 00:00:00 +0000 UTC [Roel Reiné] [Brian Markinson Cindy Busby Randy Orton Sean Rogerson Tom Stevens] [United States of America] 5.516 [Action Adventure]} {351901 12 Rounds 3: Lockdown 2015-09-11 00:00:00 +0000 UTC [Stephen Reynolds] [Daniel Cudmore Jonathan Good Lochlyn Munro Roger Cross Sarah Smyth] [United States of America] 5.9 [Action Crime Thriller]} {429351 12 Strong 2018-01-18 00:00:00 +0000 UTC [Nicolai Fuglsig] [Chris Hemsworth Michael Peña Michael Shannon Navid Negahban Trevante Rhodes] [United States of America] 6.334 [Action Drama History War]} {566387 12 Suicidal Teens 2019-01-25 00:00:00 +0000 UTC [Yukihiko Tsutsumi] [Kanna Hashimoto Kotone Furukawa Mackenyu Mahiro Takasugi Yūto Fuchino] [Japan] 6.4 [Drama Mystery]} {76203 12 Years a Slave 2013-10-18 00:00:00 +0000 UTC [Steve McQueen] [Benedict Cumberbatch Chiwetel Ejiofor Lupita Nyong'o Michael Fassbender Paul Dano] [United Kingdom United States of America] 7.936 [Drama History]} {919207 12.12: The Day 2023-11-22 00:00:00 +0000 UTC [Kim Sung-soo] [Hwang Jung-min Jung Woo-sung Kim Sung-kyun Lee Sung-min Park Hae-jun] [South Korea] 7.5 [Crime Drama History Thriller War]} {147879 12/12/12 2012-12-04 00:00:00 +0000 UTC [Jared Cohn] [Carl Donelson Jesus Guevara Jon Kondelik Sara Malakul Lane Steve Hanks] [United States of America] 3.3 [Horror]} {44115 127 Hours 2010-11-12 00:00:00 +0000 UTC [Danny Boyle] [Amber Tamblyn Clémence Poésy James Franco Kate Mara Lizzy Caplan] [France United Kingdom United States of America] 7.1 [Adventure Drama Thriller]} {1163258 12th Fail 2023-08-11 00:00:00 +0000 UTC [Vidhu Vinod Chopra] [Anant Joshi Anshumaan Pushkar Medha Shankr Priyanshu Chatterjee Vikrant Massey] [India] 8 [Drama]} {44982 13 2010-03-12 00:00:00 +0000 UTC [Gela Babluani] [50 Cent Jason Statham Mickey Rourke Ray Winstone Sam Riley] [United States of America] 5.8 [Drama Thriller]} {58857 13 Assassins 2010-09-09 00:00:00 +0000 UTC [Takashi Miike] [Goro Inagaki Kazue Fukiishi Koji Yakusho Takayuki Yamada Yûsuke Iseya] [Japan United Kingdom] 7.3 [Action Adventure Drama]} {1126128 13 Bombs 2023-12-28 00:00:00 +0000 UTC [Angga Dwimas Sasongko] [Ardhito Pramono Chicco Kurniawan Lutesha Muhammad Khan Putri Ayudya] [Indonesia South Korea] 7.2 [Action Crime Thriller]} {347751 13 Cameras 2016-04-15 00:00:00 +0000 UTC [Victor Zarcoff] [Jim Cummings Neville Archambault PJ McCabe Sarah Baldwin Sean Carrigan] [United States of America] 5 [Crime Horror Thriller]} {122369 13 Eerie 2013-03-29 00:00:00 +0000 UTC [Lowell Dean] [Brendan Fehr Brendan Fletcher Katharine Isabelle Michael Shanks Nick Moran] [Canada] 4.7 [Horror Thriller]} {1026563 13 Exorcisms 2022-11-04 00:00:00 +0000 UTC [Jacobo Martínez] [José Sacristán María Romanillos Pablo Revuelta Ruth Díaz Urko Olazábal] [Spain] 5.8 [Drama Horror]} {300671 13 Hours: The Secret Soldiers of Benghazi 2016-01-14 00:00:00 +0000 UTC [Michael Bay] [Dominic Fumusa James Badge Dale John Krasinski Max Martini Pablo Schreiber] [Malta United States of America] 7.273 [Action Drama History Thriller War]} {319337 13 Minutes 2015-04-09 00:00:00 +0000 UTC [Oliver Hirschbiegel] [Burghart Klaußner Christian Friedel Felix Eitner Johann von Bülow Katharina Schüttler] [Germany] 6.8 [Drama History]} {787723 13 Minutes 2021-10-29 00:00:00 +0000 UTC [Lindsay Gossling] [Darryl Cox Thora Birch Tokala Black Elk Trace Adkins Yancey Arias] [Canada United States of America] 6 [Action Drama Thriller]} {155084 13 Sins 2014-04-11 00:00:00 +0000 UTC [Daniel Stamm] [Devon Graye Mark Webber Ron Perlman Rutina Wesley Tom Bower] [United States of America] 6.3 [Horror Thriller]} {114587 1313: Hercules Unbound! 2012-07-01 00:00:00 +0000 UTC [David DeCoteau] [Geoff Ward Lance Leonhardt Laurene Landon Priyom Haider Tyler P. Scott] [United States of America] 2.4 [Action Fantasy]} {48015 13Hrs 2010-08-28 00:00:00 +0000 UTC [Jonathan Glendening] [Gemma Atkinson Isabella Calthorpe Joshua Bowman Peter Gadiot Tom Felton] [United Kingdom] 4.4 [Action Horror]} {407806 13th 2016-10-07 00:00:00 +0000 UTC [Ava DuVernay] [Angela Davis Cory Booker Henry Louis Gates Jelani Cobb Jr. Michelle Alexander] [United States of America] 7.9 [Documentary]} {34179 14 Blades 2010-02-04 00:00:00 +0000 UTC [Daniel Lee] [Donnie Yen Chi-Tan Kate Tsui Tsz-Shan Sammo Hung Kam-Bo Wu Chun Zhao Wei] [China Hong Kong Singapore] 6.3 [Action Drama Thriller]} {534235 14 Cameras 2018-07-13 00:00:00 +0000 UTC [Scott Hussion Seth Fuller] [Chelsea Edmundson Kodi Lane Lora Martinez-Cunningham Neville Archambault Zach Dulin] [United States of America] 5 [Horror Thriller]} {890825 14 Peaks: Nothing Is Impossible 2021-11-12 00:00:00 +0000 UTC [Torquil Jones] [Conrad Anker Jimmy Chin Klára Kolouchová Nirmal Purja Reinhold Messner] [United Kingdom United States of America] 7.3 [Documentary]} {325388 14+ 2015-02-07 00:00:00 +0000 UTC [Andrey Zaytsev] [Dmitriy Blokhin Gleb Kalyuzhny Irina Frolova Olga Ozollapinya Ulyana Vaskovich] [Russia] 7 [Drama Romance]} {1179558 15 Cameras 2023-10-13 00:00:00 +0000 UTC [Danny Madden] [Angela Wong Carbone Hilty Bowen James Babson Skyler Bible Will Madden] [United States of America] 6 [Horror Thriller]} {484638 15 Minutes of War 2019-01-30 00:00:00 +0000 UTC [Fred Grivois] [Alban Lenoir Michaël Abiteboul Olga Kurylenko Sébastien Lalanne Vincent Perez] [France] 6.818 [Action Drama History War]} {483877 15+ Coming of Age 2017-08-03 00:00:00 +0000 UTC [Atsawanai Klinaiem Naphat Chitveerapat] [Jirakit Kuariyakul Lita Janvarapa Ploynarin Sornarin Thanaset Suriyapornchaikul Yongsin Wongpanitnont] [Thailand] 4.7 [Comedy Romance]} {393945 150 Milligrams 2016-11-18 00:00:00 +0000 UTC [Emmanuelle Bercot] [Benoît Magimel Charlotte Laemmel Isabelle De Hertogh Lara Neumann Sidse Babett Knudsen] [France] 6.7 [Drama]} {40205 16 Wishes 2010-10-02 00:00:00 +0000 UTC [Peter DeLuise] [Anna Mae Wills Brenda Crichlow Debby Ryan Jean-Luc Bilodeau Karissa Tynes] [United States of America] 6.314 [Drama Family Fantasy]} {87368 17 Girls 2011-06-13 00:00:00 +0000 UTC [Delphine Coulin Muriel Coulin] [Esther Garrel Juliette Darche Louise Grinberg Roxane Duran Yara Pilartz] [France] 5.5 [Drama]} {1134754 172 Days 2023-11-23 00:00:00 +0000 UTC [Hadrah Daeng Ratu] [Abun Sungkar Amara Sophie Bryan Domani Yasmin Napper Yoriko Angeline] [Indonesia] 5.63 [Drama Romance]} {630220 18 Presents 2020-01-02 00:00:00 +0000 UTC [Francesco Amato] [Benedetta Porcaroli Edoardo Leo Marco Messeri Sara Lazzaro Vittoria Puccini] [Italy United Kingdom] 7.4 [Drama Family]} {1018306 18 Year Old Hwa-jin's Crazy Sex 2020-12-07 00:00:00 +0000 UTC [Choi Jin-chul] [Hwa Jin Jung In Tae Bong] [South Korea] 3.667 [Romance]} {852363 18 Year Old Model Rika's Fancy Walk 2020-09-28 00:00:00 +0000 UTC [Choi Jin-chul] [Joo Ah Rika Si Hoo] [South Korea] 5 [Drama Romance]} {1016113 18 Year Old Muscle Queen Seong-hye's Sex Scandal 2021-10-01 00:00:00 +0000 UTC [Choi Jin-chul] [James Jung In Sung Hye] [South Korea] 8 [Romance]} {730497 18 Year Old Seungha's Easy Piece of Cake 2020-07-30 00:00:00 +0000 UTC [Choi Jin-chul] [Chul Jin Gil Dong Seung Ha] [South Korea] 7.2 [Drama Romance]} {746375 18 Year Old Seungha's Sense Game 2020-09-22 00:00:00 +0000 UTC [Choi Jin-chul] [Sang Woo Seo Won Seung Ha] [South Korea] 5 [Drama Romance]}]`,
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

	do := func(r benchflix.Repository, c Case) {
		movies, err := r.Query(ctx, c.Query)
		if err != nil {
			b.Fatal(reflect.TypeOf(r), err)
		}

		if fmt.Sprint(movies) != c.Result {
			b.Fatal(reflect.TypeOf(r), movies)
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
