// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// CountriesColumns holds the columns for the "countries" table.
	CountriesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true},
	}
	// CountriesTable holds the schema information for the "countries" table.
	CountriesTable = &schema.Table{
		Name:       "countries",
		Columns:    CountriesColumns,
		PrimaryKey: []*schema.Column{CountriesColumns[0]},
	}
	// GenresColumns holds the columns for the "genres" table.
	GenresColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true},
	}
	// GenresTable holds the schema information for the "genres" table.
	GenresTable = &schema.Table{
		Name:       "genres",
		Columns:    GenresColumns,
		PrimaryKey: []*schema.Column{GenresColumns[0]},
	}
	// MoviesColumns holds the columns for the "movies" table.
	MoviesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "title", Type: field.TypeString},
		{Name: "added_at", Type: field.TypeTime},
		{Name: "rating", Type: field.TypeFloat64},
	}
	// MoviesTable holds the schema information for the "movies" table.
	MoviesTable = &schema.Table{
		Name:       "movies",
		Columns:    MoviesColumns,
		PrimaryKey: []*schema.Column{MoviesColumns[0]},
	}
	// PeopleColumns holds the columns for the "people" table.
	PeopleColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true},
	}
	// PeopleTable holds the schema information for the "people" table.
	PeopleTable = &schema.Table{
		Name:       "people",
		Columns:    PeopleColumns,
		PrimaryKey: []*schema.Column{PeopleColumns[0]},
	}
	// MovieDirectorsColumns holds the columns for the "movie_directors" table.
	MovieDirectorsColumns = []*schema.Column{
		{Name: "movie_id", Type: field.TypeInt64},
		{Name: "person_id", Type: field.TypeInt64},
	}
	// MovieDirectorsTable holds the schema information for the "movie_directors" table.
	MovieDirectorsTable = &schema.Table{
		Name:       "movie_directors",
		Columns:    MovieDirectorsColumns,
		PrimaryKey: []*schema.Column{MovieDirectorsColumns[0], MovieDirectorsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "movie_directors_movie_id",
				Columns:    []*schema.Column{MovieDirectorsColumns[0]},
				RefColumns: []*schema.Column{MoviesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "movie_directors_person_id",
				Columns:    []*schema.Column{MovieDirectorsColumns[1]},
				RefColumns: []*schema.Column{PeopleColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// MovieActorsColumns holds the columns for the "movie_actors" table.
	MovieActorsColumns = []*schema.Column{
		{Name: "movie_id", Type: field.TypeInt64},
		{Name: "person_id", Type: field.TypeInt64},
	}
	// MovieActorsTable holds the schema information for the "movie_actors" table.
	MovieActorsTable = &schema.Table{
		Name:       "movie_actors",
		Columns:    MovieActorsColumns,
		PrimaryKey: []*schema.Column{MovieActorsColumns[0], MovieActorsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "movie_actors_movie_id",
				Columns:    []*schema.Column{MovieActorsColumns[0]},
				RefColumns: []*schema.Column{MoviesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "movie_actors_person_id",
				Columns:    []*schema.Column{MovieActorsColumns[1]},
				RefColumns: []*schema.Column{PeopleColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// MovieCountriesColumns holds the columns for the "movie_countries" table.
	MovieCountriesColumns = []*schema.Column{
		{Name: "movie_id", Type: field.TypeInt64},
		{Name: "country_id", Type: field.TypeInt64},
	}
	// MovieCountriesTable holds the schema information for the "movie_countries" table.
	MovieCountriesTable = &schema.Table{
		Name:       "movie_countries",
		Columns:    MovieCountriesColumns,
		PrimaryKey: []*schema.Column{MovieCountriesColumns[0], MovieCountriesColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "movie_countries_movie_id",
				Columns:    []*schema.Column{MovieCountriesColumns[0]},
				RefColumns: []*schema.Column{MoviesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "movie_countries_country_id",
				Columns:    []*schema.Column{MovieCountriesColumns[1]},
				RefColumns: []*schema.Column{CountriesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// MovieGenresColumns holds the columns for the "movie_genres" table.
	MovieGenresColumns = []*schema.Column{
		{Name: "movie_id", Type: field.TypeInt64},
		{Name: "genre_id", Type: field.TypeInt64},
	}
	// MovieGenresTable holds the schema information for the "movie_genres" table.
	MovieGenresTable = &schema.Table{
		Name:       "movie_genres",
		Columns:    MovieGenresColumns,
		PrimaryKey: []*schema.Column{MovieGenresColumns[0], MovieGenresColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "movie_genres_movie_id",
				Columns:    []*schema.Column{MovieGenresColumns[0]},
				RefColumns: []*schema.Column{MoviesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "movie_genres_genre_id",
				Columns:    []*schema.Column{MovieGenresColumns[1]},
				RefColumns: []*schema.Column{GenresColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		CountriesTable,
		GenresTable,
		MoviesTable,
		PeopleTable,
		MovieDirectorsTable,
		MovieActorsTable,
		MovieCountriesTable,
		MovieGenresTable,
	}
)

func init() {
	PeopleTable.Annotation = &entsql.Annotation{
		Table: "people",
	}
	MovieDirectorsTable.ForeignKeys[0].RefTable = MoviesTable
	MovieDirectorsTable.ForeignKeys[1].RefTable = PeopleTable
	MovieActorsTable.ForeignKeys[0].RefTable = MoviesTable
	MovieActorsTable.ForeignKeys[1].RefTable = PeopleTable
	MovieCountriesTable.ForeignKeys[0].RefTable = MoviesTable
	MovieCountriesTable.ForeignKeys[1].RefTable = CountriesTable
	MovieGenresTable.ForeignKeys[0].RefTable = MoviesTable
	MovieGenresTable.ForeignKeys[1].RefTable = GenresTable
}
