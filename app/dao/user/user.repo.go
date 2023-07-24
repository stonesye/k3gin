package user

import (
	"context"
	"github.com/google/wire"
	"k3gin/app/gormx"
	"k3gin/app/logger"
	"k3gin/app/schema"
)

type UserRepo struct {
	DB *gormx.DB
}

var UserRepoSet = wire.NewSet(wire.Struct(new(UserRepo), "*"))

// Query  params : form表单提交的数据， opts ：查询条件
func (u *UserRepo) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {

	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetUserReadDB(ctx, u.DB)

	if v := params.UserName; v != "" {
		db = db.Where("username = ? ", v)
	}

	if v := params.Status; v > 0 {
		db = db.Where("status = ?", v)
	}

	// 模糊查询
	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("username LIKE ? OR realname LIKE ?", v, v)
	}

	if len(opt.SelectFields) > 0 {
		db = db.Select(opt.SelectFields)
	}

	if len(opt.OrderFields) > 0 {
		db = db.Order(opt.OrderFields)
	}

	var list Users
	err := db.Find(&list).Error
	if err != nil {
		return nil, err
	}

	logger.WithContext(ctx).Infof("%v", list)

	return &schema.UserQueryResult{Data: list.ToSchemaUsers()}, nil
}
