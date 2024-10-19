package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").
			MaxLen(255).
			NotEmpty().
			Unique(),
		field.String("name").
			MaxLen(255),
		field.String("avatar").
			MaxLen(255),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
