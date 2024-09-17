package xmixin

import (
	"context"
	"entgo.io/ent"
	"fmt"
	"regexp"
)

func hookVerifyRegexpFunc(field string, reg *regexp.Regexp) func(mutator ent.Mutator) ent.Mutator {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
			if mutation.Op().Is(ent.OpCreate | ent.OpUpdate | ent.OpUpdateOne) {
				// 1. 从身份中提取数据
				f, ok := mutation.Field(field)
				if ok && f != nil && f.(string) != "" {
					fv := f.(string)
					// 2. 验证邮箱格式
					ok := reg.MatchString(fv)
					if !ok {
						return nil, fmt.Errorf("field format error")
					}
				}
			}
			return next.Mutate(ctx, mutation)
		})
	}
}
