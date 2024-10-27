package schema

import (
	"errors"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/asaskevich/govalidator"
	"github.com/nyaruka/phonenumbers"
)

// ContactMixin holds the schema definition for the ContactMixin entity.
type ContactMixin struct {
	mixin.Schema
}

// Email validation.
func email(e string) error {
	if !govalidator.IsEmail(e) {
		return errors.New("email is invalid")
	}
	return nil
}

// Phone number validation.
func phone(p string) error {
	parsed, err := phonenumbers.Parse(p, "US")
	if err != nil {
		return errors.New("error parsing phone number")
	}
	if !phonenumbers.IsValidNumber(parsed) {
		return errors.New("phone number is invalid")
	}
	return nil
}

// Fields of the ContactMixin.
func (ContactMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(100).
			Nillable().
			Default(""),
		field.String("nickname").
			MaxLen(100).
			Nillable().
			Default(""),
		field.String("email").
			MaxLen(100).
			Nillable().
			Unique().
			Validate(email),
		field.String("phone_number").
			MaxLen(100).
			Nillable().
			Unique().
			Validate(phone),
		field.String("photo_url").
			MaxLen(255).
			Nillable().
			Unique(),
		field.String("uid").
			MaxLen(255).
			NotEmpty().
			Unique(),
	}
}

// Edges of the ContactMixin.
func (ContactMixin) Edges() []ent.Edge {
	return nil
}
