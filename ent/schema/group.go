package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Group struct {
	ent.Schema
}

// Mixins of the Groups.
func (Group) Mixin() []ent.Mixin {
	return []ent.Mixin{
		ContactMixin{},
		TimeMixin{},
	}
}

// Fields of the Groups.
func (Group) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("group_code").
			MaxLen(24).
			NotEmpty(),
		field.UUID("account_id", uuid.UUID{}).
			Unique(),
		field.String("address").
			MaxLen(255).
			NotEmpty(),
		field.Bool("is_active").
			Default(true),
	}
}

// Edges of the Groups.
func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
		edge.From("account", Account.Type).
			Ref("groups").
			Unique().
			Required().
			Field("account_id"),
	}
}
