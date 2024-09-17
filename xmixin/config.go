package xmixin

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Config
// 用于标注模型是否需要自定义配置，也可以作为中间集中管理，需要额外设计
type Config struct {
	mixin.Schema
}

func (Config) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
	}
}

func (Config) Fields() []ent.Field {
	return []ent.Field{
		field.String("xkey").MaxLen(255).Comment("配置key"),
		field.Any("xval").Comment("配置value"),
	}
}
