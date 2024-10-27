package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// TimeMixin holds the schema definition for the TimeMixin entity.
type TimeMixin struct {
	mixin.Schema
}

// Fields of the TimeMixin.
func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("create_time").
			Default(time.Now).
			Immutable(),
		field.Time("update_time").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the TimeMixin.
func (TimeMixin) Edges() []ent.Edge {
	return nil
}
