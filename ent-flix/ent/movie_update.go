// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/wroge/bench-flix/ent-flix/ent/country"
	"github.com/wroge/bench-flix/ent-flix/ent/genre"
	"github.com/wroge/bench-flix/ent-flix/ent/movie"
	"github.com/wroge/bench-flix/ent-flix/ent/person"
	"github.com/wroge/bench-flix/ent-flix/ent/predicate"
)

// MovieUpdate is the builder for updating Movie entities.
type MovieUpdate struct {
	config
	hooks    []Hook
	mutation *MovieMutation
}

// Where appends a list predicates to the MovieUpdate builder.
func (mu *MovieUpdate) Where(ps ...predicate.Movie) *MovieUpdate {
	mu.mutation.Where(ps...)
	return mu
}

// SetTitle sets the "title" field.
func (mu *MovieUpdate) SetTitle(s string) *MovieUpdate {
	mu.mutation.SetTitle(s)
	return mu
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (mu *MovieUpdate) SetNillableTitle(s *string) *MovieUpdate {
	if s != nil {
		mu.SetTitle(*s)
	}
	return mu
}

// SetAddedAt sets the "added_at" field.
func (mu *MovieUpdate) SetAddedAt(t time.Time) *MovieUpdate {
	mu.mutation.SetAddedAt(t)
	return mu
}

// SetNillableAddedAt sets the "added_at" field if the given value is not nil.
func (mu *MovieUpdate) SetNillableAddedAt(t *time.Time) *MovieUpdate {
	if t != nil {
		mu.SetAddedAt(*t)
	}
	return mu
}

// SetRating sets the "rating" field.
func (mu *MovieUpdate) SetRating(f float64) *MovieUpdate {
	mu.mutation.ResetRating()
	mu.mutation.SetRating(f)
	return mu
}

// SetNillableRating sets the "rating" field if the given value is not nil.
func (mu *MovieUpdate) SetNillableRating(f *float64) *MovieUpdate {
	if f != nil {
		mu.SetRating(*f)
	}
	return mu
}

// AddRating adds f to the "rating" field.
func (mu *MovieUpdate) AddRating(f float64) *MovieUpdate {
	mu.mutation.AddRating(f)
	return mu
}

// AddDirectorIDs adds the "directors" edge to the Person entity by IDs.
func (mu *MovieUpdate) AddDirectorIDs(ids ...int64) *MovieUpdate {
	mu.mutation.AddDirectorIDs(ids...)
	return mu
}

// AddDirectors adds the "directors" edges to the Person entity.
func (mu *MovieUpdate) AddDirectors(p ...*Person) *MovieUpdate {
	ids := make([]int64, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return mu.AddDirectorIDs(ids...)
}

// AddActorIDs adds the "actors" edge to the Person entity by IDs.
func (mu *MovieUpdate) AddActorIDs(ids ...int64) *MovieUpdate {
	mu.mutation.AddActorIDs(ids...)
	return mu
}

// AddActors adds the "actors" edges to the Person entity.
func (mu *MovieUpdate) AddActors(p ...*Person) *MovieUpdate {
	ids := make([]int64, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return mu.AddActorIDs(ids...)
}

// AddCountryIDs adds the "countries" edge to the Country entity by IDs.
func (mu *MovieUpdate) AddCountryIDs(ids ...int64) *MovieUpdate {
	mu.mutation.AddCountryIDs(ids...)
	return mu
}

// AddCountries adds the "countries" edges to the Country entity.
func (mu *MovieUpdate) AddCountries(c ...*Country) *MovieUpdate {
	ids := make([]int64, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return mu.AddCountryIDs(ids...)
}

// AddGenreIDs adds the "genres" edge to the Genre entity by IDs.
func (mu *MovieUpdate) AddGenreIDs(ids ...int64) *MovieUpdate {
	mu.mutation.AddGenreIDs(ids...)
	return mu
}

// AddGenres adds the "genres" edges to the Genre entity.
func (mu *MovieUpdate) AddGenres(g ...*Genre) *MovieUpdate {
	ids := make([]int64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return mu.AddGenreIDs(ids...)
}

// Mutation returns the MovieMutation object of the builder.
func (mu *MovieUpdate) Mutation() *MovieMutation {
	return mu.mutation
}

// ClearDirectors clears all "directors" edges to the Person entity.
func (mu *MovieUpdate) ClearDirectors() *MovieUpdate {
	mu.mutation.ClearDirectors()
	return mu
}

// RemoveDirectorIDs removes the "directors" edge to Person entities by IDs.
func (mu *MovieUpdate) RemoveDirectorIDs(ids ...int64) *MovieUpdate {
	mu.mutation.RemoveDirectorIDs(ids...)
	return mu
}

// RemoveDirectors removes "directors" edges to Person entities.
func (mu *MovieUpdate) RemoveDirectors(p ...*Person) *MovieUpdate {
	ids := make([]int64, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return mu.RemoveDirectorIDs(ids...)
}

// ClearActors clears all "actors" edges to the Person entity.
func (mu *MovieUpdate) ClearActors() *MovieUpdate {
	mu.mutation.ClearActors()
	return mu
}

// RemoveActorIDs removes the "actors" edge to Person entities by IDs.
func (mu *MovieUpdate) RemoveActorIDs(ids ...int64) *MovieUpdate {
	mu.mutation.RemoveActorIDs(ids...)
	return mu
}

// RemoveActors removes "actors" edges to Person entities.
func (mu *MovieUpdate) RemoveActors(p ...*Person) *MovieUpdate {
	ids := make([]int64, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return mu.RemoveActorIDs(ids...)
}

// ClearCountries clears all "countries" edges to the Country entity.
func (mu *MovieUpdate) ClearCountries() *MovieUpdate {
	mu.mutation.ClearCountries()
	return mu
}

// RemoveCountryIDs removes the "countries" edge to Country entities by IDs.
func (mu *MovieUpdate) RemoveCountryIDs(ids ...int64) *MovieUpdate {
	mu.mutation.RemoveCountryIDs(ids...)
	return mu
}

// RemoveCountries removes "countries" edges to Country entities.
func (mu *MovieUpdate) RemoveCountries(c ...*Country) *MovieUpdate {
	ids := make([]int64, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return mu.RemoveCountryIDs(ids...)
}

// ClearGenres clears all "genres" edges to the Genre entity.
func (mu *MovieUpdate) ClearGenres() *MovieUpdate {
	mu.mutation.ClearGenres()
	return mu
}

// RemoveGenreIDs removes the "genres" edge to Genre entities by IDs.
func (mu *MovieUpdate) RemoveGenreIDs(ids ...int64) *MovieUpdate {
	mu.mutation.RemoveGenreIDs(ids...)
	return mu
}

// RemoveGenres removes "genres" edges to Genre entities.
func (mu *MovieUpdate) RemoveGenres(g ...*Genre) *MovieUpdate {
	ids := make([]int64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return mu.RemoveGenreIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (mu *MovieUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, mu.sqlSave, mu.mutation, mu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (mu *MovieUpdate) SaveX(ctx context.Context) int {
	affected, err := mu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (mu *MovieUpdate) Exec(ctx context.Context) error {
	_, err := mu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mu *MovieUpdate) ExecX(ctx context.Context) {
	if err := mu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (mu *MovieUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(movie.Table, movie.Columns, sqlgraph.NewFieldSpec(movie.FieldID, field.TypeInt64))
	if ps := mu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := mu.mutation.Title(); ok {
		_spec.SetField(movie.FieldTitle, field.TypeString, value)
	}
	if value, ok := mu.mutation.AddedAt(); ok {
		_spec.SetField(movie.FieldAddedAt, field.TypeTime, value)
	}
	if value, ok := mu.mutation.Rating(); ok {
		_spec.SetField(movie.FieldRating, field.TypeFloat64, value)
	}
	if value, ok := mu.mutation.AddedRating(); ok {
		_spec.AddField(movie.FieldRating, field.TypeFloat64, value)
	}
	if mu.mutation.DirectorsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.DirectorsTable,
			Columns: movie.DirectorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(person.FieldID, field.TypeInt64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.RemovedDirectorsIDs(); len(nodes) > 0 && !mu.mutation.DirectorsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.DirectorsTable,
			Columns: movie.DirectorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(person.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.DirectorsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.DirectorsTable,
			Columns: movie.DirectorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(person.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if mu.mutation.ActorsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.ActorsTable,
			Columns: movie.ActorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(person.FieldID, field.TypeInt64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.RemovedActorsIDs(); len(nodes) > 0 && !mu.mutation.ActorsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.ActorsTable,
			Columns: movie.ActorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(person.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.ActorsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.ActorsTable,
			Columns: movie.ActorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(person.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if mu.mutation.CountriesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.CountriesTable,
			Columns: movie.CountriesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(country.FieldID, field.TypeInt64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.RemovedCountriesIDs(); len(nodes) > 0 && !mu.mutation.CountriesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.CountriesTable,
			Columns: movie.CountriesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(country.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.CountriesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.CountriesTable,
			Columns: movie.CountriesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(country.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if mu.mutation.GenresCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.GenresTable,
			Columns: movie.GenresPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(genre.FieldID, field.TypeInt64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.RemovedGenresIDs(); len(nodes) > 0 && !mu.mutation.GenresCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.GenresTable,
			Columns: movie.GenresPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(genre.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.GenresIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.GenresTable,
			Columns: movie.GenresPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(genre.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, mu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{movie.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	mu.mutation.done = true
	return n, nil
}

// MovieUpdateOne is the builder for updating a single Movie entity.
type MovieUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *MovieMutation
}

// SetTitle sets the "title" field.
func (muo *MovieUpdateOne) SetTitle(s string) *MovieUpdateOne {
	muo.mutation.SetTitle(s)
	return muo
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (muo *MovieUpdateOne) SetNillableTitle(s *string) *MovieUpdateOne {
	if s != nil {
		muo.SetTitle(*s)
	}
	return muo
}

// SetAddedAt sets the "added_at" field.
func (muo *MovieUpdateOne) SetAddedAt(t time.Time) *MovieUpdateOne {
	muo.mutation.SetAddedAt(t)
	return muo
}

// SetNillableAddedAt sets the "added_at" field if the given value is not nil.
func (muo *MovieUpdateOne) SetNillableAddedAt(t *time.Time) *MovieUpdateOne {
	if t != nil {
		muo.SetAddedAt(*t)
	}
	return muo
}

// SetRating sets the "rating" field.
func (muo *MovieUpdateOne) SetRating(f float64) *MovieUpdateOne {
	muo.mutation.ResetRating()
	muo.mutation.SetRating(f)
	return muo
}

// SetNillableRating sets the "rating" field if the given value is not nil.
func (muo *MovieUpdateOne) SetNillableRating(f *float64) *MovieUpdateOne {
	if f != nil {
		muo.SetRating(*f)
	}
	return muo
}

// AddRating adds f to the "rating" field.
func (muo *MovieUpdateOne) AddRating(f float64) *MovieUpdateOne {
	muo.mutation.AddRating(f)
	return muo
}

// AddDirectorIDs adds the "directors" edge to the Person entity by IDs.
func (muo *MovieUpdateOne) AddDirectorIDs(ids ...int64) *MovieUpdateOne {
	muo.mutation.AddDirectorIDs(ids...)
	return muo
}

// AddDirectors adds the "directors" edges to the Person entity.
func (muo *MovieUpdateOne) AddDirectors(p ...*Person) *MovieUpdateOne {
	ids := make([]int64, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return muo.AddDirectorIDs(ids...)
}

// AddActorIDs adds the "actors" edge to the Person entity by IDs.
func (muo *MovieUpdateOne) AddActorIDs(ids ...int64) *MovieUpdateOne {
	muo.mutation.AddActorIDs(ids...)
	return muo
}

// AddActors adds the "actors" edges to the Person entity.
func (muo *MovieUpdateOne) AddActors(p ...*Person) *MovieUpdateOne {
	ids := make([]int64, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return muo.AddActorIDs(ids...)
}

// AddCountryIDs adds the "countries" edge to the Country entity by IDs.
func (muo *MovieUpdateOne) AddCountryIDs(ids ...int64) *MovieUpdateOne {
	muo.mutation.AddCountryIDs(ids...)
	return muo
}

// AddCountries adds the "countries" edges to the Country entity.
func (muo *MovieUpdateOne) AddCountries(c ...*Country) *MovieUpdateOne {
	ids := make([]int64, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return muo.AddCountryIDs(ids...)
}

// AddGenreIDs adds the "genres" edge to the Genre entity by IDs.
func (muo *MovieUpdateOne) AddGenreIDs(ids ...int64) *MovieUpdateOne {
	muo.mutation.AddGenreIDs(ids...)
	return muo
}

// AddGenres adds the "genres" edges to the Genre entity.
func (muo *MovieUpdateOne) AddGenres(g ...*Genre) *MovieUpdateOne {
	ids := make([]int64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return muo.AddGenreIDs(ids...)
}

// Mutation returns the MovieMutation object of the builder.
func (muo *MovieUpdateOne) Mutation() *MovieMutation {
	return muo.mutation
}

// ClearDirectors clears all "directors" edges to the Person entity.
func (muo *MovieUpdateOne) ClearDirectors() *MovieUpdateOne {
	muo.mutation.ClearDirectors()
	return muo
}

// RemoveDirectorIDs removes the "directors" edge to Person entities by IDs.
func (muo *MovieUpdateOne) RemoveDirectorIDs(ids ...int64) *MovieUpdateOne {
	muo.mutation.RemoveDirectorIDs(ids...)
	return muo
}

// RemoveDirectors removes "directors" edges to Person entities.
func (muo *MovieUpdateOne) RemoveDirectors(p ...*Person) *MovieUpdateOne {
	ids := make([]int64, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return muo.RemoveDirectorIDs(ids...)
}

// ClearActors clears all "actors" edges to the Person entity.
func (muo *MovieUpdateOne) ClearActors() *MovieUpdateOne {
	muo.mutation.ClearActors()
	return muo
}

// RemoveActorIDs removes the "actors" edge to Person entities by IDs.
func (muo *MovieUpdateOne) RemoveActorIDs(ids ...int64) *MovieUpdateOne {
	muo.mutation.RemoveActorIDs(ids...)
	return muo
}

// RemoveActors removes "actors" edges to Person entities.
func (muo *MovieUpdateOne) RemoveActors(p ...*Person) *MovieUpdateOne {
	ids := make([]int64, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return muo.RemoveActorIDs(ids...)
}

// ClearCountries clears all "countries" edges to the Country entity.
func (muo *MovieUpdateOne) ClearCountries() *MovieUpdateOne {
	muo.mutation.ClearCountries()
	return muo
}

// RemoveCountryIDs removes the "countries" edge to Country entities by IDs.
func (muo *MovieUpdateOne) RemoveCountryIDs(ids ...int64) *MovieUpdateOne {
	muo.mutation.RemoveCountryIDs(ids...)
	return muo
}

// RemoveCountries removes "countries" edges to Country entities.
func (muo *MovieUpdateOne) RemoveCountries(c ...*Country) *MovieUpdateOne {
	ids := make([]int64, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return muo.RemoveCountryIDs(ids...)
}

// ClearGenres clears all "genres" edges to the Genre entity.
func (muo *MovieUpdateOne) ClearGenres() *MovieUpdateOne {
	muo.mutation.ClearGenres()
	return muo
}

// RemoveGenreIDs removes the "genres" edge to Genre entities by IDs.
func (muo *MovieUpdateOne) RemoveGenreIDs(ids ...int64) *MovieUpdateOne {
	muo.mutation.RemoveGenreIDs(ids...)
	return muo
}

// RemoveGenres removes "genres" edges to Genre entities.
func (muo *MovieUpdateOne) RemoveGenres(g ...*Genre) *MovieUpdateOne {
	ids := make([]int64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return muo.RemoveGenreIDs(ids...)
}

// Where appends a list predicates to the MovieUpdate builder.
func (muo *MovieUpdateOne) Where(ps ...predicate.Movie) *MovieUpdateOne {
	muo.mutation.Where(ps...)
	return muo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (muo *MovieUpdateOne) Select(field string, fields ...string) *MovieUpdateOne {
	muo.fields = append([]string{field}, fields...)
	return muo
}

// Save executes the query and returns the updated Movie entity.
func (muo *MovieUpdateOne) Save(ctx context.Context) (*Movie, error) {
	return withHooks(ctx, muo.sqlSave, muo.mutation, muo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (muo *MovieUpdateOne) SaveX(ctx context.Context) *Movie {
	node, err := muo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (muo *MovieUpdateOne) Exec(ctx context.Context) error {
	_, err := muo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (muo *MovieUpdateOne) ExecX(ctx context.Context) {
	if err := muo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (muo *MovieUpdateOne) sqlSave(ctx context.Context) (_node *Movie, err error) {
	_spec := sqlgraph.NewUpdateSpec(movie.Table, movie.Columns, sqlgraph.NewFieldSpec(movie.FieldID, field.TypeInt64))
	id, ok := muo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Movie.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := muo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, movie.FieldID)
		for _, f := range fields {
			if !movie.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != movie.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := muo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := muo.mutation.Title(); ok {
		_spec.SetField(movie.FieldTitle, field.TypeString, value)
	}
	if value, ok := muo.mutation.AddedAt(); ok {
		_spec.SetField(movie.FieldAddedAt, field.TypeTime, value)
	}
	if value, ok := muo.mutation.Rating(); ok {
		_spec.SetField(movie.FieldRating, field.TypeFloat64, value)
	}
	if value, ok := muo.mutation.AddedRating(); ok {
		_spec.AddField(movie.FieldRating, field.TypeFloat64, value)
	}
	if muo.mutation.DirectorsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.DirectorsTable,
			Columns: movie.DirectorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(person.FieldID, field.TypeInt64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.RemovedDirectorsIDs(); len(nodes) > 0 && !muo.mutation.DirectorsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.DirectorsTable,
			Columns: movie.DirectorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(person.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.DirectorsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.DirectorsTable,
			Columns: movie.DirectorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(person.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if muo.mutation.ActorsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.ActorsTable,
			Columns: movie.ActorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(person.FieldID, field.TypeInt64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.RemovedActorsIDs(); len(nodes) > 0 && !muo.mutation.ActorsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.ActorsTable,
			Columns: movie.ActorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(person.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.ActorsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.ActorsTable,
			Columns: movie.ActorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(person.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if muo.mutation.CountriesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.CountriesTable,
			Columns: movie.CountriesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(country.FieldID, field.TypeInt64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.RemovedCountriesIDs(); len(nodes) > 0 && !muo.mutation.CountriesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.CountriesTable,
			Columns: movie.CountriesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(country.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.CountriesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.CountriesTable,
			Columns: movie.CountriesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(country.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if muo.mutation.GenresCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.GenresTable,
			Columns: movie.GenresPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(genre.FieldID, field.TypeInt64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.RemovedGenresIDs(); len(nodes) > 0 && !muo.mutation.GenresCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.GenresTable,
			Columns: movie.GenresPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(genre.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.GenresIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   movie.GenresTable,
			Columns: movie.GenresPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(genre.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Movie{config: muo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, muo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{movie.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	muo.mutation.done = true
	return _node, nil
}
