package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type User struct {
	ent.Schema
}

// Mixins of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		ContactMixin{},
		TimeMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.UUID("group_id", uuid.UUID{}).
			Unique(),
		field.Bool("is_active").
			Default(true),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").
			Unique(),
		index.Fields("uid").
			Unique(),
		index.Fields("group_id").
			Unique(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("group", Group.Type).
			Ref("users").
			Unique().
			Field("group_id").
			Required(),
	}
}
