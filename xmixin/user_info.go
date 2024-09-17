package xmixin

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"regexp"
)

type UserInfoSex string

const (
	UserInfoSexUnknown UserInfoSex = "unknown"
	UserInfoSexBoy     UserInfoSex = "boy"
	UserInfoSexGirl    UserInfoSex = "girl"
)

func (UserInfoSex) Values() []string {
	return []string{string(UserInfoSexUnknown), string(UserInfoSexBoy), string(UserInfoSexGirl)}
}

type UserInfo struct {
	mixin.Schema
	Source []string

	NameConf struct {
		Uniq      bool
		MaxLen    int
		MinLen    int
		Regexp    *regexp.Regexp
		DefaultFn func() string
	}

	Sex      bool
	Birthday bool
}

func (u UserInfo) Fields() []ent.Field {
	f := make([]ent.Field, 0, 7)

	fname := field.String("name").Comment("用户昵称")
	if u.NameConf.DefaultFn != nil {
		fname = fname.DefaultFunc(u.NameConf.DefaultFn)
	} else {
		fname = fname.Default("")
	}
	if u.NameConf.MinLen > 0 {
		fname = fname.MinLen(u.NameConf.MinLen)
	}
	if u.NameConf.MaxLen > 0 {
		fname = fname.MaxLen(u.NameConf.MaxLen)
	}
	if u.NameConf.Regexp != nil {
		fname = fname.Match(u.NameConf.Regexp)
	}
	if u.NameConf.Uniq {
		fname = fname.Unique()
	}
	f = append(f, fname)
	f = append(f, field.String("avatar").MaxLen(255).Comment("用户头像").Default(""))
	f = append(f, field.String("remark").MaxLen(255).Comment("签名/标记/备注").Default(""))
	f = append(f, field.String("phone").MaxLen(11).Comment("手机号").Default("").SchemaType(map[string]string{dialect.MySQL: "char(11)"}))
	f = append(f, field.String("email").MaxLen(64).Comment("用户邮箱").Default(""))
	f = append(f, field.Enum("sex").GoType(UserInfoSex("")).Comment("性别").Default("unknown"))
	f = append(f, field.Time("birthday").Optional().Comment("生日"))

	return f
}

func (u UserInfo) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
		index.Fields("phone"),
		index.Fields("email"),
	}
}

func (u UserInfo) Hooks() []ent.Hook {
	h := make([]ent.Hook, 0, 2)
	h = append(h, hookVerifyRegexpFunc("email", regexp.MustCompile("^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$")))
	h = append(h, hookVerifyRegexpFunc("phone", regexp.MustCompile("^1[3456789][0-9]{9}$")))
	return h
}
