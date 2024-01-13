package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type BaseMixin struct {
	mixin.Schema
}

func (BaseMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			StructTag(`json:"id,omitempty,string"`).
			Unique().
			Immutable().
			Annotations(entsql.DefaultExpr("generate_id()")),
	}
}
