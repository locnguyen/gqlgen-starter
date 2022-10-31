package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("user_id"),
		field.String("sid").
			Unique().
			NotEmpty(),
		field.Time("expiry"),
		field.Bool("deleted").
			Default(false),
		field.Enum("type").
			Values("general", "single_use"),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Required().
			Field("user_id").
			Ref("sessions").
			Unique(),
	}
}

func (Session) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
