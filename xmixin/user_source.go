package xmixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type UserSource struct {
	mixin.Schema

	Source []string // 第一项必须是unknown否则不合法
}

func (u UserSource) Fields() []ent.Field {
	if u.Source == nil || len(u.Source) < 1 || u.Source[0] != "unknown" {
		panic("user source invalid source")
	}
	return []ent.Field{
		field.Enum("source").Values(u.Source...).Default("unknown").Comment("用户来源"),
		field.String("openid").Default("").Comment("用户的openid"),
		field.String("unionid").Default("").Comment("用户的unionid"),
	}
}

func (UserSource) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("source", "openid"),
		index.Fields("source", "unionid"),
	}
}
