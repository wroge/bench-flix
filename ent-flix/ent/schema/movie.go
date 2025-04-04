package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Movie holds the schema definition for the Movie entity.
type Movie struct {
	ent.Schema
}

// Fields of the Movie.
func (Movie) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.String("title"),
		field.Time("added_at"),
		field.Float("rating"),
	}
}

// Edges of the Movie.
func (Movie) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("directors", Person.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("actors", Person.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("countries", Country.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("genres", Genre.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
