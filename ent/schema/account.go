package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Account struct {
	ent.Schema
}

// Mixins of the Accounts.
func (Account) Mixin() []ent.Mixin {
	return []ent.Mixin{
		ContactMixin{},
		TimeMixin{},
	}
}

// Fields of the Accounts.
func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("api_key").
			MaxLen(255).
			NotEmpty(),
		field.Bool("is_active").
			Default(true),
	}
}

// Edges of the Accounts.
func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("groups", Group.Type),
	}
}
