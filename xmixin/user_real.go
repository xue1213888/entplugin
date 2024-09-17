package xmixin

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type UserReal struct {
	mixin.Schema
}

func (UserReal) Fields() []ent.Field {
	return []ent.Field{
		field.String("real_name").Comment("姓名").Default(""),
		field.String("id_card").MaxLen(18).Comment("身份证号").Default("").SchemaType(map[string]string{dialect.MySQL: "char(18)"}).Annotations(entsql.Annotation{Size: 18}),
		field.Enum("real_sex").GoType(UserInfoSex("")).Comment("性别").Default("unknown"),
		field.Time("real_birthday").Optional().Comment("生日"),
		field.Uint8("real_province").Default(0).Comment("省级编码"),
		field.Uint8("real_city").Default(0).Comment("市级编码"),
		field.Uint8("real_district").Default(0).Comment("区级编码"),
	}
}

func (UserReal) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
				if mutation.Op().Is(ent.OpCreate | ent.OpUpdate | ent.OpUpdateOne) {
					// 1. 从身份中提取数据
					idCard, ok := mutation.Field("id_card")
					if ok && idCard != nil && idCard.(string) != "" {
						id := idCard.(string)
						if len(id) != 18 {
							return nil, fmt.Errorf("id card length is not 18")
						}
						ok, err := regexp.MatchString("^[1-9][0-9]{5}(19|20)[0-9]{2}(0[1-9]|1[012])(0[1-9]|[12][0-9]|3[01])[0-9]{3}[0-9xX]$", id)
						if err != nil {
							return nil, fmt.Errorf("id card format error: %v", err)
						}
						if !ok {
							return nil, fmt.Errorf("id card format error")
						}
						switch id[17] {
						case '0', '2', '4', '6', '8':
							err = mutation.SetField("real_sex", UserInfoSexGirl)
						default:
							err = mutation.SetField("real_sex", UserInfoSexBoy)
						}
						if err != nil {
							return nil, fmt.Errorf("id set sex failed: %v", err)
						}
						birthday := id[6:14]
						birth, err := time.Parse("20060102", birthday)
						if err != nil {
							return nil, fmt.Errorf("parse birthday failed: %v", err)
						}
						err = mutation.SetField("real_birthday", birth)
						if err != nil {
							return nil, fmt.Errorf("id set birthday failed: %v", err)
						}
						province := id[0:2]
						city := id[2:4]
						district := id[4:6]
						provinceID, _ := strconv.Atoi(province)
						cityID, _ := strconv.Atoi(city)
						districtID, _ := strconv.Atoi(district)
						err = mutation.SetField("real_province", uint8(provinceID))
						if err != nil {
							return nil, fmt.Errorf("id set province failed: %v", err)
						}
						err = mutation.SetField("real_city", uint8(cityID))
						if err != nil {
							return nil, fmt.Errorf("id set city failed: %v", err)
						}
						err = mutation.SetField("real_district", uint8(districtID))
						if err != nil {
							return nil, fmt.Errorf("id set district failed: %v", err)
						}
					}
				}
				return next.Mutate(ctx, mutation)
			})
		},
	}
}
