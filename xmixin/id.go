package xmixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/xue1213888/entplugin/snowflakeid"
)

// ID holds the schema definition for the ID entity.
type ID struct {
	mixin.Schema
}

// Fields of the ID.
func (ID) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").DefaultFunc(func() int64 {
			return snowflakeid.ID()
		}).Immutable().Immutable(),
	}
}

// Edges of the ID.
func (ID) Edges() []ent.Edge {
	return nil
}
