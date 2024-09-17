package xmixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type UserIntegral struct {
	mixin.Schema
}

func (UserIntegral) Fields() []ent.Field {
	return []ent.Field{
		field.String("flag").Default("unknown").Comment("积分标记").MaxLen(32),
		field.Uint64("integral").Default(0).Comment("用户积分"),
	}
}
