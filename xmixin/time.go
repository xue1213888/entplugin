package xmixin

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"fmt"
	"reflect"
	"time"
)

var zeroTime time.Time

// Time 为模型提供记录创建时间、更新时间、删除时间
// 创建时间是不可变的
// 更新时间在update的时候修改
// 删除时间在delete的时候修改
type Time struct {
	mixin.Schema
	Optional bool
}

// Fields of the Time.
func (t Time) Fields() []ent.Field {
	f := make([]ent.Field, 0, 3)
	f = append(f, field.Time("created_at").Default(time.Now).Immutable())
	f = append(f, field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now))
	if t.Optional {
		f = append(f, field.Time("deleted_at").Optional())
	} else {
		f = append(f, field.Int64("deleted_at").Default(0).Comment("秒级时间戳"))
	}
	return f
}

// Edges of the Time.
func (Time) Edges() []ent.Edge {
	return nil
}

type softDeleteKey struct{}

func SkipSoftDelete(parent context.Context) context.Context {
	return context.WithValue(parent, softDeleteKey{}, true)
}

func (t Time) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.TraverseFunc(func(ctx context.Context, query ent.Query) error {
			if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
				return nil
			}
			nq, ok := query.(interface {
				WhereP(ps ...func(sql *sql.Selector))
			})
			if !ok {
				return nil
			}
			if t.Optional {
				nq.WhereP(sql.FieldIsNull(t.Fields()[2].Descriptor().Name))
			} else {
				nq.WhereP(sql.FieldEQ(t.Fields()[2].Descriptor().Name, 0))
			}
			return nil
		}),
	}
}

func (t Time) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				defer func() {
					if err := recover(); err != nil {
						fmt.Println(err)
					}
				}()
				if m.Op().Is(ent.OpDelete | ent.OpDeleteOne) {
					if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
						return next.Mutate(ctx, m)
					}
					// 需要更新deleted_at的
					mx, ok := m.(interface {
						SetOp(op ent.Op)
						SetDeletedAt(time.Time)
						WhereP(ps ...func(sql *sql.Selector))
					})
					if !ok {
						return nil, fmt.Errorf("unexpected mutation type %T", m)
					}
					if t.Optional {
						mx.WhereP(sql.FieldIsNull(t.Fields()[2].Descriptor().Name))
					} else {
						mx.WhereP(sql.FieldEQ(t.Fields()[2].Descriptor().Name, 0))
					}
					mx.SetOp(ent.OpUpdate)
					mx.SetDeletedAt(time.Now())
					client := reflect.ValueOf(m).MethodByName("Client").Call(nil)[0]
					res := client.MethodByName("Mutate").Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(m)})
					nErr := res[1].Interface()
					if nErr != nil {
						return res[0].Interface(), nErr.(error)
					}
					return res[0].Interface(), nil
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}
