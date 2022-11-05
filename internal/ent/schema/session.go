package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String("token").
			Unique().
			NotEmpty(),
		field.Bytes("data").
			NotEmpty(),
		field.Time("expiry"),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Session) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("expiry"),
	}
}
