package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Accounts holds the schema definition for the Accounts entity.
type Accounts struct {
	ent.Schema
}

// Fields of the Accounts.
func (Accounts) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(255).
			NotEmpty(),
		field.String("email").
			MaxLen(255).
			NotEmpty(),
		field.String("api_key").
			MaxLen(255).
			NotEmpty(),
		field.String("uid").
			MaxLen(255).
			NotEmpty(),
		field.Bool("is_active").
			Default(true),
		field.Time("creation_time").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Accounts.
func (Accounts) Edges() []ent.Edge {
	return nil
}
