package user

import (
	"context"
	"gorm.io/gorm"
	"k3gin/app/gormx"
	"k3gin/app/schema"
	"k3gin/app/util/structure"
)

type User struct {
	ID       uint64 `gorm:"column:user_id"`
	UserName string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	RealName string `gorm:"column:realname"`
	Status   int    `gorm:"column:status"`
}

// TableName 实现此方法可以解决结构体名称和数据库表对应关系的问题
func (u User) TableName() string {
	return "admin_users"
}

// ToSchemaUser 将数据库映射的结构体 转换乘 schema user
func (u User) ToSchemaUser() *schema.User {
	item := new(schema.User)
	structure.Copy(u, item)
	return item
}

type Users []*User

func (u Users) ToSchemaUsers() []*schema.User {

	list := make([]*schema.User, len(u))

	for i, item := range u {
		list[i] = item.ToSchemaUser()
	}
	return list
}

func GetUserWriteDB(ctx context.Context, db *gormx.DB) *gorm.DB {
	return db.WDB
}

func GetUserReadDB(ctx context.Context, db *gormx.DB) *gorm.DB {
	return db.RDB
}
