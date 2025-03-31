package schema

import (
	"entgo.io/ent"
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
		edge.To("directors", Person.Type),
		edge.To("actors", Person.Type),
		edge.To("countries", Country.Type),
		edge.To("genres", Genre.Type),
	}
}
