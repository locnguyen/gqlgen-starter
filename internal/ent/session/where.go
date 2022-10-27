// Code generated by ent, DO NOT EDIT.

package session

import (
	"gqlgen-starter/internal/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		v := make([]any, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		v := make([]any, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Sid applies equality check predicate on the "sid" field. It's identical to SidEQ.
func Sid(v string) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSid), v))
	})
}

// Expiry applies equality check predicate on the "expiry" field. It's identical to ExpiryEQ.
func Expiry(v time.Time) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldExpiry), v))
	})
}

// Deleted applies equality check predicate on the "deleted" field. It's identical to DeletedEQ.
func Deleted(v bool) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDeleted), v))
	})
}

// SidEQ applies the EQ predicate on the "sid" field.
func SidEQ(v string) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSid), v))
	})
}

// SidNEQ applies the NEQ predicate on the "sid" field.
func SidNEQ(v string) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSid), v))
	})
}

// SidIn applies the In predicate on the "sid" field.
func SidIn(vs ...string) predicate.Session {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldSid), v...))
	})
}

// SidNotIn applies the NotIn predicate on the "sid" field.
func SidNotIn(vs ...string) predicate.Session {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldSid), v...))
	})
}

// SidGT applies the GT predicate on the "sid" field.
func SidGT(v string) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSid), v))
	})
}

// SidGTE applies the GTE predicate on the "sid" field.
func SidGTE(v string) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSid), v))
	})
}

// SidLT applies the LT predicate on the "sid" field.
func SidLT(v string) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSid), v))
	})
}

// SidLTE applies the LTE predicate on the "sid" field.
func SidLTE(v string) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSid), v))
	})
}

// SidContains applies the Contains predicate on the "sid" field.
func SidContains(v string) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldSid), v))
	})
}

// SidHasPrefix applies the HasPrefix predicate on the "sid" field.
func SidHasPrefix(v string) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldSid), v))
	})
}

// SidHasSuffix applies the HasSuffix predicate on the "sid" field.
func SidHasSuffix(v string) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldSid), v))
	})
}

// SidEqualFold applies the EqualFold predicate on the "sid" field.
func SidEqualFold(v string) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldSid), v))
	})
}

// SidContainsFold applies the ContainsFold predicate on the "sid" field.
func SidContainsFold(v string) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldSid), v))
	})
}

// ExpiryEQ applies the EQ predicate on the "expiry" field.
func ExpiryEQ(v time.Time) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldExpiry), v))
	})
}

// ExpiryNEQ applies the NEQ predicate on the "expiry" field.
func ExpiryNEQ(v time.Time) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldExpiry), v))
	})
}

// ExpiryIn applies the In predicate on the "expiry" field.
func ExpiryIn(vs ...time.Time) predicate.Session {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldExpiry), v...))
	})
}

// ExpiryNotIn applies the NotIn predicate on the "expiry" field.
func ExpiryNotIn(vs ...time.Time) predicate.Session {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldExpiry), v...))
	})
}

// ExpiryGT applies the GT predicate on the "expiry" field.
func ExpiryGT(v time.Time) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldExpiry), v))
	})
}

// ExpiryGTE applies the GTE predicate on the "expiry" field.
func ExpiryGTE(v time.Time) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldExpiry), v))
	})
}

// ExpiryLT applies the LT predicate on the "expiry" field.
func ExpiryLT(v time.Time) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldExpiry), v))
	})
}

// ExpiryLTE applies the LTE predicate on the "expiry" field.
func ExpiryLTE(v time.Time) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldExpiry), v))
	})
}

// DeletedEQ applies the EQ predicate on the "deleted" field.
func DeletedEQ(v bool) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDeleted), v))
	})
}

// DeletedNEQ applies the NEQ predicate on the "deleted" field.
func DeletedNEQ(v bool) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldDeleted), v))
	})
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v Type) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldType), v))
	})
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v Type) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldType), v))
	})
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...Type) predicate.Session {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldType), v...))
	})
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...Type) predicate.Session {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Session(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldType), v...))
	})
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(UserTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(UserInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Session) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Session) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Session) predicate.Session {
	return predicate.Session(func(s *sql.Selector) {
		p(s.Not())
	})
}
