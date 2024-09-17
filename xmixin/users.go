package xmixin

import (
	"context"
	"crypto/md5"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type UserSex string

const (
	UserSexUnknown UserSex = "unknown"
	UserSexBoy     UserSex = "boy"
	UserSexGirl    UserSex = "girl"
)

func (UserSex) Values() []string {
	return []string{string(UserSexUnknown), string(UserSexBoy), string(UserSexGirl)}
}

type User struct {
	mixin.Schema

	Optional bool

	NameUniq      bool
	NameDefaultFn func() string
	NameMaxLen    int
	NameMinLen    int
	NameRegexp    *regexp.Regexp

	Source           []string
	UsernamePassword bool
	Email            bool
	Phone            bool
	RealInfo         bool

	Avatar          bool
	AvatarDefaultFn func() string
	Remark          bool
	RemarkDefaultFn func() string

	Sex      bool
	Birthday bool
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
	}
}

func (u User) Fields() []ent.Field {
	f := make([]ent.Field, 0)

	fname := field.String("name").Comment("用户昵称")
	if u.NameDefaultFn != nil {
		fname = fname.DefaultFunc(u.NameDefaultFn)
	} else {
		fname = fname.Default("")
	}
	if u.NameMinLen > 0 {
		fname = fname.MinLen(u.NameMinLen)
	}
	if u.NameMaxLen > 0 {
		fname = fname.MaxLen(u.NameMaxLen)
	}
	if u.NameRegexp != nil {
		fname = fname.Match(u.NameRegexp)
	}
	if u.NameUniq {
		fname = fname.Unique()
	}
	f = append(f, fname)

	if u.Avatar {
		favatar := field.String("avatar").MaxLen(255).Comment("用户头像")
		if u.AvatarDefaultFn != nil {
			favatar = favatar.DefaultFunc(u.AvatarDefaultFn)
		} else {
			favatar = favatar.Default("")
		}
		f = append(f, favatar)
	}

	if u.Remark {
		fremark := field.String("remark").MaxLen(255).Comment("备注")
		if u.RemarkDefaultFn != nil {
			fremark = fremark.DefaultFunc(u.RemarkDefaultFn)
		} else {
			fremark = fremark.Default("该用户太懒了，什么都没写")
		}

		f = append(f, fremark)
	}

	if len(u.Source) > 1 && u.Source[0] == "known" {
		f = append(f, field.Enum("source").Values(u.Source...).Default("unknown").Comment("用户来源"))
		f = append(f, field.String("openid").Default("").Comment("用户的openid"))
		f = append(f, field.String("unionid").Default("").Comment("用户的unionid"))
	}
	if u.UsernamePassword {
		if u.Optional {
			f = append(f, field.String("username").MinLen(10).MaxLen(32).Comment("用户登录账号(数字字母下划线,10~32)").Optional().Unique().Match(regexp.MustCompile(`^[a-zA-Z0-9_]{10,32}$`)))
		} else {
			f = append(f, field.String("username").MinLen(10).MaxLen(32).Comment("用户登录账号(数字字母下划线,10~32)").Default(""))
		}
		f = append(f, field.String("password").MaxLen(32).Comment("用户密码").
			SchemaType(map[string]string{dialect.MySQL: "char(32)"}).Default(""))
	}

	if u.Phone {
		f = append(f, field.String("phone").MaxLen(11).Comment("手机号").Default(""))
	}
	if u.Email {
		f = append(f, field.String("email").MaxLen(64).Comment("用户邮箱，最长64位").Default(""))
	}
	if u.RealInfo {
		f = append(f, field.String("real_name").Comment("姓名").Default(""))
		f = append(f, field.String("id_card").MaxLen(18).Comment("身份证号").Default("").SchemaType(map[string]string{dialect.MySQL: "char(18)"}).Annotations(entsql.Annotation{Size: 18}))
		f = append(f, field.Enum("real_sex").GoType(UserSex("")).Comment("性别").Default("unknown"))
		if u.Optional {
			f = append(f, field.Time("real_birthday").Optional().Comment("生日"))
		} else {
			f = append(f, field.Int64("real_birthday").Default(0).Comment("生日"))
		}
		f = append(f, field.Uint8("real_province").Default(0).Comment("省级编码"))
		f = append(f, field.Uint8("real_city").Default(0).Comment("市级编码"))
		f = append(f, field.Uint8("real_district").Default(0).Comment("区级编码"))
	}

	if u.Sex {
		f = append(f, field.Enum("sex").GoType(UserSex("")).Comment("性别").Default("unknown"))
	}
	if u.Birthday {
		if u.Optional {
			f = append(f, field.Time("birthday").Optional().Comment("生日"))
		} else {
			f = append(f, field.Int64("birthday").Default(0).Comment("生日"))
		}
	}

	f = append(f, field.String("register_ip").Default("").Comment("注册ip"))
	return f
}

func (u User) Indexes() []ent.Index {
	idx := make([]ent.Index, 0)
	idx = append(idx, index.Fields("name"))
	if len(u.Source) > 1 && u.Source[0] == "known" {
		idx = append(idx, index.Fields("source", "openid"))
		idx = append(idx, index.Fields("source", "unionid"))
	}

	if u.UsernamePassword {
		idx = append(idx, index.Fields("username", "password"))
	}

	if u.Phone {
		idx = append(idx, index.Fields("phone"))
	}
	if u.Email {
		idx = append(idx, index.Fields("email"))
	}

	return idx
}

func (u User) Hooks() []ent.Hook {
	h := make([]ent.Hook, 0, 2)
	if u.UsernamePassword {
		h = append(h, func(next ent.Mutator) ent.Mutator {
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
		})
	}
	if u.RealInfo {
		h = append(h, func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
				if mutation.Op().Is(ent.OpCreate | ent.OpUpdate | ent.OpUpdateOne) {
					// 1. 从身份中提取数据
					idCard, ok := mutation.Field("id_card")
					if ok && idCard != nil && idCard.(string) != "" {
						id := idCard.(string)
						if len(id) != 18 {
							return nil, fmt.Errorf("id card length is not 18")
						}
						var err error
						switch id[17] {
						case '0', '2', '4', '6', '8':
							err = mutation.SetField("real_sex", UserSexGirl)
						default:
							err = mutation.SetField("real_sex", UserSexBoy)
						}
						if err != nil {
							return nil, fmt.Errorf("id set sex failed: %v", err)
						}
						birthday := id[6:14]
						birth, err := time.Parse("20060102", birthday)
						if err != nil {
							return nil, fmt.Errorf("parse birthday failed: %v", err)
						}
						if u.Optional {
							err = mutation.SetField("real_birthday", birth)
							if err != nil {
								return nil, fmt.Errorf("id set birthday failed: %v", err)
							}
						} else {
							err = mutation.SetField("real_birthday", birth.Unix())
							if err != nil {
								return nil, fmt.Errorf("id set birthday failed: %v", err)
							}
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
		})
	}

	return h
}
