package xmixin

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/privacy"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Audit
// 用来标注模型审核状态，需要实现审核状态的实体
type Audit struct {
	mixin.Schema
}

var AuditReasonMap = map[string]string{
	"unknown":  "未知",
	"pending":  "待审核",
	"waiting":  "审核中",
	"approved": "已通过",
	"rejected": "已拒绝",
	"removed":  "下架锁定",
}

func (Audit) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("audit").Values("unknown", "pending", "waiting", "approved", "rejected", "removed").Comment("审核状态"),
		field.String("reason").MaxLen(255).Comment("审核原因"),
	}
}

func (Audit) Policy() ent.Policy {
	return privacy.Policy{
		Query: privacy.QueryPolicy{
			privacy.ContextQueryMutationRule(func(ctx context.Context) error {
				return nil
			}),
		},
		Mutation: nil,
	}
}
