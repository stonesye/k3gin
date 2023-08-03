package service

import (
	"context"
	"github.com/google/wire"
	"k3gin/app/dao/user"
	"k3gin/app/schema"
)

type UserSrv struct {
	UserRepo *user.UserRepo
}

var UserSrvSet = wire.NewSet(wire.Struct(new(UserSrv), "*"))

// Query  params : form提交的数据 option : 查询条件
func (u *UserSrv) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	return u.UserRepo.Query(ctx, params, opts...)
}
