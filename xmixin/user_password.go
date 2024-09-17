package xmixin

import (
	"context"
	"crypto/md5"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"fmt"
	"regexp"
)

type UserPassword struct {
	mixin.Schema
}

func (UserPassword) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").MinLen(10).MaxLen(32).Match(regexp.MustCompile(`^[a-zA-Z0-9_]{10,32}$`)).Comment("用户登录账号(数字字母下划线,10~32)").Unique().Optional(),
		field.String("password").MaxLen(32).Comment("用户密码").SchemaType(map[string]string{dialect.MySQL: "char(32)"}).Default(""),
	}
}

func (UserPassword) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
				if mutation.Op().Is(ent.OpCreate | ent.OpUpdate | ent.OpUpdateOne) {
					passwordField, ok := mutation.Field("password")
					if ok && passwordField != nil && passwordField.(string) != "" {
						pass := fmt.Sprintf("%x", md5.Sum([]byte(passwordField.(string))))
						err := mutation.SetField("password", pass)
						if err != nil {
							return nil, fmt.Errorf("set password failed: %v", err)
						}
					}
				}
				return next.Mutate(ctx, mutation)
			})
		},
	}
}

func (UserPassword) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username", "password"),
	}
}
