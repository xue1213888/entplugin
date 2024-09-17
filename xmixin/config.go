package xmixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Config
// 用于标注模型是否需要自定义配置，也可以作为中间集中管理，需要额外设计
type Config struct {
	mixin.Schema
}

func (Config) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").MaxLen(255).Comment("配置key"),
		field.Any("val").Comment("配置value"),
	}
}
