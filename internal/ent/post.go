// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"gqlgen-starter/internal/ent/post"
	"gqlgen-starter/internal/ent/user"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
)

// Post is the model entity for the Post schema.
type Post struct {
	config `json:"-"`
	// ID of the ent.
	ID int64 `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// AuthorID holds the value of the "author_id" field.
	AuthorID int64 `json:"author_id,omitempty"`
	// Content holds the value of the "content" field.
	Content string `json:"content,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PostQuery when eager-loading is set.
	Edges PostEdges `json:"edges"`
}

// PostEdges holds the relations/edges for other nodes in the graph.
type PostEdges struct {
	// Author holds the value of the author edge.
	Author *User `json:"author,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// AuthorOrErr returns the Author value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PostEdges) AuthorOrErr() (*User, error) {
	if e.loadedTypes[0] {
		if e.Author == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.Author, nil
	}
	return nil, &NotLoadedError{edge: "author"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Post) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case post.FieldID, post.FieldAuthorID:
			values[i] = new(sql.NullInt64)
		case post.FieldContent:
			values[i] = new(sql.NullString)
		case post.FieldCreateTime, post.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Post", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Post fields.
func (po *Post) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case post.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			po.ID = int64(value.Int64)
		case post.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				po.CreateTime = value.Time
			}
		case post.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				po.UpdateTime = value.Time
			}
		case post.FieldAuthorID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field author_id", values[i])
			} else if value.Valid {
				po.AuthorID = value.Int64
			}
		case post.FieldContent:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field content", values[i])
			} else if value.Valid {
				po.Content = value.String
			}
		}
	}
	return nil
}

// QueryAuthor queries the "author" edge of the Post entity.
func (po *Post) QueryAuthor() *UserQuery {
	return (&PostClient{config: po.config}).QueryAuthor(po)
}

// Update returns a builder for updating this Post.
// Note that you need to call Post.Unwrap() before calling this method if this Post
// was returned from a transaction, and the transaction was committed or rolled back.
func (po *Post) Update() *PostUpdateOne {
	return (&PostClient{config: po.config}).UpdateOne(po)
}

// Unwrap unwraps the Post entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (po *Post) Unwrap() *Post {
	_tx, ok := po.config.driver.(*txDriver)
	if !ok {
		panic("ent: Post is not a transactional entity")
	}
	po.config.driver = _tx.drv
	return po
}

// String implements the fmt.Stringer.
func (po *Post) String() string {
	var builder strings.Builder
	builder.WriteString("Post(")
	builder.WriteString(fmt.Sprintf("id=%v, ", po.ID))
	builder.WriteString("create_time=")
	builder.WriteString(po.CreateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("update_time=")
	builder.WriteString(po.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("author_id=")
	builder.WriteString(fmt.Sprintf("%v", po.AuthorID))
	builder.WriteString(", ")
	builder.WriteString("content=")
	builder.WriteString(po.Content)
	builder.WriteByte(')')
	return builder.String()
}

// Posts is a parsable slice of Post.
type Posts []*Post

func (po Posts) config(cfg config) {
	for _i := range po {
		po[_i].config = cfg
	}
}