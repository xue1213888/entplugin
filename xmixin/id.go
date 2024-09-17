package xmixin

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/xue1213888/entplugin/snowflakeid"
)

// ID 是一个雪花算法ID生成器，在一个项目中使用时安全合法的
// 分布式部署的时候需要使用 snowflakeid.SetNode(nodeId) 来切换节点ID
type ID struct {
	mixin.Schema
}

func (ID) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
	}
}

// Fields of the ID.
func (ID) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").DefaultFunc(func() int64 {
			return snowflakeid.ID()
		}).Immutable().Unique().Comment("雪花ID"),
	}
}

// Edges of the ID.
func (ID) Edges() []ent.Edge {
	return nil
}
